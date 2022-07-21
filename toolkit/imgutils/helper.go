package imgutils

import (
	"container/heap"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"sort"
)

const (
	numDimensions = 3
)

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

type point [numDimensions]int

type block struct {
	minCorner, maxCorner point
	points               []point
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

func newBlock(p []point) *block {
	return &block{
		minCorner: point{0x00, 0x00, 0x00},
		maxCorner: point{0xFF, 0xFF, 0xFF},
		points:    p,
	}
}

func (b *block) longestSideIndex() int {
	m := b.maxCorner[0] - b.minCorner[0]
	maxIndex := 0
	for i := 1; i < numDimensions; i++ {
		diff := b.maxCorner[i] - b.minCorner[i]
		if diff > m {
			m = diff
			maxIndex = i
		}
	}
	return maxIndex
}

func (b *block) longestSideLength() int {
	i := b.longestSideIndex()
	return b.maxCorner[i] - b.minCorner[i]
}

func (b *block) shrink() {
	for j := 0; j < numDimensions; j++ {
		b.minCorner[j] = b.points[0][j]
		b.maxCorner[j] = b.points[0][j]
	}
	for i := 1; i < len(b.points); i++ {
		for j := 0; j < numDimensions; j++ {
			b.minCorner[j] = min(b.minCorner[j], b.points[i][j])
			b.maxCorner[j] = max(b.maxCorner[j], b.points[i][j])
		}
	}
}

type pointSorter struct {
	points []point
	by     func(p1, p2 *point) bool
}

func (p *pointSorter) Len() int {
	return len(p.points)
}

func (p *pointSorter) Swap(i, j int) {
	p.points[i], p.points[j] = p.points[j], p.points[i]
}

func (p *pointSorter) Less(i, j int) bool {
	return p.by(&p.points[i], &p.points[j])
}

// A priorityQueue implements heap.Interface and holds blocks.
type priorityQueue []*block

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].longestSideLength() > pq[j].longestSideLength()
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*block)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[:n-1]
	return item
}

func (pq *priorityQueue) top() interface{} {
	n := len(*pq)
	if n == 0 {
		return nil
	}
	return (*pq)[n-1]
}

// clip clips r against each image's bounds (after translating into
// the destination image's co-ordinate space) and shifts the point
// sp by the same amount as the change in r.Min.
func clip(dst draw.Image, r *image.Rectangle, src image.Image, sp *image.Point) {
	orig := r.Min
	*r = r.Intersect(dst.Bounds())
	*r = r.Intersect(src.Bounds().Add(orig.Sub(*sp)))
	dx := r.Min.X - orig.X
	dy := r.Min.Y - orig.Y
	if dx == 0 && dy == 0 {
		return
	}
	(*sp).X += dx
	(*sp).Y += dy
}

// MedianCutQuantizer constructs a palette with a maximum of
// NumColor colors by iteratively splitting clusters of color
// points mapped on a three-dimensional (RGB) Euclidian space.
// Once the number of clusters is within the specified bounds,
// the resulting color is computed by averaging those within
// each grouping.
type MedianCutQuantizer struct {
	NumColor int
}

func (q *MedianCutQuantizer) medianCut(points []point) color.Palette {
	if q.NumColor == 0 {
		return color.Palette{}
	}

	initialBlock := newBlock(points)
	initialBlock.shrink()
	pq := &priorityQueue{}
	heap.Init(pq)
	heap.Push(pq, initialBlock)

	for pq.Len() < q.NumColor && len(pq.top().(*block).points) > 1 {
		longestBlock := heap.Pop(pq).(*block)
		points := longestBlock.points
		li := longestBlock.longestSideIndex()
		// TODO: Instead of sorting the entire slice, finding the median using an
		// algorithm like introselect would give much better performance.
		sort.Sort(&pointSorter{
			points: points,
			by:     func(p1, p2 *point) bool { return p1[li] < p2[li] },
		})
		median := len(points) / 2
		block1 := newBlock(points[:median])
		block2 := newBlock(points[median:])
		block1.shrink()
		block2.shrink()
		heap.Push(pq, block1)
		heap.Push(pq, block2)
	}

	palette := make(color.Palette, q.NumColor)
	var n int
	for n = 0; pq.Len() > 0; n++ {
		block := heap.Pop(pq).(*block)
		var sum [numDimensions]int
		for i := 0; i < len(block.points); i++ {
			for j := 0; j < numDimensions; j++ {
				sum[j] += block.points[i][j]
			}
		}
		palette[n] = color.RGBA64{
			R: uint16(sum[0] / len(block.points)),
			G: uint16(sum[1] / len(block.points)),
			B: uint16(sum[2] / len(block.points)),
			A: 0xFFFF,
		}
	}
	// Trim to only the colors present in the image, which
	// could be less than NumColor.
	return palette[:n]
}

func (q *MedianCutQuantizer) Quantize(dst *image.Paletted, r image.Rectangle, src image.Image, sp image.Point) {
	clip(dst, &r, src, &sp)
	if r.Empty() {
		return
	}

	points := make([]point, r.Dx()*r.Dy())
	colorSet := make(map[uint32]color.Color, q.NumColor)
	i := 0
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			c := src.At(x, y)
			r, g, b, _ := c.RGBA()
			colorSet[(r>>8)<<16|(g>>8)<<8|b>>8] = c
			points[i][0] = int(r)
			points[i][1] = int(g)
			points[i][2] = int(b)
			i++
		}
	}
	if len(colorSet) <= q.NumColor {
		// No need to quantize since the total number of colors
		// fits within the palette.
		dst.Palette = make(color.Palette, len(colorSet))
		i := 0
		for _, c := range colorSet {
			dst.Palette[i] = c
			i++
		}
	} else {
		dst.Palette = q.medianCut(points)
	}

	for y := 0; y < r.Dy(); y++ {
		for x := 0; x < r.Dx(); x++ {
			// TODO: this should be done more efficiently.
			dst.Set(sp.X+x, sp.Y+y, src.At(r.Min.X+x, r.Min.Y+y))
		}
	}
}

func inPalette(p color.Palette, c color.Color) int {
	ret := -1
	for i, v := range p {
		if v == c {
			return i
		}
	}
	return ret
}

func getSubPalette(m image.Image) color.Palette {
	p := color.Palette{color.RGBA{0x00, 0x00, 0x00, 0x00}}
	p9 := color.Palette(palette.Plan9)
	b := m.Bounds()
	black := false
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := m.At(x, y)
			cc := p9.Convert(c)
			if cc == p9[0] {
				black = true
			}
			if inPalette(p, cc) == -1 {
				p = append(p, cc)
			}
		}
	}
	if len(p) < 256 && black == true {
		p[0] = color.RGBA{0x00, 0x00, 0x00, 0x00} // transparent
		p = append(p, p9[0])
	}
	return p
}

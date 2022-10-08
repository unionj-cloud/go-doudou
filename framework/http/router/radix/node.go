package radix

import (
	"sort"
	"strings"

	gstrings "github.com/savsgio/gotils/strings"
	"github.com/valyala/bytebufferpool"
	"github.com/valyala/fasthttp"
)

func newNode(path string) *node {
	return &node{
		nType: static,
		path:  path,
	}
}

// conflict raises a panic with some details
func (n *nodeWildcard) conflict(path, fullPath string) error {
	prefix := fullPath[:strings.LastIndex(fullPath, path)] + n.path

	return newRadixError(errWildcardConflict, path, fullPath, n.path, prefix)
}

// wildPathConflict raises a panic with some details
func (n *node) wildPathConflict(path, fullPath string) error {
	pathSeg := strings.SplitN(path, "/", 2)[0]
	prefix := fullPath[:strings.LastIndex(fullPath, path)] + n.path

	return newRadixError(errWildPathConflict, pathSeg, fullPath, n.path, prefix)
}

// clone clones the current node in a new pointer
func (n node) clone() *node {
	cloneNode := new(node)
	cloneNode.nType = n.nType
	cloneNode.path = n.path
	cloneNode.tsr = n.tsr
	cloneNode.handler = n.handler

	if len(n.children) > 0 {
		cloneNode.children = make([]*node, len(n.children))

		for i, child := range n.children {
			cloneNode.children[i] = child.clone()
		}
	}

	if n.wildcard != nil {
		cloneNode.wildcard = &nodeWildcard{
			path:     n.wildcard.path,
			paramKey: n.wildcard.paramKey,
			handler:  n.wildcard.handler,
		}
	}

	if len(n.paramKeys) > 0 {
		cloneNode.paramKeys = make([]string, len(n.paramKeys))
		copy(cloneNode.paramKeys, n.paramKeys)
	}

	cloneNode.paramRegex = n.paramRegex

	return cloneNode
}

func (n *node) split(i int) {
	cloneChild := n.clone()
	cloneChild.nType = static
	cloneChild.path = cloneChild.path[i:]
	cloneChild.paramKeys = nil
	cloneChild.paramRegex = nil

	n.path = n.path[:i]
	n.handler = nil
	n.tsr = false
	n.wildcard = nil
	n.children = append(n.children[:0], cloneChild)
}

func (n *node) findEndIndexAndValues(path string) (int, []string) {
	index := n.paramRegex.FindStringSubmatchIndex(path)
	if len(index) == 0 {
		return -1, nil
	}

	end := index[1]

	index = index[2:]
	values := make([]string, len(index)/2)

	i := 0
	for j := range index {
		if (j+1)%2 != 0 {
			continue
		}

		values[i] = gstrings.Copy(path[index[j-1]:index[j]])

		i++
	}

	return end, values
}

func (n *node) setHandler(handler fasthttp.RequestHandler, fullPath string) (*node, error) {
	if n.handler != nil || n.tsr {
		return n, newRadixError(errSetHandler, fullPath)
	}

	n.handler = handler
	foundTSR := false

	// Set TSR in method
	for i := range n.children {
		child := n.children[i]

		if child.path != "/" {
			continue
		}

		child.tsr = true
		foundTSR = true

		break
	}

	if n.path != "/" && !foundTSR {
		childTSR := newNode("/")
		childTSR.tsr = true
		n.children = append(n.children, childTSR)
	}

	return n, nil
}

func (n *node) insert(path, fullPath string, handler fasthttp.RequestHandler) (*node, error) {
	end := segmentEndIndex(path, true)
	child := newNode(path)

	wp := findWildPath(path, fullPath)
	if wp != nil {
		j := end
		if wp.start > 0 {
			j = wp.start
		}

		child.path = path[:j]

		if wp.start > 0 {
			n.children = append(n.children, child)

			return child.insert(path[j:], fullPath, handler)
		}

		switch wp.pType {
		case param:
			n.hasWildChild = true

			child.nType = wp.pType
			child.paramKeys = wp.keys
			child.paramRegex = wp.regex
		case wildcard:
			if len(path) == end && n.path[len(n.path)-1] != '/' {
				return nil, newRadixError(errWildcardSlash, fullPath)
			} else if len(path) != end {
				return nil, newRadixError(errWildcardNotAtEnd, fullPath)
			}

			if n.path != "/" && n.path[len(n.path)-1] == '/' {
				n.split(len(n.path) - 1)
				n.tsr = true

				n = n.children[0]
			}

			if n.wildcard != nil {
				if n.wildcard.path == path {
					return n, newRadixError(errSetWildcardHandler, fullPath)
				}

				return nil, n.wildcard.conflict(path, fullPath)
			}

			n.wildcard = &nodeWildcard{
				path:     wp.path,
				paramKey: wp.keys[0],
				handler:  handler,
			}

			return n, nil
		}

		path = path[wp.end:]

		if len(path) > 0 {
			n.children = append(n.children, child)

			return child.insert(path, fullPath, handler)
		}
	}

	child.handler = handler
	n.children = append(n.children, child)

	if child.path == "/" {
		// Add TSR when split a edge and the remain path to insert is "/"
		n.tsr = true
	} else if strings.HasSuffix(child.path, "/") {
		child.split(len(child.path) - 1)
		child.tsr = true
	} else {
		childTSR := newNode("/")
		childTSR.tsr = true
		child.children = append(child.children, childTSR)
	}

	return child, nil
}

// add adds the handler to node for the given path
func (n *node) add(path, fullPath string, handler fasthttp.RequestHandler) (*node, error) {
	if len(path) == 0 {
		return n.setHandler(handler, fullPath)
	}

	for _, child := range n.children {
		i := longestCommonPrefix(path, child.path)
		if i == 0 {
			continue
		}

		switch child.nType {
		case static:
			if len(child.path) > i {
				child.split(i)
			}

			if len(path) > i {
				return child.add(path[i:], fullPath, handler)
			}
		case param:
			wp := findWildPath(path, fullPath)

			isParam := wp.start == 0 && wp.pType == param
			hasHandler := child.handler != nil || handler == nil

			if len(path) == wp.end && isParam && hasHandler {
				// The current segment is a param and it's duplicated
				if child.path == path {
					return child, newRadixError(errSetHandler, fullPath)
				}

				return nil, child.wildPathConflict(path, fullPath)
			}

			if len(path) > i {
				if child.path == wp.path {
					return child.add(path[i:], fullPath, handler)
				}

				return n.insert(path, fullPath, handler)
			}
		}

		if path == "/" {
			n.tsr = true
		}

		return child.setHandler(handler, fullPath)
	}

	return n.insert(path, fullPath, handler)
}

func (n *node) getFromChild(path string, ctx *fasthttp.RequestCtx) (fasthttp.RequestHandler, bool) {
	for _, child := range n.children {
		switch child.nType {
		case static:

			// Checks if the first byte is equal
			// It's faster than compare strings
			if path[0] != child.path[0] {
				continue
			}

			if len(path) > len(child.path) {
				if path[:len(child.path)] != child.path {
					continue
				}

				h, tsr := child.getFromChild(path[len(child.path):], ctx)
				if h != nil || tsr {
					return h, tsr
				}
			} else if path == child.path {
				switch {
				case child.tsr:
					return nil, true
				case child.handler != nil:
					return child.handler, false
				case child.wildcard != nil:
					if ctx != nil {
						ctx.SetUserValue(child.wildcard.paramKey, "")
					}

					return child.wildcard.handler, false
				}

				return nil, false
			}

		case param:
			end := segmentEndIndex(path, false)
			values := []string{gstrings.Copy(path[:end])}

			if child.paramRegex != nil {
				end, values = child.findEndIndexAndValues(path[:end])
				if end == -1 {
					continue
				}
			}

			if len(path) > end {
				h, tsr := child.getFromChild(path[end:], ctx)
				if tsr {
					return nil, tsr
				} else if h != nil {
					if ctx != nil {
						for i, key := range child.paramKeys {
							ctx.SetUserValue(key, values[i])
						}
					}

					return h, false
				}

			} else if len(path) == end {
				switch {
				case child.tsr:
					return nil, true
				case child.handler == nil:
					// try another child
					continue
				case ctx != nil:
					for i, key := range child.paramKeys {
						ctx.SetUserValue(key, values[i])
					}
				}

				return child.handler, false
			}

		default:
			panic("invalid node type")
		}
	}

	if n.wildcard != nil {
		if ctx != nil {
			ctx.SetUserValue(n.wildcard.paramKey, gstrings.Copy(path))
		}

		return n.wildcard.handler, false
	}

	return nil, false
}

func (n *node) find(path string, buf *bytebufferpool.ByteBuffer) (bool, bool) {
	if len(path) > len(n.path) {
		if !strings.EqualFold(path[:len(n.path)], n.path) {
			return false, false
		}

		path = path[len(n.path):]
		buf.WriteString(n.path)

		found, tsr := n.findFromChild(path, buf)
		if found {
			return found, tsr
		}

		bufferRemoveString(buf, n.path)

	} else if strings.EqualFold(path, n.path) {
		buf.WriteString(n.path)

		if n.tsr {
			if n.path == "/" {
				bufferRemoveString(buf, n.path)
			} else {
				buf.WriteByte('/')
			}

			return true, true
		}

		if n.handler != nil {
			return true, false
		} else {
			bufferRemoveString(buf, n.path)
		}
	}

	return false, false
}

func (n *node) findFromChild(path string, buf *bytebufferpool.ByteBuffer) (bool, bool) {
	for _, child := range n.children {
		switch child.nType {
		case static:
			found, tsr := child.find(path, buf)
			if found {
				return found, tsr
			}

		case param:
			end := segmentEndIndex(path, false)

			if child.paramRegex != nil {
				end, _ = child.findEndIndexAndValues(path[:end])
				if end == -1 {
					continue
				}
			}

			buf.WriteString(path[:end])

			if len(path) > end {
				found, tsr := child.findFromChild(path[end:], buf)
				if found {
					return found, tsr
				}

			} else if len(path) == end {
				if child.tsr {
					buf.WriteByte('/')

					return true, true
				}

				if child.handler != nil {
					return true, false
				}
			}

			bufferRemoveString(buf, path[:end])

		default:
			panic("invalid node type")
		}
	}

	if n.wildcard != nil {
		buf.WriteString(path)

		return true, false
	}

	return false, false
}

// sort sorts the current node and their children
func (n *node) sort() {
	for _, child := range n.children {
		child.sort()
	}

	sort.Sort(n)
}

// Len returns the total number of children the node has
func (n *node) Len() int {
	return len(n.children)
}

// Swap swaps the order of children nodes
func (n *node) Swap(i, j int) {
	n.children[i], n.children[j] = n.children[j], n.children[i]
}

// Less checks if the node 'i' has less priority than the node 'j'
func (n *node) Less(i, j int) bool {
	if n.children[i].nType < n.children[j].nType {
		return true
	} else if n.children[i].nType > n.children[j].nType {
		return false
	}

	return len(n.children[i].children) > len(n.children[j].children)
}

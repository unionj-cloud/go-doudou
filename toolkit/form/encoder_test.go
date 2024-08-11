package form

import (
	"errors"
	"github.com/goccy/go-reflect"
	"strings"
	"testing"
	"time"

	. "github.com/go-playground/assert/v2"
)

// NOTES:
// - Run "go test" to run tests
// - Run "gocov test | gocov report" to report on test converage by file
// - Run "gocov test | gocov annotate -" to report on all code and functions, those ,marked with "MISS" were never called
//
// or
//
// -- may be a good idea to change to output path to somewherelike /tmp
// go test -coverprofile cover.out && go tool cover -html=cover.out -o cover.html
//
//
// go test -cpuprofile cpu.out
// ./validator.test -test.bench=. -test.cpuprofile=cpu.prof
// go tool pprof validator.test cpu.prof
//
//
// go test -memprofile mem.out

func TestEncoderInt(t *testing.T) {

	type TestInt struct {
		Int              int
		Int8             int8
		Int16            int16
		Int32            int32
		Int64            int64
		IntPtr           *int
		Int8Ptr          *int8
		Int16Ptr         *int16
		Int32Ptr         *int32
		Int64Ptr         *int64
		IntArray         []int
		IntPtrArray      []*int
		IntArrayArray    [][]int
		IntPtrArrayArray [][]*int
		IntMap           map[int]int
		IntPtrMap        map[*int]*int
		NoValue          int
		NoPtrValue       *int
	}

	i := int(3)
	i8 := int8(3)
	i16 := int16(3)
	i32 := int32(3)
	i64 := int64(3)

	zero := int(0)
	one := int(1)
	two := int(2)
	three := int(3)

	test := TestInt{
		Int:              i,
		Int8:             i8,
		Int16:            i16,
		Int32:            i32,
		Int64:            i64,
		IntPtr:           &i,
		Int8Ptr:          &i8,
		Int16Ptr:         &i16,
		Int32Ptr:         &i32,
		Int64Ptr:         &i64,
		IntArray:         []int{one, two, three},
		IntPtrArray:      []*int{&one, &two, &three},
		IntArrayArray:    [][]int{{one, zero, three}},
		IntPtrArrayArray: [][]*int{{&one, &zero, &three}},
		IntMap:           map[int]int{one: three, zero: two},
		IntPtrMap:        map[*int]*int{&one: &three, &zero: &two},
	}

	encoder := NewEncoder()
	values, errs := encoder.Encode(test)
	Equal(t, errs, nil)
	Equal(t, len(values), 25)

	val, ok := values["Int8"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["Int8"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["Int16"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["Int32"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["Int64"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["Int8"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["Int8Ptr"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["Int16Ptr"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["Int32Ptr"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["Int64Ptr"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["IntArray"]
	Equal(t, ok, true)
	Equal(t, len(val), 3)
	Equal(t, val[0], "1")
	Equal(t, val[1], "2")
	Equal(t, val[2], "3")

	val, ok = values["IntPtrArray[0]"]
	Equal(t, ok, true)
	Equal(t, val[0], "1")

	val, ok = values["IntPtrArray[1]"]
	Equal(t, ok, true)
	Equal(t, val[0], "2")

	val, ok = values["IntPtrArray[2]"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["IntArrayArray[0][0]"]
	Equal(t, ok, true)
	Equal(t, val[0], "1")

	val, ok = values["IntArrayArray[0][1]"]
	Equal(t, ok, true)
	Equal(t, val[0], "0")

	val, ok = values["IntArrayArray[0][2]"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["IntPtrArrayArray[0][0]"]
	Equal(t, ok, true)
	Equal(t, val[0], "1")

	val, ok = values["IntPtrArrayArray[0][1]"]
	Equal(t, ok, true)
	Equal(t, val[0], "0")

	val, ok = values["IntPtrArrayArray[0][2]"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["IntMap[0]"]
	Equal(t, ok, true)
	Equal(t, val[0], "2")

	val, ok = values["IntMap[1]"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["IntPtrMap[0]"]
	Equal(t, ok, true)
	Equal(t, val[0], "2")

	val, ok = values["IntPtrMap[1]"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["NoValue"]
	Equal(t, ok, true)
	Equal(t, val[0], "0")

	_, ok = values["NoPtrValue"]
	Equal(t, ok, false)
}

func TestEncoderUint(t *testing.T) {

	type TestUint struct {
		Uint              uint
		Uint8             uint8
		Uint16            uint16
		Uint32            uint32
		Uint64            uint64
		UintPtr           *uint
		Uint8Ptr          *uint8
		Uint16Ptr         *uint16
		Uint32Ptr         *uint32
		Uint64Ptr         *uint64
		UintArray         []uint
		UintPtrArray      []*uint
		UintArrayArray    [][]uint
		UintPtrArrayArray [][]*uint
		UintMap           map[uint]uint
		UintPtrMap        map[*uint]*uint
		NoValue           uint
		NoPtrValue        *uint
	}

	i := uint(3)
	i8 := uint8(3)
	i16 := uint16(3)
	i32 := uint32(3)
	i64 := uint64(3)

	zero := uint(0)
	one := uint(1)
	two := uint(2)
	three := uint(3)

	test := TestUint{
		Uint:              i,
		Uint8:             i8,
		Uint16:            i16,
		Uint32:            i32,
		Uint64:            i64,
		UintPtr:           &i,
		Uint8Ptr:          &i8,
		Uint16Ptr:         &i16,
		Uint32Ptr:         &i32,
		Uint64Ptr:         &i64,
		UintArray:         []uint{one, two, three},
		UintPtrArray:      []*uint{&one, &two, &three},
		UintArrayArray:    [][]uint{{one, zero, three}},
		UintPtrArrayArray: [][]*uint{{&one, &zero, &three}},
		UintMap:           map[uint]uint{one: three, zero: two},
		UintPtrMap:        map[*uint]*uint{&one: &three, &zero: &two},
	}

	encoder := NewEncoder()
	values, errs := encoder.Encode(test)

	Equal(t, errs, nil)
	Equal(t, len(values), 25)

	val, ok := values["Uint8"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["Uint8"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["Uint16"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["Uint32"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["Uint64"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["Uint8"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["Uint8Ptr"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["Uint16Ptr"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["Uint32Ptr"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["Uint64Ptr"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["UintArray"]
	Equal(t, ok, true)
	Equal(t, len(val), 3)
	Equal(t, val[0], "1")
	Equal(t, val[1], "2")
	Equal(t, val[2], "3")

	val, ok = values["UintPtrArray[0]"]
	Equal(t, ok, true)
	Equal(t, val[0], "1")

	val, ok = values["UintPtrArray[1]"]
	Equal(t, ok, true)
	Equal(t, val[0], "2")

	val, ok = values["UintPtrArray[2]"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["UintArrayArray[0][0]"]
	Equal(t, ok, true)
	Equal(t, val[0], "1")

	val, ok = values["UintArrayArray[0][1]"]
	Equal(t, ok, true)
	Equal(t, val[0], "0")

	val, ok = values["UintArrayArray[0][2]"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["UintPtrArrayArray[0][0]"]
	Equal(t, ok, true)
	Equal(t, val[0], "1")

	val, ok = values["UintPtrArrayArray[0][1]"]
	Equal(t, ok, true)
	Equal(t, val[0], "0")

	val, ok = values["UintPtrArrayArray[0][2]"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["UintMap[0]"]
	Equal(t, ok, true)
	Equal(t, val[0], "2")

	val, ok = values["UintMap[1]"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["UintPtrMap[0]"]
	Equal(t, ok, true)
	Equal(t, val[0], "2")

	val, ok = values["UintPtrMap[1]"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["NoValue"]
	Equal(t, ok, true)
	Equal(t, val[0], "0")

	_, ok = values["NoPtrValue"]
	Equal(t, ok, false)
}

func TestEncoderString(t *testing.T) {

	type TestString struct {
		String              string
		StringPtr           *string
		StringArray         []string
		StringPtrArray      []*string
		StringArrayArray    [][]string
		StringPtrArrayArray [][]*string
		StringMap           map[string]string
		StringPtrMap        map[*string]*string
		NoValue             string
	}

	one := "1"
	two := "2"
	three := "3"

	test := TestString{
		String:              three,
		StringPtr:           &two,
		StringArray:         []string{one, "", three},
		StringPtrArray:      []*string{&one, nil, &three},
		StringArrayArray:    [][]string{{one, "", three}, nil, {one}},
		StringPtrArrayArray: [][]*string{{&one, nil, &three}, nil, {&one}},
		StringMap:           map[string]string{one: three, three: two},
		StringPtrMap:        map[*string]*string{&one: &three, &three: &two},
	}

	encoder := NewEncoder()
	values, errs := encoder.Encode(test)

	Equal(t, errs, nil)
	Equal(t, len(values), 17)

	val, ok := values["String"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["StringPtr"]
	Equal(t, ok, true)
	Equal(t, val[0], "2")

	val, ok = values["StringArray"]
	Equal(t, ok, true)
	Equal(t, len(val), 3)
	Equal(t, val[0], "1")
	Equal(t, val[1], "")
	Equal(t, val[2], "3")

	val, ok = values["StringPtr"]
	Equal(t, ok, true)
	Equal(t, val[0], "2")

	val, ok = values["StringPtrArray[0]"]
	Equal(t, ok, true)
	Equal(t, val[0], "1")

	_, ok = values["StringPtrArray[1]"]
	Equal(t, ok, false)

	val, ok = values["StringPtrArray[2]"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["StringPtrArray[2]"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["StringArrayArray[0][0]"]
	Equal(t, ok, true)
	Equal(t, val[0], "1")

	val, ok = values["StringArrayArray[0][1]"]
	Equal(t, ok, true)
	Equal(t, val[0], "")

	val, ok = values["StringArrayArray[0][2]"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	_, ok = values["StringArrayArray[1][1]"]
	Equal(t, ok, false)

	val, ok = values["StringArrayArray[2][0]"]
	Equal(t, ok, true)
	Equal(t, val[0], "1")

	val, ok = values["StringPtrArrayArray[0][0]"]
	Equal(t, ok, true)
	Equal(t, val[0], "1")

	_, ok = values["StringPtrArrayArray[0][1]"]
	Equal(t, ok, false)

	val, ok = values["StringPtrArrayArray[0][2]"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	_, ok = values["StringPtrArrayArray[1][1]"]
	Equal(t, ok, false)

	val, ok = values["StringPtrArrayArray[2][0]"]
	Equal(t, ok, true)
	Equal(t, val[0], "1")

	val, ok = values["StringMap[1]"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["StringMap[3]"]
	Equal(t, ok, true)
	Equal(t, val[0], "2")

	val, ok = values["StringPtrMap[1]"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["StringPtrMap[3]"]
	Equal(t, ok, true)
	Equal(t, val[0], "2")

	val, ok = values["NoValue"]
	Equal(t, ok, true)
	Equal(t, val[0], "")
}

func TestEncoderFloat(t *testing.T) {

	type TestFloat struct {
		Float32              float32
		Float32Ptr           *float32
		Float64              float64
		Float64Ptr           *float64
		Float32Array         []float32
		Float64Array         []float64
		Float32PtrArray      []*float32
		Float64PtrArray      []*float64
		Float32ArrayArray    [][]float32
		Float64ArrayArray    [][]float64
		Float32PtrArrayArray [][]*float32
		Float64PtrArrayArray [][]*float64
		Float32Map           map[float32]float32
		Float64Map           map[float64]float64
		Float32PtrMap        map[*float32]*float32
		Float64PtrMap        map[*float64]*float64
		NoValue              float32
	}

	one32 := float32(1.1)
	two32 := float32(2.2)
	three32 := float32(3.3)
	one64 := float64(1.1)
	two64 := float64(2.2)
	three64 := float64(3.3)

	test := TestFloat{
		Float32:              three32,
		Float32Ptr:           &three32,
		Float64:              three64,
		Float64Ptr:           &three64,
		Float32Array:         []float32{one32, two32, three32},
		Float64Array:         []float64{one64, two64, three64},
		Float32PtrArray:      []*float32{&one32, &two32, &three32},
		Float64PtrArray:      []*float64{&one64, &two64, &three64},
		Float32ArrayArray:    [][]float32{{one32, 0, three32}, nil, {one32}},
		Float64ArrayArray:    [][]float64{{one64, 0, three64}, nil, {one64}},
		Float32PtrArrayArray: [][]*float32{{&one32, nil, &three32}, nil, {&one32}},
		Float64PtrArrayArray: [][]*float64{{&one64, nil, &three64}, nil, {&one64}},
		Float32Map:           map[float32]float32{one32: three32, three32: two32},
		Float64Map:           map[float64]float64{one64: three64, three64: two64},
		Float32PtrMap:        map[*float32]*float32{&one32: &three32, &three32: &two32},
		Float64PtrMap:        map[*float64]*float64{&one64: &three64, &three64: &two64},
	}

	encoder := NewEncoder()
	values, errs := encoder.Encode(test)

	Equal(t, errs, nil)
	Equal(t, len(values), 35)

	val, ok := values["Float32"]
	Equal(t, ok, true)
	Equal(t, val[0], "3.3")

	val, ok = values["Float32Ptr"]
	Equal(t, ok, true)
	Equal(t, val[0], "3.3")

	val, ok = values["Float64"]
	Equal(t, ok, true)
	Equal(t, val[0], "3.3")

	val, ok = values["Float64Ptr"]
	Equal(t, ok, true)
	Equal(t, val[0], "3.3")

	val, ok = values["Float32Array"]
	Equal(t, ok, true)
	Equal(t, len(val), 3)
	Equal(t, val[0], "1.1")
	Equal(t, val[1], "2.2")
	Equal(t, val[2], "3.3")

	val, ok = values["Float64Array"]
	Equal(t, ok, true)
	Equal(t, len(val), 3)
	Equal(t, val[0], "1.1")
	Equal(t, val[1], "2.2")
	Equal(t, val[2], "3.3")

	val, ok = values["Float32PtrArray[0]"]
	Equal(t, ok, true)
	Equal(t, val[0], "1.1")

	val, ok = values["Float32PtrArray[1]"]
	Equal(t, ok, true)
	Equal(t, val[0], "2.2")

	val, ok = values["Float32PtrArray[2]"]
	Equal(t, ok, true)
	Equal(t, val[0], "3.3")

	val, ok = values["Float64PtrArray[0]"]
	Equal(t, ok, true)
	Equal(t, val[0], "1.1")

	val, ok = values["Float64PtrArray[1]"]
	Equal(t, ok, true)
	Equal(t, val[0], "2.2")

	val, ok = values["Float64PtrArray[2]"]
	Equal(t, ok, true)
	Equal(t, val[0], "3.3")

	val, ok = values["Float32ArrayArray[0][0]"]
	Equal(t, ok, true)
	Equal(t, val[0], "1.1")

	val, ok = values["Float32ArrayArray[0][1]"]
	Equal(t, ok, true)
	Equal(t, val[0], "0")

	val, ok = values["Float32ArrayArray[0][2]"]
	Equal(t, ok, true)
	Equal(t, val[0], "3.3")

	_, ok = values["Float32ArrayArray[1][0]"]
	Equal(t, ok, false)

	val, ok = values["Float32ArrayArray[2][0]"]
	Equal(t, ok, true)
	Equal(t, val[0], "1.1")

	val, ok = values["Float64ArrayArray[0][0]"]
	Equal(t, ok, true)
	Equal(t, val[0], "1.1")

	val, ok = values["Float64ArrayArray[0][1]"]
	Equal(t, ok, true)
	Equal(t, val[0], "0")

	val, ok = values["Float64ArrayArray[0][2]"]
	Equal(t, ok, true)
	Equal(t, val[0], "3.3")

	_, ok = values["Float64ArrayArray[1][0]"]
	Equal(t, ok, false)

	val, ok = values["Float64ArrayArray[2][0]"]
	Equal(t, ok, true)
	Equal(t, val[0], "1.1")

	val, ok = values["Float32PtrArrayArray[0][0]"]
	Equal(t, ok, true)
	Equal(t, val[0], "1.1")

	_, ok = values["Float32PtrArrayArray[0][1]"]
	Equal(t, ok, false)

	val, ok = values["Float32PtrArrayArray[0][2]"]
	Equal(t, ok, true)
	Equal(t, val[0], "3.3")

	_, ok = values["Float32PtrArrayArray[1][0]"]
	Equal(t, ok, false)

	val, ok = values["Float32PtrArrayArray[2][0]"]
	Equal(t, ok, true)
	Equal(t, val[0], "1.1")

	val, ok = values["Float64PtrArrayArray[0][0]"]
	Equal(t, ok, true)
	Equal(t, val[0], "1.1")

	_, ok = values["Float64PtrArrayArray[0][1]"]
	Equal(t, ok, false)

	val, ok = values["Float64PtrArrayArray[0][2]"]
	Equal(t, ok, true)
	Equal(t, val[0], "3.3")

	_, ok = values["Float64PtrArrayArray[1][0]"]
	Equal(t, ok, false)

	val, ok = values["Float64PtrArrayArray[2][0]"]
	Equal(t, ok, true)
	Equal(t, val[0], "1.1")

	val, ok = values["Float32Map[1.1]"]
	Equal(t, ok, true)
	Equal(t, val[0], "3.3")

	val, ok = values["Float32Map[3.3]"]
	Equal(t, ok, true)
	Equal(t, val[0], "2.2")

	val, ok = values["Float64Map[1.1]"]
	Equal(t, ok, true)
	Equal(t, val[0], "3.3")

	val, ok = values["Float64Map[3.3]"]
	Equal(t, ok, true)
	Equal(t, val[0], "2.2")

	val, ok = values["Float32PtrMap[1.1]"]
	Equal(t, ok, true)
	Equal(t, val[0], "3.3")

	val, ok = values["Float32PtrMap[3.3]"]
	Equal(t, ok, true)
	Equal(t, val[0], "2.2")

	val, ok = values["Float64PtrMap[1.1]"]
	Equal(t, ok, true)
	Equal(t, val[0], "3.3")

	val, ok = values["Float64PtrMap[3.3]"]
	Equal(t, ok, true)
	Equal(t, val[0], "2.2")

	val, ok = values["NoValue"]
	Equal(t, ok, true)
	Equal(t, val[0], "0")
}

func TestEncoderBool(t *testing.T) {

	type TestBool struct {
		Bool              bool
		BoolPtr           *bool
		BoolArray         []bool
		BoolPtrArray      []*bool
		BoolArrayArray    [][]bool
		BoolPtrArrayArray [][]*bool
		BoolMap           map[bool]bool
		BoolPtrMap        map[*bool]*bool
		NoValue           bool
	}

	tr := true
	fa := false

	test := TestBool{
		Bool:              tr,
		BoolPtr:           &tr,
		BoolArray:         []bool{fa, tr, tr},
		BoolPtrArray:      []*bool{&fa, nil, &tr},
		BoolArrayArray:    [][]bool{{tr, fa, tr}, nil, {tr}},
		BoolPtrArrayArray: [][]*bool{{&tr, nil, &tr}, nil, {&tr}},
		BoolMap:           map[bool]bool{tr: fa, fa: true},
		BoolPtrMap:        map[*bool]*bool{&tr: &fa, &fa: &tr},
	}

	encoder := NewEncoder()
	values, errs := encoder.Encode(test)

	Equal(t, errs, nil)
	Equal(t, len(values), 17)

	val, ok := values["Bool"]
	Equal(t, ok, true)
	Equal(t, val[0], "true")

	val, ok = values["BoolPtr"]
	Equal(t, ok, true)
	Equal(t, val[0], "true")

	val, ok = values["BoolArray"]
	Equal(t, ok, true)
	Equal(t, len(val), 3)
	Equal(t, val[0], "false")
	Equal(t, val[1], "true")
	Equal(t, val[2], "true")

	val, ok = values["BoolPtrArray[0]"]
	Equal(t, ok, true)
	Equal(t, val[0], "false")

	_, ok = values["BoolPtrArray[1]"]
	Equal(t, ok, false)

	val, ok = values["BoolPtrArray[2]"]
	Equal(t, ok, true)
	Equal(t, val[0], "true")

	val, ok = values["BoolArrayArray[0][0]"]
	Equal(t, ok, true)
	Equal(t, val[0], "true")

	val, ok = values["BoolArrayArray[0][1]"]
	Equal(t, ok, true)
	Equal(t, val[0], "false")

	val, ok = values["BoolArrayArray[0][2]"]
	Equal(t, ok, true)
	Equal(t, val[0], "true")

	_, ok = values["BoolArrayArray[1][0]"]
	Equal(t, ok, false)

	val, ok = values["BoolArrayArray[2][0]"]
	Equal(t, ok, true)
	Equal(t, val[0], "true")

	val, ok = values["BoolPtrArrayArray[0][0]"]
	Equal(t, ok, true)
	Equal(t, val[0], "true")

	_, ok = values["BoolPtrArrayArray[0][1]"]
	Equal(t, ok, false)

	val, ok = values["BoolPtrArrayArray[0][2]"]
	Equal(t, ok, true)
	Equal(t, val[0], "true")

	_, ok = values["BoolPtrArrayArray[1][0]"]
	Equal(t, ok, false)

	val, ok = values["BoolPtrArrayArray[2][0]"]
	Equal(t, ok, true)
	Equal(t, val[0], "true")

	val, ok = values["BoolMap[true]"]
	Equal(t, ok, true)
	Equal(t, val[0], "false")

	val, ok = values["BoolMap[false]"]
	Equal(t, ok, true)
	Equal(t, val[0], "true")

	val, ok = values["BoolPtrMap[true]"]
	Equal(t, ok, true)
	Equal(t, val[0], "false")

	val, ok = values["BoolPtrMap[false]"]
	Equal(t, ok, true)
	Equal(t, val[0], "true")
}

func TestEncoderStruct(t *testing.T) {

	type Phone struct {
		Number string
	}

	type TestStruct struct {
		Name      string `form:"name"`
		Ignore    string `form:"-"`
		NonNilPtr *Phone
		Phone     []Phone
		PhonePtr  []*Phone

		Anonymous struct {
			Value     string
			Ignore    string `form:"-"`
			unexposed string
		}
		Time            time.Time
		TimePtr         *time.Time
		unexposed       string
		Invalid         interface{}
		ExistingMap     map[string]string `form:"mp"`
		MapNoValue      map[int]int
		Array           []string
		ZeroLengthArray []string
		IfaceNonNil     interface{}
		IfaceInvalid    interface{}
		TimeMapKey      map[time.Time]string
		ArrayMap        []map[int]int
		ArrayTime       []time.Time
	}

	tm, err := time.Parse("2006-01-02", "2016-01-02")
	Equal(t, err, nil)

	test := &TestStruct{
		Name:      "joeybloggs",
		Ignore:    "ignore",
		NonNilPtr: new(Phone),
		Phone: []Phone{
			{Number: "1(111)111-1111"},
			{Number: "9(999)999-9999"},
		},
		PhonePtr: []*Phone{
			{Number: "1(111)111-1111"},
			{Number: "9(999)999-9999"},
		},
		Anonymous: struct {
			Value     string
			Ignore    string `form:"-"`
			unexposed string
		}{
			Value: "Anon",
		},
		Time:            tm,
		TimePtr:         &tm,
		unexposed:       "unexposed field",
		ExistingMap:     map[string]string{"existingkey": "existingvalue", "key": "value"},
		Array:           []string{"1", "2"},
		ZeroLengthArray: []string{},
		IfaceNonNil:     new(Phone),
		IfaceInvalid:    nil,
		TimeMapKey:      map[time.Time]string{tm: "time"},
		ArrayMap:        []map[int]int{{1: 3}},
		ArrayTime:       []time.Time{tm},
	}

	encoder := NewEncoder()
	encoder.SetTagName("form")
	encoder.RegisterCustomTypeFunc(func(x interface{}) ([]string, error) {
		return []string{x.(time.Time).Format("2006-01-02")}, nil
	}, time.Time{})

	values, errs := encoder.Encode(test)

	Equal(t, errs, nil)
	Equal(t, len(values), 16)

	val, ok := values["name"]
	Equal(t, ok, true)
	Equal(t, val[0], "joeybloggs")

	_, ok = values["Ignore"]
	Equal(t, ok, false)

	val, ok = values["NonNilPtr.Number"]
	Equal(t, ok, true)
	Equal(t, val[0], "")

	val, ok = values["Phone[0].Number"]
	Equal(t, ok, true)
	Equal(t, val[0], "1(111)111-1111")

	val, ok = values["Phone[1].Number"]
	Equal(t, ok, true)
	Equal(t, val[0], "9(999)999-9999")

	val, ok = values["PhonePtr[0].Number"]
	Equal(t, ok, true)
	Equal(t, val[0], "1(111)111-1111")

	val, ok = values["PhonePtr[1].Number"]
	Equal(t, ok, true)
	Equal(t, val[0], "9(999)999-9999")

	val, ok = values["Anonymous.Value"]
	Equal(t, ok, true)
	Equal(t, val[0], "Anon")

	_, ok = values["Anonymous.Ignore"]
	Equal(t, ok, false)

	_, ok = values["Anonymous.unexposed"]
	Equal(t, ok, false)

	val, ok = values["Time"]
	Equal(t, ok, true)
	Equal(t, val[0], tm.Format("2006-01-02"))

	val, ok = values["TimePtr"]
	Equal(t, ok, true)
	Equal(t, val[0], tm.Format("2006-01-02"))

	_, ok = values["unexposed"]
	Equal(t, ok, false)

	val, ok = values["mp[existingkey]"]
	Equal(t, ok, true)
	Equal(t, val[0], "existingvalue")

	val, ok = values["mp[key]"]
	Equal(t, ok, true)
	Equal(t, val[0], "value")

	val, ok = values["Array"]
	Equal(t, ok, true)
	Equal(t, len(val), 2)
	Equal(t, val[0], "1")
	Equal(t, val[1], "2")

	_, ok = values["ZeroLengthArray"]
	Equal(t, ok, false)

	val, ok = values["IfaceNonNil.Number"]
	Equal(t, ok, true)
	Equal(t, val[0], "")

	_, ok = values["IfaceInvalid"]
	Equal(t, ok, false)

	val, ok = values["IfaceNonNil.Number"]
	Equal(t, ok, true)
	Equal(t, val[0], "")

	val, ok = values["TimeMapKey["+tm.Format("2006-01-02")+"]"]
	Equal(t, ok, true)
	Equal(t, val[0], "time")

	val, ok = values["ArrayMap[0][1]"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["ArrayTime[0]"]
	Equal(t, ok, true)
	Equal(t, val[0], tm.Format("2006-01-02"))

	s := struct {
		Value     string
		Ignore    string `form:"-"`
		unexposed string
		ArrayTime []time.Time
	}{
		Value:     "tval",
		Ignore:    "ignore",
		unexposed: "unexp",
		ArrayTime: []time.Time{tm},
	}

	encoder = NewEncoder()
	values, errs = encoder.Encode(&s)
	Equal(t, errs, nil)
	Equal(t, len(values), 2)

	val, ok = values["Value"]
	Equal(t, ok, true)
	Equal(t, val[0], "tval")

	_, ok = values["Ignore"]
	Equal(t, ok, false)

	_, ok = values["unexposed"]
	Equal(t, ok, false)

	val, ok = values["ArrayTime[0]"]
	Equal(t, ok, true)
	Equal(t, val[0], tm.Format(time.RFC3339))
}

func TestEncoderStructCustomNamespace(t *testing.T) {

	type Phone struct {
		Number string
	}

	type TestStruct struct {
		Name      string `form:"name"`
		Ignore    string `form:"-"`
		NonNilPtr *Phone
		Phone     []Phone
		PhonePtr  []*Phone

		Anonymous struct {
			Value     string
			Ignore    string `form:"-"`
			unexposed string
		}
		Time            time.Time
		TimePtr         *time.Time
		unexposed       string
		Invalid         interface{}
		ExistingMap     map[string]string `form:"mp"`
		MapNoValue      map[int]int
		Array           []string
		ZeroLengthArray []string
		IfaceNonNil     interface{}
		IfaceInvalid    interface{}
		TimeMapKey      map[time.Time]string
		ArrayMap        []map[int]int
		ArrayTime       []time.Time
	}

	tm, err := time.Parse("2006-01-02", "2016-01-02")
	Equal(t, err, nil)

	test := &TestStruct{
		Name:      "joeybloggs",
		Ignore:    "ignore",
		NonNilPtr: new(Phone),
		Phone: []Phone{
			{Number: "1(111)111-1111"},
			{Number: "9(999)999-9999"},
		},
		PhonePtr: []*Phone{
			{Number: "1(111)111-1111"},
			{Number: "9(999)999-9999"},
		},
		Anonymous: struct {
			Value     string
			Ignore    string `form:"-"`
			unexposed string
		}{
			Value: "Anon",
		},
		Time:            tm,
		TimePtr:         &tm,
		unexposed:       "unexposed field",
		ExistingMap:     map[string]string{"existingkey": "existingvalue", "key": "value"},
		Array:           []string{"1", "2"},
		ZeroLengthArray: []string{},
		IfaceNonNil:     new(Phone),
		IfaceInvalid:    nil,
		TimeMapKey:      map[time.Time]string{tm: "time"},
		ArrayMap:        []map[int]int{{1: 3}},
		ArrayTime:       []time.Time{tm},
	}

	encoder := NewEncoder()
	encoder.SetTagName("form")
	encoder.RegisterCustomTypeFunc(func(x interface{}) ([]string, error) {
		return []string{x.(time.Time).Format("2006-01-02")}, nil
	}, time.Time{})
	encoder.SetNamespacePrefix("[")
	encoder.SetNamespaceSuffix("]")

	values, errs := encoder.Encode(test)

	Equal(t, errs, nil)
	Equal(t, len(values), 16)

	val, ok := values["name"]
	Equal(t, ok, true)
	Equal(t, val[0], "joeybloggs")

	_, ok = values["Ignore"]
	Equal(t, ok, false)

	val, ok = values["NonNilPtr[Number]"]
	Equal(t, ok, true)
	Equal(t, val[0], "")

	val, ok = values["Phone[0][Number]"]
	Equal(t, ok, true)
	Equal(t, val[0], "1(111)111-1111")

	val, ok = values["Phone[1][Number]"]
	Equal(t, ok, true)
	Equal(t, val[0], "9(999)999-9999")

	val, ok = values["PhonePtr[0][Number]"]
	Equal(t, ok, true)
	Equal(t, val[0], "1(111)111-1111")

	val, ok = values["PhonePtr[1][Number]"]
	Equal(t, ok, true)
	Equal(t, val[0], "9(999)999-9999")

	val, ok = values["Anonymous[Value]"]
	Equal(t, ok, true)
	Equal(t, val[0], "Anon")

	_, ok = values["Anonymous[Ignore]"]
	Equal(t, ok, false)

	_, ok = values["Anonymous.unexposed"]
	Equal(t, ok, false)

	val, ok = values["Time"]
	Equal(t, ok, true)
	Equal(t, val[0], tm.Format("2006-01-02"))

	val, ok = values["TimePtr"]
	Equal(t, ok, true)
	Equal(t, val[0], tm.Format("2006-01-02"))

	_, ok = values["unexposed"]
	Equal(t, ok, false)

	val, ok = values["mp[existingkey]"]
	Equal(t, ok, true)
	Equal(t, val[0], "existingvalue")

	val, ok = values["mp[key]"]
	Equal(t, ok, true)
	Equal(t, val[0], "value")

	val, ok = values["Array"]
	Equal(t, ok, true)
	Equal(t, len(val), 2)
	Equal(t, val[0], "1")
	Equal(t, val[1], "2")

	_, ok = values["ZeroLengthArray"]
	Equal(t, ok, false)

	val, ok = values["IfaceNonNil[Number]"]
	Equal(t, ok, true)
	Equal(t, val[0], "")

	_, ok = values["IfaceInvalid"]
	Equal(t, ok, false)

	val, ok = values["IfaceNonNil[Number]"]
	Equal(t, ok, true)
	Equal(t, val[0], "")

	val, ok = values["TimeMapKey["+tm.Format("2006-01-02")+"]"]
	Equal(t, ok, true)
	Equal(t, val[0], "time")

	val, ok = values["ArrayMap[0][1]"]
	Equal(t, ok, true)
	Equal(t, val[0], "3")

	val, ok = values["ArrayTime[0]"]
	Equal(t, ok, true)
	Equal(t, val[0], tm.Format("2006-01-02"))

	s := struct {
		Value     string
		Ignore    string `form:"-"`
		unexposed string
		ArrayTime []time.Time
	}{
		Value:     "tval",
		Ignore:    "ignore",
		unexposed: "unexp",
		ArrayTime: []time.Time{tm},
	}

	encoder = NewEncoder()
	values, errs = encoder.Encode(&s)
	Equal(t, errs, nil)
	Equal(t, len(values), 2)

	val, ok = values["Value"]
	Equal(t, ok, true)
	Equal(t, val[0], "tval")

	_, ok = values["Ignore"]
	Equal(t, ok, false)

	_, ok = values["unexposed"]
	Equal(t, ok, false)

	val, ok = values["ArrayTime[0]"]
	Equal(t, ok, true)
	Equal(t, val[0], tm.Format(time.RFC3339))
}

func TestEncoderMap(t *testing.T) {
	inner := map[string]interface{}{
		"inner": "1",
	}

	outer := map[string]interface{}{
		"outer": inner,
	}

	enc := NewEncoder()
	values, _ := enc.Encode(outer)

	val, ok := values["[outer][inner]"]
	Equal(t, ok, true)
	Equal(t, val[0], "1")
}

func TestDecodeAllNonStructTypes(t *testing.T) {

	encoder := NewEncoder()

	// test integers

	i := int(3)
	i8 := int8(2)
	i16 := int16(1)
	i32 := int32(0)
	i64 := int64(3)

	values, err := encoder.Encode(i)
	Equal(t, err, nil)
	Equal(t, values[""][0], "3")

	values, err = encoder.Encode(i8)
	Equal(t, err, nil)
	Equal(t, values[""][0], "2")

	values, err = encoder.Encode(i16)
	Equal(t, err, nil)
	Equal(t, values[""][0], "1")

	values, err = encoder.Encode(i32)
	Equal(t, err, nil)
	Equal(t, values[""][0], "0")

	values, err = encoder.Encode(i64)
	Equal(t, err, nil)
	Equal(t, values[""][0], "3")

	// test unsigned integers

	ui := uint(3)
	ui8 := uint8(2)
	ui16 := uint16(1)
	ui32 := uint32(0)
	ui64 := uint64(3)

	values, err = encoder.Encode(ui)
	Equal(t, err, nil)
	Equal(t, values[""][0], "3")

	values, err = encoder.Encode(ui8)
	Equal(t, err, nil)
	Equal(t, values[""][0], "2")

	values, err = encoder.Encode(ui16)
	Equal(t, err, nil)
	Equal(t, values[""][0], "1")

	values, err = encoder.Encode(ui32)
	Equal(t, err, nil)
	Equal(t, values[""][0], "0")

	values, err = encoder.Encode(ui64)
	Equal(t, err, nil)
	Equal(t, values[""][0], "3")

	// test bool

	ok := true
	values, err = encoder.Encode(ok)
	Equal(t, err, nil)
	Equal(t, values[""][0], "true")

	ok = false
	values, err = encoder.Encode(ok)
	Equal(t, err, nil)
	Equal(t, values[""][0], "false")

	// test time

	tm, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	Equal(t, err, nil)

	values, err = encoder.Encode(tm)
	Equal(t, err, nil)
	Equal(t, values[""][0], "2006-01-02T15:04:05Z")

	// test basic array

	arr := []string{"arr1", "arr2"}

	values, err = encoder.Encode(arr)
	Equal(t, err, nil)
	Equal(t, len(values), 1)
	Equal(t, values[""][0], "arr1")
	Equal(t, values[""][1], "arr2")

	// test ptr array

	s1 := "arr1"
	s2 := "arr2"
	arrPtr := []*string{&s1, &s2}

	values, err = encoder.Encode(arrPtr)
	Equal(t, err, nil)
	Equal(t, len(values), 2)
	Equal(t, values["[0]"][0], "arr1")
	Equal(t, values["[1]"][0], "arr2")

	// test map

	m := map[string]string{"key1": "val1", "key2": "val2"}

	values, err = encoder.Encode(m)
	Equal(t, err, nil)
	Equal(t, len(values), 2)
	Equal(t, values["[key1]"][0], "val1")
	Equal(t, values["[key2]"][0], "val2")
}

func TestEncoderNativeTime(t *testing.T) {

	type TestError struct {
		Time        time.Time
		TimeNoValue time.Time
	}

	tm, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	Equal(t, err, nil)

	test := TestError{
		Time: tm,
	}

	encoder := NewEncoder()
	values, errs := encoder.Encode(&test)
	Equal(t, errs, nil)

	val, ok := values["Time"]
	Equal(t, ok, true)
	Equal(t, val[0], tm.Format(time.RFC3339))

	val, ok = values["TimeNoValue"]
	Equal(t, ok, true)
	Equal(t, val[0], "0001-01-01T00:00:00Z")
}

func TestEncoderErrors(t *testing.T) {

	tm, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	Equal(t, err, nil)

	type TestError struct {
		Time      time.Time
		BadMapKey map[time.Time]string
		Iface     map[interface{}]string
		Struct    map[struct{}]string
	}

	test := TestError{
		BadMapKey: map[time.Time]string{tm: "time"},
		Iface:     map[interface{}]string{nil: "time"},
		Struct:    map[struct{}]string{{}: "str"},
	}

	encoder := NewEncoder()
	encoder.RegisterCustomTypeFunc(func(x interface{}) ([]string, error) {
		return nil, errors.New("Bad Type Conversion")
	}, time.Time{})

	values, errs := encoder.Encode(&test)
	NotEqual(t, errs, nil)

	Equal(t, len(values), 0)

	e := errs.Error()
	NotEqual(t, e, "")

	ee := errs.(EncodeErrors)
	Equal(t, len(ee), 3)

	k := ee["Time"]
	Equal(t, k.Error(), "Bad Type Conversion")

	k = ee["BadMapKey"]
	Equal(t, k.Error(), "Bad Type Conversion")

	k = ee["Struct"]
	Equal(t, k.Error(), "Unsupported Map Key '<struct {} Value>' Namespace 'Struct'")
}

func TestEncoderPanicsAndBadValues(t *testing.T) {

	encoder := NewEncoder()

	values, err := encoder.Encode(nil)
	NotEqual(t, err, nil)
	Equal(t, values, nil)

	_, ok := err.(*InvalidEncodeError)
	Equal(t, ok, true)
	Equal(t, err.Error(), "form: Encode(nil)")

	type TestStruct struct {
		Value string
	}

	var tst *TestStruct

	values, err = encoder.Encode(tst)
	NotEqual(t, err, nil)
	Equal(t, values, nil)

	_, ok = err.(*InvalidEncodeError)
	Equal(t, ok, true)
	Equal(t, err.Error(), "form: Encode(nil *form.TestStruct)")
}

func TestEncoderExplicit(t *testing.T) {

	type Test struct {
		Name string `form:"Name"`
		Age  int
	}

	test := &Test{
		Name: "Joeybloggs",
		Age:  3,
	}

	encoder := NewEncoder()
	encoder.SetMode(ModeExplicit)

	values, err := encoder.Encode(test)
	Equal(t, err, nil)
	Equal(t, len(values), 1)
	Equal(t, values["Name"][0], "Joeybloggs")
}

func TestEncoderRegisterTagNameFunc(t *testing.T) {

	type Test struct {
		Name string `json:"name"`
		Age  int    `json:"-"`
	}

	test := &Test{
		Name: "Joeybloggs",
		Age:  3,
	}

	encoder := NewEncoder()
	encoder.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("json")

		if commaIndex := strings.Index(name, ","); commaIndex != -1 {
			name = name[:commaIndex]
		}

		return name
	})

	values, err := encoder.Encode(test)
	Equal(t, err, nil)
	Equal(t, len(values), 1)
	Equal(t, values["name"][0], "Joeybloggs")
}

func TestEncoderEmbedModes(t *testing.T) {

	type A struct {
		Field string
	}

	type B struct {
		A
		Field string
	}

	b := B{
		A: A{
			Field: "A Val",
		},
		Field: "B Val",
	}

	encoder := NewEncoder()

	values, err := encoder.Encode(b)
	Equal(t, err, nil)
	Equal(t, len(values), 1)
	Equal(t, values["Field"][0], "B Val")
	Equal(t, values["Field"][1], "A Val")

	encoder.SetAnonymousMode(AnonymousSeparate)
	values, err = encoder.Encode(b)
	Equal(t, err, nil)
	Equal(t, len(values), 2)
	Equal(t, values["Field"][0], "B Val")
	Equal(t, values["A.Field"][0], "A Val")
}

func TestOmitEmpty(t *testing.T) {

	type NotComparable struct {
		Slice []string
	}

	type Test struct {
		String  string            `form:",omitempty"`
		Array   []string          `form:",omitempty"`
		Map     map[string]string `form:",omitempty"`
		String2 string            `form:"str,omitempty"`
		Array2  []string          `form:"arr,omitempty"`
		Map2    map[string]string `form:"map,omitempty"`
		NotComparable			  `form:",omitempty"`
	}

	var tst Test

	encoder := NewEncoder()

	values, err := encoder.Encode(tst)
	Equal(t, err, nil)
	Equal(t, len(values), 0)

	type Test2 struct {
		String  string
		Array   []string
		Map     map[string]string
		String2 string            `form:"str,omitempty"`
		Array2  []string          `form:"arr,omitempty"`
		Map2    map[string]string `form:"map,omitempty"`
	}

	var tst2 Test2

	values, err = encoder.Encode(tst2)
	Equal(t, err, nil)
	Equal(t, len(values), 1)
	Equal(t, values["String"][0], "")

	type Test3 struct {
		String  string
		Array   []string
		Map     map[string]string
		String2 string `form:"str"`
		Array2  []string
		Map2    map[string]string
	}

	var tst3 Test3

	values, err = encoder.Encode(tst3)
	Equal(t, err, nil)
	Equal(t, len(values), 2)
	Equal(t, values["String"][0], "")
	Equal(t, values["str"][0], "")

	type T struct {
		X      *uint8    `form:"x,omitempty"`
		Array  []*string `form:"arr,omitempty"`
		Array2 []*string `form:"arr2,dive,omitempty"`
	}
	x := uint8(0)
	s := ""
	tst4 := T{
		X:     &x,
		Array: []*string{&s},
	}

	values, err = encoder.Encode(tst4)
	Equal(t, err, nil)
	Equal(t, len(values), 2)
	Equal(t, values["x"][0], "0")
	Equal(t, values["arr[0]"][0], "")
}

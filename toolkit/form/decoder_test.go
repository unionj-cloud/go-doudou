package form

import (
	"errors"
	"fmt"
	"net/url"
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

func TestDecoderInt(t *testing.T) {

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
		NoURLValue       int
		IntNoValues      int
		Int8NoValues     int8
		Int16NoValues    int16
		Int32NoValues    int32
		Int64NoValues    int64
	}

	values := url.Values{
		"Int":                    []string{"3"},
		"Int8":                   []string{"3"},
		"Int16":                  []string{"3"},
		"Int32":                  []string{"3"},
		"Int64":                  []string{"3"},
		"IntPtr":                 []string{"3"},
		"Int8Ptr":                []string{"3"},
		"Int16Ptr":               []string{"3"},
		"Int32Ptr":               []string{"3"},
		"Int64Ptr":               []string{"3"},
		"IntArray":               []string{"1", "2", "3"},
		"IntPtrArray[0]":         []string{"1"},
		"IntPtrArray[2]":         []string{"3"},
		"IntArrayArray[0][0]":    []string{"1"},
		"IntArrayArray[0][2]":    []string{"3"},
		"IntArrayArray[2][0]":    []string{"1"},
		"IntPtrArrayArray[0][0]": []string{"1"},
		"IntPtrArrayArray[0][2]": []string{"3"},
		"IntPtrArrayArray[2][0]": []string{"1"},
		"IntMap[1]":              []string{"3"},
		"IntPtrMap[1]":           []string{"3"},
	}

	var test TestInt

	test.IntArray = make([]int, 4)

	decoder := NewDecoder()
	errs := decoder.Decode(&test, values)
	Equal(t, errs, nil)

	Equal(t, test.Int, int(3))
	Equal(t, test.Int8, int8(3))
	Equal(t, test.Int16, int16(3))
	Equal(t, test.Int32, int32(3))
	Equal(t, test.Int64, int64(3))

	Equal(t, *test.IntPtr, int(3))
	Equal(t, *test.Int8Ptr, int8(3))
	Equal(t, *test.Int16Ptr, int16(3))
	Equal(t, *test.Int32Ptr, int32(3))
	Equal(t, *test.Int64Ptr, int64(3))

	Equal(t, len(test.IntArray), 7)
	Equal(t, test.IntArray[0], int(0))
	Equal(t, test.IntArray[1], int(0))
	Equal(t, test.IntArray[2], int(0))
	Equal(t, test.IntArray[3], int(0))
	Equal(t, test.IntArray[4], int(1))
	Equal(t, test.IntArray[5], int(2))
	Equal(t, test.IntArray[6], int(3))

	Equal(t, len(test.IntPtrArray), 3)
	Equal(t, *test.IntPtrArray[0], int(1))
	Equal(t, test.IntPtrArray[1], nil)
	Equal(t, *test.IntPtrArray[2], int(3))

	Equal(t, len(test.IntArrayArray), 3)
	Equal(t, len(test.IntArrayArray[0]), 3)
	Equal(t, len(test.IntArrayArray[1]), 0)
	Equal(t, len(test.IntArrayArray[2]), 1)
	Equal(t, test.IntArrayArray[0][0], int(1))
	Equal(t, test.IntArrayArray[0][1], int(0))
	Equal(t, test.IntArrayArray[0][2], int(3))
	Equal(t, test.IntArrayArray[2][0], int(1))

	Equal(t, len(test.IntPtrArrayArray), 3)
	Equal(t, len(test.IntPtrArrayArray[0]), 3)
	Equal(t, len(test.IntPtrArrayArray[1]), 0)
	Equal(t, len(test.IntPtrArrayArray[2]), 1)
	Equal(t, *test.IntPtrArrayArray[0][0], int(1))
	Equal(t, test.IntPtrArrayArray[0][1], nil)
	Equal(t, *test.IntPtrArrayArray[0][2], int(3))
	Equal(t, *test.IntPtrArrayArray[2][0], int(1))

	Equal(t, len(test.IntMap), 1)
	Equal(t, len(test.IntPtrMap), 1)

	v, ok := test.IntMap[1]
	Equal(t, ok, true)
	Equal(t, v, int(3))

	Equal(t, test.NoURLValue, int(0))

	Equal(t, test.IntNoValues, int(0))
	Equal(t, test.Int8NoValues, int8(0))
	Equal(t, test.Int16NoValues, int16(0))
	Equal(t, test.Int32NoValues, int32(0))
	Equal(t, test.Int64NoValues, int64(0))
}

func TestDecoderUint(t *testing.T) {

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
		NoURLValue        uint
		UintNoValues      uint
		Uint8NoValues     uint8
		Uint16NoValues    uint16
		Uint32NoValues    uint32
		Uint64NoValues    uint64
	}

	values := url.Values{
		"Uint":                    []string{"3"},
		"Uint8":                   []string{"3"},
		"Uint16":                  []string{"3"},
		"Uint32":                  []string{"3"},
		"Uint64":                  []string{"3"},
		"UintPtr":                 []string{"3"},
		"Uint8Ptr":                []string{"3"},
		"Uint16Ptr":               []string{"3"},
		"Uint32Ptr":               []string{"3"},
		"Uint64Ptr":               []string{"3"},
		"UintArray":               []string{"1", "2", "3"},
		"UintPtrArray[0]":         []string{"1"},
		"UintPtrArray[2]":         []string{"3"},
		"UintArrayArray[0][0]":    []string{"1"},
		"UintArrayArray[0][2]":    []string{"3"},
		"UintArrayArray[2][0]":    []string{"1"},
		"UintPtrArrayArray[0][0]": []string{"1"},
		"UintPtrArrayArray[0][2]": []string{"3"},
		"UintPtrArrayArray[2][0]": []string{"1"},
		"UintMap[1]":              []string{"3"},
		"UintPtrMap[1]":           []string{"3"},
	}

	var test TestUint

	test.UintArray = make([]uint, 4)

	decoder := NewDecoder()
	errs := decoder.Decode(&test, values)
	Equal(t, errs, nil)

	Equal(t, test.Uint, uint(3))
	Equal(t, test.Uint8, uint8(3))
	Equal(t, test.Uint16, uint16(3))
	Equal(t, test.Uint32, uint32(3))
	Equal(t, test.Uint64, uint64(3))

	Equal(t, *test.UintPtr, uint(3))
	Equal(t, *test.Uint8Ptr, uint8(3))
	Equal(t, *test.Uint16Ptr, uint16(3))
	Equal(t, *test.Uint32Ptr, uint32(3))
	Equal(t, *test.Uint64Ptr, uint64(3))

	Equal(t, len(test.UintArray), 7)
	Equal(t, test.UintArray[0], uint(0))
	Equal(t, test.UintArray[1], uint(0))
	Equal(t, test.UintArray[2], uint(0))
	Equal(t, test.UintArray[3], uint(0))
	Equal(t, test.UintArray[4], uint(1))
	Equal(t, test.UintArray[5], uint(2))
	Equal(t, test.UintArray[6], uint(3))

	Equal(t, len(test.UintPtrArray), 3)
	Equal(t, *test.UintPtrArray[0], uint(1))
	Equal(t, test.UintPtrArray[1], nil)
	Equal(t, *test.UintPtrArray[2], uint(3))

	Equal(t, len(test.UintArrayArray), 3)
	Equal(t, len(test.UintArrayArray[0]), 3)
	Equal(t, len(test.UintArrayArray[1]), 0)
	Equal(t, len(test.UintArrayArray[2]), 1)
	Equal(t, test.UintArrayArray[0][0], uint(1))
	Equal(t, test.UintArrayArray[0][1], uint(0))
	Equal(t, test.UintArrayArray[0][2], uint(3))
	Equal(t, test.UintArrayArray[2][0], uint(1))

	Equal(t, len(test.UintPtrArrayArray), 3)
	Equal(t, len(test.UintPtrArrayArray[0]), 3)
	Equal(t, len(test.UintPtrArrayArray[1]), 0)
	Equal(t, len(test.UintPtrArrayArray[2]), 1)
	Equal(t, *test.UintPtrArrayArray[0][0], uint(1))
	Equal(t, test.UintPtrArrayArray[0][1], nil)
	Equal(t, *test.UintPtrArrayArray[0][2], uint(3))
	Equal(t, *test.UintPtrArrayArray[2][0], uint(1))

	Equal(t, len(test.UintMap), 1)
	Equal(t, len(test.UintPtrMap), 1)

	v, ok := test.UintMap[1]
	Equal(t, ok, true)
	Equal(t, v, uint(3))

	Equal(t, test.NoURLValue, uint(0))

	Equal(t, test.UintNoValues, uint(0))
	Equal(t, test.Uint8NoValues, uint8(0))
	Equal(t, test.Uint16NoValues, uint16(0))
	Equal(t, test.Uint32NoValues, uint32(0))
	Equal(t, test.Uint64NoValues, uint64(0))
}

func TestDecoderString(t *testing.T) {

	type TestString struct {
		String              string
		StringPtr           *string
		StringArray         []string
		StringPtrArray      []*string
		StringArrayArray    [][]string
		StringPtrArrayArray [][]*string
		StringMap           map[string]string
		StringPtrMap        map[*string]*string
		NoURLValue          string
	}

	values := url.Values{
		"String":                    []string{"3"},
		"StringPtr":                 []string{"3"},
		"StringArray":               []string{"1", "2", "3"},
		"StringPtrArray[0]":         []string{"1"},
		"StringPtrArray[2]":         []string{"3"},
		"StringArrayArray[0][0]":    []string{"1"},
		"StringArrayArray[0][2]":    []string{"3"},
		"StringArrayArray[2][0]":    []string{"1"},
		"StringPtrArrayArray[0][0]": []string{"1"},
		"StringPtrArrayArray[0][2]": []string{"3"},
		"StringPtrArrayArray[2][0]": []string{"1"},
		"StringMap[1]":              []string{"3"},
		"StringPtrMap[1]":           []string{"3"},
	}

	var test TestString

	test.StringArray = make([]string, 4)

	decoder := NewDecoder()
	errs := decoder.Decode(&test, values)
	Equal(t, errs, nil)

	Equal(t, test.String, "3")

	Equal(t, *test.StringPtr, "3")

	Equal(t, len(test.StringArray), 7)
	Equal(t, test.StringArray[0], "")
	Equal(t, test.StringArray[1], "")
	Equal(t, test.StringArray[2], "")
	Equal(t, test.StringArray[3], "")
	Equal(t, test.StringArray[4], "1")
	Equal(t, test.StringArray[5], "2")
	Equal(t, test.StringArray[6], "3")

	Equal(t, len(test.StringPtrArray), 3)
	Equal(t, *test.StringPtrArray[0], "1")
	Equal(t, test.StringPtrArray[1], nil)
	Equal(t, *test.StringPtrArray[2], "3")

	Equal(t, len(test.StringArrayArray), 3)
	Equal(t, len(test.StringArrayArray[0]), 3)
	Equal(t, len(test.StringArrayArray[1]), 0)
	Equal(t, len(test.StringArrayArray[2]), 1)
	Equal(t, test.StringArrayArray[0][0], "1")
	Equal(t, test.StringArrayArray[0][1], "")
	Equal(t, test.StringArrayArray[0][2], "3")
	Equal(t, test.StringArrayArray[2][0], "1")

	Equal(t, len(test.StringPtrArrayArray), 3)
	Equal(t, len(test.StringPtrArrayArray[0]), 3)
	Equal(t, len(test.StringPtrArrayArray[1]), 0)
	Equal(t, len(test.StringPtrArrayArray[2]), 1)
	Equal(t, *test.StringPtrArrayArray[0][0], "1")
	Equal(t, test.StringPtrArrayArray[0][1], nil)
	Equal(t, *test.StringPtrArrayArray[0][2], "3")
	Equal(t, *test.StringPtrArrayArray[2][0], "1")

	Equal(t, len(test.StringMap), 1)
	Equal(t, len(test.StringPtrMap), 1)

	v, ok := test.StringMap["1"]
	Equal(t, ok, true)
	Equal(t, v, "3")

	Equal(t, test.NoURLValue, "")
}

func TestDecoderFloat(t *testing.T) {

	type TestFloat struct {
		Float32              float32
		Float32Ptr           *float32
		Float64              float64
		Float64Ptr           *float64
		Float32Array         []float32
		Float32PtrArray      []*float32
		Float32ArrayArray    [][]float32
		Float32PtrArrayArray [][]*float32
		Float32Map           map[float32]float32
		Float32PtrMap        map[*float32]*float32
		Float32NoValue       float32
		Float64NoValue       float64
	}

	values := url.Values{
		"Float32":                    []string{"3.3"},
		"Float32Ptr":                 []string{"3.3"},
		"Float64":                    []string{"3.3"},
		"Float64Ptr":                 []string{"3.3"},
		"Float32Array":               []string{"1.1", "2.2", "3.3"},
		"Float32PtrArray[0]":         []string{"1.1"},
		"Float32PtrArray[2]":         []string{"3.3"},
		"Float32ArrayArray[0][0]":    []string{"1.1"},
		"Float32ArrayArray[0][2]":    []string{"3.3"},
		"Float32ArrayArray[2][0]":    []string{"1.1"},
		"Float32PtrArrayArray[0][0]": []string{"1.1"},
		"Float32PtrArrayArray[0][2]": []string{"3.3"},
		"Float32PtrArrayArray[2][0]": []string{"1.1"},
		"Float32Map[1.1]":            []string{"3.3"},
		"Float32PtrMap[1.1]":         []string{"3.3"},
	}

	var test TestFloat

	test.Float32Array = make([]float32, 4)

	decoder := NewDecoder()
	errs := decoder.Decode(&test, values)
	Equal(t, errs, nil)

	Equal(t, test.Float32, float32(3.3))
	Equal(t, test.Float64, float64(3.3))

	Equal(t, *test.Float32Ptr, float32(3.3))
	Equal(t, *test.Float64Ptr, float64(3.3))

	Equal(t, len(test.Float32Array), 7)
	Equal(t, test.Float32Array[0], float32(0.0))
	Equal(t, test.Float32Array[1], float32(0.0))
	Equal(t, test.Float32Array[2], float32(0.0))
	Equal(t, test.Float32Array[3], float32(0.0))
	Equal(t, test.Float32Array[4], float32(1.1))
	Equal(t, test.Float32Array[5], float32(2.2))
	Equal(t, test.Float32Array[6], float32(3.3))

	Equal(t, len(test.Float32PtrArray), 3)
	Equal(t, *test.Float32PtrArray[0], float32(1.1))
	Equal(t, test.Float32PtrArray[1], nil)
	Equal(t, *test.Float32PtrArray[2], float32(3.3))

	Equal(t, len(test.Float32ArrayArray), 3)
	Equal(t, len(test.Float32ArrayArray[0]), 3)
	Equal(t, len(test.Float32ArrayArray[1]), 0)
	Equal(t, len(test.Float32ArrayArray[2]), 1)
	Equal(t, test.Float32ArrayArray[0][0], float32(1.1))
	Equal(t, test.Float32ArrayArray[0][1], float32(0.0))
	Equal(t, test.Float32ArrayArray[0][2], float32(3.3))
	Equal(t, test.Float32ArrayArray[2][0], float32(1.1))

	Equal(t, len(test.Float32PtrArrayArray), 3)
	Equal(t, len(test.Float32PtrArrayArray[0]), 3)
	Equal(t, len(test.Float32PtrArrayArray[1]), 0)
	Equal(t, len(test.Float32PtrArrayArray[2]), 1)
	Equal(t, *test.Float32PtrArrayArray[0][0], float32(1.1))
	Equal(t, test.Float32PtrArrayArray[0][1], nil)
	Equal(t, *test.Float32PtrArrayArray[0][2], float32(3.3))
	Equal(t, *test.Float32PtrArrayArray[2][0], float32(1.1))

	Equal(t, len(test.Float32Map), 1)
	Equal(t, len(test.Float32PtrMap), 1)

	v, ok := test.Float32Map[float32(1.1)]
	Equal(t, ok, true)
	Equal(t, v, float32(3.3))

	Equal(t, test.Float32NoValue, float32(0.0))
	Equal(t, test.Float64NoValue, float64(0.0))
}

func TestDecoderBool(t *testing.T) {

	type TestBool struct {
		Bool              bool
		BoolPtr           *bool
		BoolPtrNil        *bool
		BoolPtrEmpty      *bool
		BoolArray         []bool
		BoolPtrArray      []*bool
		BoolArrayArray    [][]bool
		BoolPtrArrayArray [][]*bool
		BoolMap           map[bool]bool
		BoolPtrMap        map[*bool]*bool
		NoURLValue        bool
	}

	values := url.Values{
		"Bool":                    []string{"true"},
		"BoolPtr":                 []string{"true"},
		"BoolPtrEmpty":            []string{""},
		"BoolArray":               []string{"off", "t", "on"},
		"BoolPtrArray[0]":         []string{"true"},
		"BoolPtrArray[2]":         []string{"T"},
		"BoolArrayArray[0][0]":    []string{"TRUE"},
		"BoolArrayArray[0][2]":    []string{"True"},
		"BoolArrayArray[2][0]":    []string{"true"},
		"BoolPtrArrayArray[0][0]": []string{"true"},
		"BoolPtrArrayArray[0][2]": []string{"t"},
		"BoolPtrArrayArray[2][0]": []string{"1"},
		"BoolMap[true]":           []string{"true"},
		"BoolPtrMap[t]":           []string{"true"},
	}

	var test TestBool

	test.BoolArray = make([]bool, 4)

	decoder := NewDecoder()
	errs := decoder.Decode(&test, values)
	Equal(t, errs, nil)

	Equal(t, test.Bool, true)

	Equal(t, *test.BoolPtr, true)
	Equal(t, test.BoolPtrNil, nil)
	NotEqual(t, test.BoolPtrEmpty, nil)
	Equal(t, *test.BoolPtrEmpty, false)

	Equal(t, len(test.BoolArray), 7)
	Equal(t, test.BoolArray[0], false)
	Equal(t, test.BoolArray[1], false)
	Equal(t, test.BoolArray[2], false)
	Equal(t, test.BoolArray[3], false)
	Equal(t, test.BoolArray[4], false)
	Equal(t, test.BoolArray[5], true)
	Equal(t, test.BoolArray[6], true)

	Equal(t, len(test.BoolPtrArray), 3)
	Equal(t, *test.BoolPtrArray[0], true)
	Equal(t, test.BoolPtrArray[1], nil)
	Equal(t, *test.BoolPtrArray[2], true)

	Equal(t, len(test.BoolArrayArray), 3)
	Equal(t, len(test.BoolArrayArray[0]), 3)
	Equal(t, len(test.BoolArrayArray[1]), 0)
	Equal(t, len(test.BoolArrayArray[2]), 1)
	Equal(t, test.BoolArrayArray[0][0], true)
	Equal(t, test.BoolArrayArray[0][1], false)
	Equal(t, test.BoolArrayArray[0][2], true)
	Equal(t, test.BoolArrayArray[2][0], true)

	Equal(t, len(test.BoolPtrArrayArray), 3)
	Equal(t, len(test.BoolPtrArrayArray[0]), 3)
	Equal(t, len(test.BoolPtrArrayArray[1]), 0)
	Equal(t, len(test.BoolPtrArrayArray[2]), 1)
	Equal(t, *test.BoolPtrArrayArray[0][0], true)
	Equal(t, test.BoolPtrArrayArray[0][1], nil)
	Equal(t, *test.BoolPtrArrayArray[0][2], true)
	Equal(t, *test.BoolPtrArrayArray[2][0], true)

	Equal(t, len(test.BoolMap), 1)
	Equal(t, len(test.BoolPtrMap), 1)

	v, ok := test.BoolMap[true]
	Equal(t, ok, true)
	Equal(t, v, true)

	Equal(t, test.NoURLValue, false)
}

func TestDecoderEqualStructMapValue(t *testing.T) {
	type PhoneStruct struct {
		Number string
	}

	type PhoneMap map[string]string

	type TestStruct struct {
		PhoneStruct PhoneStruct `form:"Phone"`
		PhoneMap    PhoneMap    `form:"Phone"`
	}

	testCases := []struct {
		NamespacePrefix string
		NamespaceSuffix string
		Values          url.Values
	}{{
		NamespacePrefix: ".",
		Values: url.Values{
			"Phone.Number":  []string{"111"},
			"Phone[Number]": []string{"222"},
		},
	}, {
		NamespacePrefix: "[",
		NamespaceSuffix: "]",
		Values: url.Values{
			"Phone[Number]": []string{"111"},
		},
	}}

	for _, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("Namespace_%s%s", tc.NamespacePrefix, tc.NamespaceSuffix), func(t *testing.T) {
			decoder := NewDecoder()
			decoder.SetNamespacePrefix(tc.NamespacePrefix)
			decoder.SetNamespaceSuffix(tc.NamespaceSuffix)

			var test TestStruct

			err := decoder.Decode(&test, tc.Values)
			Equal(t, err, nil)

			Equal(t, test.PhoneStruct.Number, "111")

			if tc.NamespacePrefix == "." {
				Equal(t, test.PhoneMap["Number"], "222")
			} else {
				Equal(t, test.PhoneMap["Number"], "111")
			}
		})
	}
}

func TestDecoderStruct(t *testing.T) {

	type Phone struct {
		Number string
	}

	type TestMapKeys struct {
		MapIfaceKey   map[interface{}]string
		MapFloat32Key map[float32]float32
		MapFloat64Key map[float64]float64
		MapNestedInt  map[int]map[int]int
		MapInt8       map[int8]int8
		MapInt16      map[int16]int16
		MapInt32      map[int32]int32
		MapUint8      map[uint8]uint8
		MapUint16     map[uint16]uint16
		MapUint32     map[uint32]uint32
	}

	type TestStruct struct {
		Name      string `form:"name"`
		Phone     []Phone
		PhonePtr  []*Phone
		NonNilPtr *Phone
		Ignore    string `form:"-"`
		Anonymous struct {
			Value     string
			Ignore    string `form:"-"`
			unexposed string
		}
		Time                       time.Time
		TimePtr                    *time.Time
		Invalid                    interface{}
		ExistingMap                map[string]string `form:"mp"`
		MapNoValue                 map[int]int
		TestMapKeys                TestMapKeys
		NilArray                   []string
		TooSmallArray              []string
		TooSmallCapOKArray         []string
		ZeroLengthArray            []string
		TooSmallNumberedArray      []string
		TooSmallCapOKNumberedArray []string
		BigEnoughNumberedArray     []string
		IfaceNonNil                interface{}
		IfaceInvalid               interface{}
		TimeMapKey                 map[time.Time]string
		ExistingArray              []string
		ExistingArrayIndex         []string
	}

	defaultValues := url.Values{
		"name":                          []string{"joeybloggs"},
		"Ignore":                        []string{"ignore"},
		"Time":                          []string{"2016-01-02"},
		"TimePtr":                       []string{"2016-01-02"},
		"mp[key]":                       []string{"value"},
		"NilArray":                      []string{"1", "2"},
		"TooSmallArray":                 []string{"1", "2"},
		"TooSmallCapOKArray":            []string{"1", "2"},
		"ZeroLengthArray":               []string{},
		"TooSmallNumberedArray[2]":      []string{"2"},
		"TooSmallCapOKNumberedArray[2]": []string{"2"},
		"BigEnoughNumberedArray[2]":     []string{"1"},
		"TimeMapKey[2016-01-02]":        []string{"time"},
		"ExistingArray":                 []string{"arr2"},
		"ExistingArrayIndex[1]":         []string{"arr2"},
	}

	testCases := []struct {
		NamespacePrefix string
		NamespaceSuffix string
		Values          url.Values
	}{{
		NamespacePrefix: ".",
		Values: url.Values{
			"Phone[0].Number":                []string{"1(111)111-1111"},
			"Phone[1].Number":                []string{"9(999)999-9999"},
			"PhonePtr[0].Number":             []string{"1(111)111-1111"},
			"PhonePtr[1].Number":             []string{"9(999)999-9999"},
			"NonNilPtr.Number":               []string{"9(999)999-9999"},
			"Anonymous.Value":                []string{"Anon"},
			"TestMapKeys.MapIfaceKey[key]":   []string{"3"},
			"TestMapKeys.MapFloat32Key[0.0]": []string{"3.3"},
			"TestMapKeys.MapFloat64Key[0.0]": []string{"3.3"},
			"TestMapKeys.MapNestedInt[1][2]": []string{"3"},
			"TestMapKeys.MapInt8[0]":         []string{"3"},
			"TestMapKeys.MapInt16[0]":        []string{"3"},
			"TestMapKeys.MapInt32[0]":        []string{"3"},
			"TestMapKeys.MapUint8[0]":        []string{"3"},
			"TestMapKeys.MapUint16[0]":       []string{"3"},
			"TestMapKeys.MapUint32[0]":       []string{"3"},
		},
	}, {
		NamespacePrefix: "[",
		NamespaceSuffix: "]",
		Values: url.Values{
			"Phone[0][Number]":                []string{"1(111)111-1111"},
			"Phone[1][Number]":                []string{"9(999)999-9999"},
			"PhonePtr[0][Number]":             []string{"1(111)111-1111"},
			"PhonePtr[1][Number]":             []string{"9(999)999-9999"},
			"NonNilPtr[Number]":               []string{"9(999)999-9999"},
			"Anonymous[Value]":                []string{"Anon"},
			"TestMapKeys[MapIfaceKey][key]":   []string{"3"},
			"TestMapKeys[MapFloat32Key][0.0]": []string{"3.3"},
			"TestMapKeys[MapFloat64Key][0.0]": []string{"3.3"},
			"TestMapKeys[MapNestedInt][1][2]": []string{"3"},
			"TestMapKeys[MapInt8][0]":         []string{"3"},
			"TestMapKeys[MapInt16][0]":        []string{"3"},
			"TestMapKeys[MapInt32][0]":        []string{"3"},
			"TestMapKeys[MapUint8][0]":        []string{"3"},
			"TestMapKeys[MapUint16][0]":       []string{"3"},
			"TestMapKeys[MapUint32][0]":       []string{"3"},
		},
	}}

	for _, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("Namespace_%s%s", tc.NamespacePrefix, tc.NamespaceSuffix), func(t *testing.T) {
			decoder := NewDecoder()
			decoder.SetNamespacePrefix(tc.NamespacePrefix)
			decoder.SetNamespaceSuffix(tc.NamespaceSuffix)

			values := url.Values{}

			for key, vals := range defaultValues {
				values[key] = vals
			}
			for key, vals := range tc.Values {
				values[key] = vals
			}

			decoder.SetTagName("form")
			decoder.RegisterCustomTypeFunc(func(vals []string) (interface{}, error) {
				return time.Parse("2006-01-02", vals[0])
			}, time.Time{})

			var test TestStruct
			test.ExistingMap = map[string]string{"existingkey": "existingvalue"}
			test.NonNilPtr = new(Phone)
			test.IfaceNonNil = new(Phone)
			test.IfaceInvalid = nil
			test.TooSmallArray = []string{"0"}
			test.TooSmallCapOKArray = make([]string, 0, 10)
			test.TooSmallNumberedArray = []string{"0"}
			test.TooSmallCapOKNumberedArray = make([]string, 0, 10)
			test.BigEnoughNumberedArray = make([]string, 3, 10)
			test.ExistingArray = []string{"arr1"}
			test.ExistingArrayIndex = []string{"arr1"}

			errs := decoder.Decode(&test, values)
			Equal(t, errs, nil)

			Equal(t, test.Name, "joeybloggs")
			Equal(t, test.Ignore, "")
			Equal(t, len(test.Phone), 2)
			Equal(t, test.Phone[0].Number, "1(111)111-1111")
			Equal(t, test.Phone[1].Number, "9(999)999-9999")
			Equal(t, len(test.PhonePtr), 2)
			Equal(t, test.PhonePtr[0].Number, "1(111)111-1111")
			Equal(t, test.PhonePtr[1].Number, "9(999)999-9999")
			Equal(t, test.NonNilPtr.Number, "9(999)999-9999")
			Equal(t, test.Anonymous.Value, "Anon")
			Equal(t, len(test.ExistingMap), 2)
			Equal(t, test.ExistingMap["existingkey"], "existingvalue")
			Equal(t, test.ExistingMap["key"], "value")
			Equal(t, len(test.NilArray), 2)
			Equal(t, test.NilArray[0], "1")
			Equal(t, test.NilArray[1], "2")
			Equal(t, len(test.TooSmallArray), 3)
			Equal(t, test.TooSmallArray[0], "0")
			Equal(t, test.TooSmallArray[1], "1")
			Equal(t, test.TooSmallArray[2], "2")
			Equal(t, len(test.ZeroLengthArray), 0)
			Equal(t, len(test.TooSmallNumberedArray), 3)
			Equal(t, test.TooSmallNumberedArray[0], "0")
			Equal(t, test.TooSmallNumberedArray[1], "")
			Equal(t, test.TooSmallNumberedArray[2], "2")
			Equal(t, len(test.BigEnoughNumberedArray), 3)
			Equal(t, cap(test.BigEnoughNumberedArray), 10)
			Equal(t, test.BigEnoughNumberedArray[0], "")
			Equal(t, test.BigEnoughNumberedArray[1], "")
			Equal(t, test.BigEnoughNumberedArray[2], "1")
			Equal(t, len(test.TooSmallCapOKArray), 2)
			Equal(t, cap(test.TooSmallCapOKArray), 10)
			Equal(t, test.TooSmallCapOKArray[0], "1")
			Equal(t, test.TooSmallCapOKArray[1], "2")
			Equal(t, len(test.TooSmallCapOKNumberedArray), 3)
			Equal(t, cap(test.TooSmallCapOKNumberedArray), 10)
			Equal(t, test.TooSmallCapOKNumberedArray[0], "")
			Equal(t, test.TooSmallCapOKNumberedArray[1], "")
			Equal(t, test.TooSmallCapOKNumberedArray[2], "2")

			Equal(t, len(test.ExistingArray), 2)
			Equal(t, test.ExistingArray[0], "arr1")
			Equal(t, test.ExistingArray[1], "arr2")

			Equal(t, len(test.ExistingArrayIndex), 2)
			Equal(t, test.ExistingArrayIndex[0], "arr1")
			Equal(t, test.ExistingArrayIndex[1], "arr2")

			tm, _ := time.Parse("2006-01-02", "2016-01-02")
			Equal(t, test.Time.Equal(tm), true)
			Equal(t, test.TimePtr.Equal(tm), true)

			NotEqual(t, test.TimeMapKey, nil)
			Equal(t, len(test.TimeMapKey), 1)

			_, ok := test.TimeMapKey[tm]
			Equal(t, ok, true)

			s := struct {
				Value     string
				Ignore    string `form:"-"`
				unexposed string
			}{}

			errs = decoder.Decode(&s, defaultValues)
			Equal(t, errs, nil)
			Equal(t, s.Value, "")
			Equal(t, s.Ignore, "")
			Equal(t, s.unexposed, "")

			Equal(t, test.TestMapKeys.MapIfaceKey["key"], "3")
			Equal(t, test.TestMapKeys.MapFloat32Key[float32(0.0)], float32(3.3))
			Equal(t, test.TestMapKeys.MapFloat64Key[float64(0.0)], float64(3.3))

			Equal(t, test.TestMapKeys.MapInt8[int8(0)], int8(3))
			Equal(t, test.TestMapKeys.MapInt16[int16(0)], int16(3))
			Equal(t, test.TestMapKeys.MapInt32[int32(0)], int32(3))

			Equal(t, test.TestMapKeys.MapUint8[uint8(0)], uint8(3))
			Equal(t, test.TestMapKeys.MapUint16[uint16(0)], uint16(3))
			Equal(t, test.TestMapKeys.MapUint32[uint32(0)], uint32(3))

			Equal(t, len(test.TestMapKeys.MapNestedInt), 1)
			Equal(t, len(test.TestMapKeys.MapNestedInt[1]), 1)
			Equal(t, test.TestMapKeys.MapNestedInt[1][2], 3)
		})
	}
}

func TestDecoderNativeTime(t *testing.T) {

	type TestError struct {
		Time        time.Time
		TimeNoValue time.Time
		TimePtr     *time.Time
	}

	values := url.Values{
		"Time":        []string{"2006-01-02T15:04:05Z"},
		"TimeNoValue": []string{""},
		"TimePtr":     []string{"2006-01-02T15:04:05Z"},
	}

	var test TestError

	decoder := NewDecoder()

	errs := decoder.Decode(&test, values)
	Equal(t, errs, nil)

	tm, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	Equal(t, test.Time.Equal(tm), true)
	Equal(t, test.TimeNoValue.Equal(tm), false)

	NotEqual(t, test.TimePtr, nil)
	Equal(t, (*test.TimePtr).Equal(tm), true)
}

func TestDecoderErrors(t *testing.T) {

	type TestError struct {
		Bool                  bool `form:"bool"`
		Int                   int
		Int8                  int8
		Int16                 int16
		Int32                 int32
		Uint                  uint
		Uint8                 uint8
		Uint16                uint16
		Uint32                uint32
		Float32               float32
		Float64               float64
		String                string
		Time                  time.Time
		MapBadIntKey          map[int]int
		MapBadInt8Key         map[int8]int8
		MapBadInt16Key        map[int16]int16
		MapBadInt32Key        map[int32]int32
		MapBadUintKey         map[uint]uint
		MapBadUint8Key        map[uint8]uint8
		MapBadUint16Key       map[uint16]uint16
		MapBadUint32Key       map[uint32]uint32
		MapBadFloat32Key      map[float32]float32
		MapBadFloat64Key      map[float64]float64
		MapBadBoolKey         map[bool]bool
		MapBadKeyType         map[complex64]int
		BadArrayValue         []int
		BadMapKey             map[time.Time]string
		OverflowNilArray      []int
		OverFlowExistingArray []int
		BadArrayIndex         []int
	}

	values := url.Values{
		"bool":                       []string{"uh-huh"},
		"Int":                        []string{"bad"},
		"Int8":                       []string{"bad"},
		"Int16":                      []string{"bad"},
		"Int32":                      []string{"bad"},
		"Uint":                       []string{"bad"},
		"Uint8":                      []string{"bad"},
		"Uint16":                     []string{"bad"},
		"Uint32":                     []string{"bad"},
		"Float32":                    []string{"bad"},
		"Float64":                    []string{"bad"},
		"String":                     []string{"str bad return val"},
		"Time":                       []string{"bad"},
		"MapBadIntKey[key]":          []string{"1"},
		"MapBadInt8Key[key]":         []string{"1"},
		"MapBadInt16Key[key]":        []string{"1"},
		"MapBadInt32Key[key]":        []string{"1"},
		"MapBadUintKey[key]":         []string{"1"},
		"MapBadUint8Key[key]":        []string{"1"},
		"MapBadUint16Key[key]":       []string{"1"},
		"MapBadUint32Key[key]":       []string{"1"},
		"MapBadFloat32Key[key]":      []string{"1.1"},
		"MapBadFloat64Key[key]":      []string{"1.1"},
		"MapBadBoolKey[uh-huh]":      []string{"true"},
		"MapBadKeyType[1.4]":         []string{"5"},
		"BadArrayValue[0]":           []string{"badintval"},
		"BadMapKey[badtime]":         []string{"badtime"},
		"OverflowNilArray[999]":      []string{"idx 1000"},
		"OverFlowExistingArray[999]": []string{"idx 1000"},
		"BadArrayIndex[bad index]":   []string{"bad idx"},
	}

	testCases := []struct {
		NamespacePrefix string
		NamespaceSuffix string
	}{{
		NamespacePrefix: ".",
	}, {
		NamespacePrefix: "[",
		NamespaceSuffix: "]",
	}}

	for _, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("Namespace_%s%s", tc.NamespacePrefix, tc.NamespaceSuffix), func(t *testing.T) {
			decoder := NewDecoder()
			decoder.SetNamespacePrefix(tc.NamespacePrefix)
			decoder.SetNamespaceSuffix(tc.NamespaceSuffix)

			decoder.SetMaxArraySize(4)
			decoder.RegisterCustomTypeFunc(func(vals []string) (interface{}, error) {
				return nil, errors.New("Bad Type Conversion")
			}, "")

			test := TestError{
				OverFlowExistingArray: make([]int, 2),
			}

			errs := decoder.Decode(&test, values)
			NotEqual(t, errs, nil)

			e := errs.Error()
			NotEqual(t, e, "")

			err := errs.(DecodeErrors)
			Equal(t, len(err), 30)

			k := err["bool"]
			Equal(t, k.Error(), "Invalid Boolean Value 'uh-huh' Type 'bool' Namespace 'bool'")

			k = err["Int"]
			Equal(t, k.Error(), "Invalid Integer Value 'bad' Type 'int' Namespace 'Int'")

			k = err["Int8"]
			Equal(t, k.Error(), "Invalid Integer Value 'bad' Type 'int8' Namespace 'Int8'")

			k = err["Int16"]
			Equal(t, k.Error(), "Invalid Integer Value 'bad' Type 'int16' Namespace 'Int16'")

			k = err["Int32"]
			Equal(t, k.Error(), "Invalid Integer Value 'bad' Type 'int32' Namespace 'Int32'")

			k = err["Uint"]
			Equal(t, k.Error(), "Invalid Unsigned Integer Value 'bad' Type 'uint' Namespace 'Uint'")

			k = err["Uint8"]
			Equal(t, k.Error(), "Invalid Unsigned Integer Value 'bad' Type 'uint8' Namespace 'Uint8'")

			k = err["Uint16"]
			Equal(t, k.Error(), "Invalid Unsigned Integer Value 'bad' Type 'uint16' Namespace 'Uint16'")

			k = err["Uint32"]
			Equal(t, k.Error(), "Invalid Unsigned Integer Value 'bad' Type 'uint32' Namespace 'Uint32'")

			k = err["Float32"]
			Equal(t, k.Error(), "Invalid Float Value 'bad' Type 'float32' Namespace 'Float32'")

			k = err["Float64"]
			Equal(t, k.Error(), "Invalid Float Value 'bad' Type 'float64' Namespace 'Float64'")

			k = err["String"]
			Equal(t, k.Error(), "Bad Type Conversion")

			k = err["Time"]
			Equal(t, k.Error(), "parsing time \"bad\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"bad\" as \"2006\"")

			k = err["MapBadIntKey"]
			Equal(t, k.Error(), "Invalid Integer Value 'key' Type 'int' Namespace 'MapBadIntKey'")

			k = err["MapBadInt8Key"]
			Equal(t, k.Error(), "Invalid Integer Value 'key' Type 'int8' Namespace 'MapBadInt8Key'")

			k = err["MapBadInt16Key"]
			Equal(t, k.Error(), "Invalid Integer Value 'key' Type 'int16' Namespace 'MapBadInt16Key'")

			k = err["MapBadInt32Key"]
			Equal(t, k.Error(), "Invalid Integer Value 'key' Type 'int32' Namespace 'MapBadInt32Key'")

			k = err["MapBadUintKey"]
			Equal(t, k.Error(), "Invalid Unsigned Integer Value 'key' Type 'uint' Namespace 'MapBadUintKey'")

			k = err["MapBadUint8Key"]
			Equal(t, k.Error(), "Invalid Unsigned Integer Value 'key' Type 'uint8' Namespace 'MapBadUint8Key'")

			k = err["MapBadUint16Key"]
			Equal(t, k.Error(), "Invalid Unsigned Integer Value 'key' Type 'uint16' Namespace 'MapBadUint16Key'")

			k = err["MapBadUint32Key"]
			Equal(t, k.Error(), "Invalid Unsigned Integer Value 'key' Type 'uint32' Namespace 'MapBadUint32Key'")

			k = err["MapBadFloat32Key"]
			Equal(t, k.Error(), "Invalid Float Value 'key' Type 'float32' Namespace 'MapBadFloat32Key'")

			k = err["MapBadFloat64Key"]
			Equal(t, k.Error(), "Invalid Float Value 'key' Type 'float64' Namespace 'MapBadFloat64Key'")

			k = err["MapBadBoolKey"]
			Equal(t, k.Error(), "Invalid Boolean Value 'uh-huh' Type 'bool' Namespace 'MapBadBoolKey'")

			k = err["MapBadKeyType"]
			Equal(t, k.Error(), "Unsupported Map Key '1.4', Type 'complex64' Namespace 'MapBadKeyType'")

			k = err["BadArrayValue[0]"]
			Equal(t, k.Error(), "Invalid Integer Value 'badintval' Type 'int' Namespace 'BadArrayValue[0]'")

			k = err["OverflowNilArray"]
			Equal(t, k.Error(), "Array size of '1000' is larger than the maximum currently set on the decoder of '4'. To increase this limit please see, SetMaxArraySize(size uint)")

			k = err["OverFlowExistingArray"]
			Equal(t, k.Error(), "Array size of '1000' is larger than the maximum currently set on the decoder of '4'. To increase this limit please see, SetMaxArraySize(size uint)")

			k = err["BadArrayIndex"]
			Equal(t, k.Error(), "invalid slice index 'bad index'")

			type TestError2 struct {
				BadMapKey map[time.Time]string
			}

			values2 := url.Values{
				"BadMapKey[badtime]": []string{"badtime"},
			}

			var test2 TestError2
			decoder2 := NewDecoder()
			decoder2.RegisterCustomTypeFunc(func(vals []string) (interface{}, error) {
				return time.Parse("2006-01-02", vals[0])
			}, time.Time{})

			errs = decoder2.Decode(&test2, values2)
			NotEqual(t, errs, nil)

			e = errs.Error()
			NotEqual(t, e, "")

			k = err["BadMapKey"]
			Equal(t, k.Error(), "Unsupported Map Key 'badtime', Type 'time.Time' Namespace 'BadMapKey'")
		})
	}
}

func TestDecodeAllTypes(t *testing.T) {

	values := url.Values{
		"": []string{"3"},
	}

	decoder := NewDecoder()

	var i int

	errs := decoder.Decode(&i, values)
	Equal(t, errs, nil)
	Equal(t, i, 3)

	var i8 int

	errs = decoder.Decode(&i8, values)
	Equal(t, errs, nil)
	Equal(t, i8, 3)

	var i16 int

	errs = decoder.Decode(&i16, values)
	Equal(t, errs, nil)
	Equal(t, i16, 3)

	var i32 int

	errs = decoder.Decode(&i32, values)
	Equal(t, errs, nil)
	Equal(t, i32, 3)

	var i64 int

	errs = decoder.Decode(&i64, values)
	Equal(t, errs, nil)
	Equal(t, i64, 3)

	var ui int

	errs = decoder.Decode(&ui, values)
	Equal(t, errs, nil)
	Equal(t, ui, 3)

	var ui8 int

	errs = decoder.Decode(&ui8, values)
	Equal(t, errs, nil)
	Equal(t, ui8, 3)

	var ui16 int

	errs = decoder.Decode(&ui16, values)
	Equal(t, errs, nil)
	Equal(t, ui16, 3)

	var ui32 int

	errs = decoder.Decode(&ui32, values)
	Equal(t, errs, nil)
	Equal(t, ui32, 3)

	var ui64 int

	errs = decoder.Decode(&ui64, values)
	Equal(t, errs, nil)
	Equal(t, ui64, 3)

	values = url.Values{
		"": []string{"3.4"},
	}

	var f32 float32

	errs = decoder.Decode(&f32, values)
	Equal(t, errs, nil)
	Equal(t, f32, float32(3.4))

	var f64 float64

	errs = decoder.Decode(&f64, values)
	Equal(t, errs, nil)
	Equal(t, f64, float64(3.4))

	values = url.Values{
		"": []string{"true"},
	}

	var b bool

	errs = decoder.Decode(&b, values)
	Equal(t, errs, nil)
	Equal(t, b, true)

	tm, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")

	values = url.Values{
		"": []string{"2006-01-02T15:04:05Z"},
	}

	var dt time.Time

	errs = decoder.Decode(&dt, values)

	Equal(t, errs, nil)
	Equal(t, dt, tm)

	values = url.Values{
		"": []string{"arr1", "arr2"},
	}

	// basic array
	var arr []string

	errs = decoder.Decode(&arr, values)
	Equal(t, errs, nil)
	Equal(t, len(arr), 2)
	Equal(t, arr[0], "arr1")
	Equal(t, arr[1], "arr2")

	// pre-populated array

	// fmt.Println("Decoding...")
	errs = decoder.Decode(&arr, values)
	Equal(t, errs, nil)
	Equal(t, len(arr), 4)
	Equal(t, arr[0], "arr1")
	Equal(t, arr[1], "arr2")
	Equal(t, arr[2], "arr1")
	Equal(t, arr[3], "arr2")

	// basic array Ptr
	var arrPtr []*string

	errs = decoder.Decode(&arrPtr, values)
	Equal(t, errs, nil)
	Equal(t, len(arrPtr), 2)
	Equal(t, *arrPtr[0], "arr1")
	Equal(t, *arrPtr[1], "arr2")

	// pre-populated array Ptr

	// fmt.Println("Decoding...")
	errs = decoder.Decode(&arrPtr, values)
	Equal(t, errs, nil)
	Equal(t, len(arr), 4)
	Equal(t, *arrPtr[0], "arr1")
	Equal(t, *arrPtr[1], "arr2")
	Equal(t, *arrPtr[2], "arr1")
	Equal(t, *arrPtr[3], "arr2")

	// indexed array

	values = url.Values{
		"[0]": []string{"newVal1"},
		"[1]": []string{"newVal2"},
	}

	errs = decoder.Decode(&arr, values)
	Equal(t, errs, nil)
	Equal(t, len(arr), 4)
	Equal(t, arr[0], "newVal1")
	Equal(t, arr[1], "newVal2")
	Equal(t, arr[2], "arr1")
	Equal(t, arr[3], "arr2")

	values = url.Values{
		"[key1]": []string{"val1"},
		"[key2]": []string{"val2"},
	}

	// basic map
	var m map[string]string

	errs = decoder.Decode(&m, values)
	Equal(t, errs, nil)
	Equal(t, len(m), 2)
	Equal(t, m["key1"], "val1")
	Equal(t, m["key2"], "val2")

	// existing map

	errs = decoder.Decode(&m, values)
	Equal(t, errs, nil)
	Equal(t, len(m), 2)
	Equal(t, m["key1"], "val1")
	Equal(t, m["key2"], "val2")

	// basic map, adding more keys

	values = url.Values{
		"[key3]": []string{"val3"},
	}
	errs = decoder.Decode(&m, values)
	Equal(t, errs, nil)
	Equal(t, len(m), 3)
	Equal(t, m["key3"], "val3")

	// array of struct

	type Phone struct {
		Number string
		Label  string
	}

	values = url.Values{
		"[0].Number": []string{"999"},
		"[1].Label":  []string{"label2"},
		"[1].Number": []string{"111"},
		"[0].Label":  []string{"label1"},
	}

	var phones []Phone

	errs = decoder.Decode(&phones, values)
	Equal(t, errs, nil)
	Equal(t, len(phones), 2)
	Equal(t, phones[0].Number, "999")
	Equal(t, phones[0].Label, "label1")
	Equal(t, phones[1].Number, "111")
	Equal(t, phones[1].Label, "label2")
}

func TestDecoderPanicsAndBadValues(t *testing.T) {

	type Phone struct {
		Number string
	}

	type TestError struct {
		Phone  []Phone
		Phone2 []Phone
		Phone3 []Phone
	}

	values := url.Values{
		"Phone[0.Number": []string{"1(111)111-1111"},
	}

	var test TestError

	decoder := NewDecoder()

	PanicMatches(t, func() { _ = decoder.Decode(&test, values) }, "Invalid formatting for key 'Phone[0.Number' missing ']' bracket")

	i := 1
	err := decoder.Decode(i, values)
	NotEqual(t, err, nil)

	_, ok := err.(*InvalidDecoderError)
	Equal(t, ok, true)
	Equal(t, err.Error(), "form: Decode(non-pointer int)")

	err = decoder.Decode(nil, values)
	NotEqual(t, err, nil)

	_, ok = err.(*InvalidDecoderError)
	Equal(t, ok, true)
	Equal(t, err.Error(), "form: Decode(nil)")

	var ts *TestError

	err = decoder.Decode(ts, values)
	NotEqual(t, err, nil)

	_, ok = err.(*InvalidDecoderError)
	Equal(t, ok, true)
	Equal(t, err.Error(), "form: Decode(nil *form.TestError)")

	values = url.Values{
		"Phone0].Number": []string{"1(111)111-1111"},
	}

	PanicMatches(t, func() { _ = decoder.Decode(&test, values) }, "Invalid formatting for key 'Phone0].Number' missing '[' bracket")

	values = url.Values{
		"Phone[[0.Number": []string{"1(111)111-1111"},
	}

	PanicMatches(t, func() { _ = decoder.Decode(&test, values) }, "Invalid formatting for key 'Phone[[0.Number' missing ']' bracket")

	values = url.Values{
		"Phone0]].Number": []string{"1(111)111-1111"},
	}

	PanicMatches(t, func() { _ = decoder.Decode(&test, values) }, "Invalid formatting for key 'Phone0]].Number' missing '[' bracket")
}

func TestDecoderMapKeys(t *testing.T) {

	type TestMapKeys struct {
		MapIfaceKey   map[interface{}]string
		MapFloat32Key map[float32]float32
		MapFloat64Key map[float64]float64
		MapNestedInt  map[int]map[int]int
		MapInt8       map[int8]int8
		MapInt16      map[int16]int16
		MapInt32      map[int32]int32
		MapUint8      map[uint8]uint8
		MapUint16     map[uint16]uint16
		MapUint32     map[uint32]uint32
	}

	values := url.Values{
		"MapIfaceKey[key]":   []string{"3"},
		"MapFloat32Key[0.0]": []string{"3.3"},
		"MapFloat64Key[0.0]": []string{"3.3"},
		"MapNestedInt[1][2]": []string{"3"},
		"MapInt8[0]":         []string{"3"},
		"MapInt16[0]":        []string{"3"},
		"MapInt32[0]":        []string{"3"},
		"MapUint8[0]":        []string{"3"},
		"MapUint16[0]":       []string{"3"},
		"MapUint32[0]":       []string{"3"},
	}

	var test TestMapKeys

	testCases := []struct {
		NamespacePrefix string
		NamespaceSuffix string
	}{
		{
			NamespacePrefix: ".",
		},
		{
			NamespacePrefix: "[",
			NamespaceSuffix: "]",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("Namespace_%s%s", tc.NamespacePrefix, tc.NamespaceSuffix), func(t *testing.T) {
			decoder := NewDecoder()
			decoder.SetNamespacePrefix(tc.NamespacePrefix)
			decoder.SetNamespaceSuffix(tc.NamespaceSuffix)

			errs := decoder.Decode(&test, values)
			Equal(t, errs, nil)

			Equal(t, test.MapIfaceKey["key"], "3")
			Equal(t, test.MapFloat32Key[float32(0.0)], float32(3.3))
			Equal(t, test.MapFloat64Key[float64(0.0)], float64(3.3))

			Equal(t, test.MapInt8[int8(0)], int8(3))
			Equal(t, test.MapInt16[int16(0)], int16(3))
			Equal(t, test.MapInt32[int32(0)], int32(3))

			Equal(t, test.MapUint8[uint8(0)], uint8(3))
			Equal(t, test.MapUint16[uint16(0)], uint16(3))
			Equal(t, test.MapUint32[uint32(0)], uint32(3))

			Equal(t, len(test.MapNestedInt), 1)
			Equal(t, len(test.MapNestedInt[1]), 1)
			Equal(t, test.MapNestedInt[1][2], 3)
		})
	}
}

func TestDecoderStructRecursion(t *testing.T) {

	type Nested struct {
		Value  string
		Nested *Nested
	}

	type TestRecursive struct {
		Nested    Nested
		NestedPtr *Nested
		NestedTwo Nested
	}

	defaultValues := url.Values{
		"Value": []string{"value"},
	}

	testCases := []struct {
		Values          url.Values
		NamespacePrefix string
		NamespaceSuffix string
	}{{
		NamespacePrefix: ".",
		Values: url.Values{
			"Nested.Value":           []string{"value"},
			"NestedPtr.Value":        []string{"value"},
			"NestedTwo.Nested.Value": []string{"value"},
		},
	}, {
		NamespacePrefix: "[",
		NamespaceSuffix: "]",
		Values: url.Values{
			"Nested[Value]":            []string{"value"},
			"NestedPtr[Value]":         []string{"value"},
			"NestedTwo[Nested][Value]": []string{"value"},
		},
	}}

	for _, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("Namespace_%s%s", tc.NamespacePrefix, tc.NamespaceSuffix), func(t *testing.T) {
			values := url.Values{}

			for key, vals := range defaultValues {
				values[key] = vals
			}
			for key, vals := range tc.Values {
				values[key] = vals
			}

			decoder := NewDecoder()
			decoder.SetNamespacePrefix(tc.NamespacePrefix)
			decoder.SetNamespaceSuffix(tc.NamespaceSuffix)

			var test TestRecursive

			errs := decoder.Decode(&test, values)
			Equal(t, errs, nil)

			Equal(t, test.Nested.Value, "value")
			Equal(t, test.NestedPtr.Value, "value")
			Equal(t, test.Nested.Nested, nil)
			Equal(t, test.NestedTwo.Nested.Value, "value")
		})
	}

}

func TestDecoderFormDecode(t *testing.T) {

	type Struct2 struct {
		Foo string
		Bar string
	}

	type Struct2Wrapper struct {
		InnerSlice []Struct2
	}

	sliceValues := map[string][]string{
		"InnerSlice[0].Foo": {"foo-is-set"},
	}

	singleValues := map[string][]string{
		"Foo": {"foo-is-set"},
	}

	fd := NewDecoder()

	dst := Struct2Wrapper{}
	err := fd.Decode(&dst, sliceValues)
	Equal(t, err, nil)
	NotEqual(t, dst.InnerSlice, nil)
	Equal(t, dst.InnerSlice[0].Foo, "foo-is-set")

	dst2 := Struct2{}
	err = fd.Decode(&dst2, singleValues)
	Equal(t, err, nil)
	Equal(t, dst2.Foo, "foo-is-set")
}

func TestDecoderArrayKeysSort(t *testing.T) {

	type Struct struct {
		Array []int
	}

	values := map[string][]string{

		"Array[2]":  {"2"},
		"Array[10]": {"10"},
	}

	var test Struct

	d := NewDecoder()

	err := d.Decode(&test, values)
	Equal(t, err, nil)

	Equal(t, len(test.Array), 11)
	Equal(t, test.Array[2], int(2))
	Equal(t, test.Array[10], int(10))
}

func TestDecoderIncreasingKeys(t *testing.T) {

	type Struct struct {
		Array []int
	}

	values := map[string][]string{
		"Array[2]": {"2"},
	}

	var test Struct

	d := NewDecoder()

	err := d.Decode(&test, values)
	Equal(t, err, nil)

	Equal(t, len(test.Array), 3)
	Equal(t, test.Array[2], int(2))

	values["Array[10]"] = []string{"10"}

	var test2 Struct

	err = d.Decode(&test2, values)
	Equal(t, err, nil)

	Equal(t, len(test2.Array), 11)
	Equal(t, test2.Array[2], int(2))
	Equal(t, test2.Array[10], int(10))
}

func TestDecoderInterface(t *testing.T) {

	var iface interface{}

	d := NewDecoder()

	values := map[string][]string{
		"": {"2"},
	}

	var i int

	iface = &i

	err := d.Decode(iface, values)
	Equal(t, err, nil)
	Equal(t, i, 2)

	iface = i

	err = d.Decode(iface, values)
	NotEqual(t, err, nil)

	_, ok := err.(*InvalidDecoderError)
	Equal(t, ok, true)
	Equal(t, err.Error(), "form: Decode(non-pointer int)")

	values = map[string][]string{
		"Value": {"testVal"},
	}

	type test struct {
		Value string
	}

	var tst test

	iface = &tst

	err = d.Decode(iface, values)
	Equal(t, err, nil)
	Equal(t, tst.Value, "testVal")

	iface = tst

	err = d.Decode(iface, values)
	NotEqual(t, err, nil)

	_, ok = err.(*InvalidDecoderError)
	Equal(t, ok, true)
	Equal(t, err.Error(), "form: Decode(non-pointer form.test)")
}

func TestDecoderPointerToPointer(t *testing.T) {

	values := map[string][]string{
		"Value": {"testVal"},
	}

	type Test struct {
		Value string
	}

	var tst *Test

	d := NewDecoder()
	err := d.Decode(&tst, values)
	Equal(t, err, nil)
	Equal(t, tst.Value, "testVal")
}

func TestDecoderExplicit(t *testing.T) {

	type Test struct {
		Name string `form:"Name"`
		Age  int
	}

	values := map[string][]string{
		"Name": {"Joeybloggs"},
		"Age":  {"3"},
	}

	var test Test

	d := NewDecoder()
	d.SetMode(ModeExplicit)

	err := d.Decode(&test, values)
	Equal(t, err, nil)
	Equal(t, test.Name, "Joeybloggs")
	Equal(t, test.Age, 0)
}

func TestDecoderStructWithJSONTag(t *testing.T) {
	type Test struct {
		Name string `json:"name,omitempty"`
		Age  int    `json:",omitempty"`
	}

	values := map[string][]string{
		"name": {"Joeybloggs"},
		"Age":  {"3"},
	}

	var test Test

	d := NewDecoder()
	d.SetTagName("json")

	err := d.Decode(&test, values)
	Equal(t, err, nil)
	Equal(t, test.Name, "Joeybloggs")
	Equal(t, test.Age, int(3))
}

func TestDecoderRegisterTagNameFunc(t *testing.T) {

	type Test struct {
		Value  string `json:"val,omitempty"`
		Ignore string `json:"-"`
	}

	values := url.Values{
		"val":    []string{"joeybloggs"},
		"Ignore": []string{"ignore"},
	}

	var test Test

	decoder := NewDecoder()
	decoder.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("json")

		if commaIndex := strings.Index(name, ","); commaIndex != -1 {
			name = name[:commaIndex]
		}

		return name
	})

	err := decoder.Decode(&test, values)
	Equal(t, err, nil)
	Equal(t, test.Value, "joeybloggs")
	Equal(t, test.Ignore, "")
}

func TestDecoderEmbedModes(t *testing.T) {

	type A struct {
		Field string
	}

	type B struct {
		A
		Field string
	}

	var b B

	decoder := NewDecoder()

	values := url.Values{
		"Field": []string{"Value"},
	}

	err := decoder.Decode(&b, values)
	Equal(t, err, nil)
	Equal(t, b.Field, "Value")
	Equal(t, b.A.Field, "Value")

	values = url.Values{
		"Field":   []string{"B Val"},
		"A.Field": []string{"A Val"},
	}

	err = decoder.Decode(&b, values)
	Equal(t, err, nil)
	Equal(t, b.Field, "B Val")
	Equal(t, b.A.Field, "A Val")
}

func TestInterfaceDecoding(t *testing.T) {

	type Test struct {
		Iface interface{}
	}

	var b Test

	values := url.Values{
		"Iface": []string{"1"},
	}

	decoder := NewDecoder()
	err := decoder.Decode(&b, values)
	Equal(t, err, nil)
	Equal(t, b.Iface, "1")
}

func TestDecodeArrayBug(t *testing.T) {
	var data struct {
		A [2]string
		B [2]string
		C [2]string
		D [3]string
		E [3]string
		F [3]string
		G [3]string
	}
	decoder := NewDecoder()
	err := decoder.Decode(&data, url.Values{
		// Mixed types
		"A":    {"10"},
		"A[1]": {"20"},
		// overflow
		"B":    {"10", "20", "30"},
		"B[1]": {"31", "10", "20"},
		"B[2]": {"40"},
		// invalid array index
		"C[q]": {""},
		// index and mix tests
		"D":    {"10"},
		"E":    {"10", "20"},
		"F":    {"10", "", "20"},
		"G":    {"10"},
		"G[2]": {"20"},
	})
	NotEqual(t, err, nil)
	Equal(t, err.Error(), "Field Namespace:C ERROR:invalid array index 'q'")
	Equal(t, data.A[0], "10")
	Equal(t, data.A[1], "20")
	Equal(t, data.B[0], "10")
	Equal(t, data.B[1], "31")
	Equal(t, data.C[0], "")
	Equal(t, data.C[1], "")
	Equal(t, data.D[0], "10")
	Equal(t, data.D[1], "")
	Equal(t, data.D[2], "")
	Equal(t, data.E[0], "10")
	Equal(t, data.E[1], "20")
	Equal(t, data.E[2], "")
	Equal(t, data.F[0], "10")
	Equal(t, data.F[1], "")
	Equal(t, data.F[2], "20")
	Equal(t, data.G[0], "10")
	Equal(t, data.G[1], "")
	Equal(t, data.G[2], "20")
}

func TestDecoder_RegisterCustomTypeFuncOnSlice(t *testing.T) {
	type customString string

	type TestStruct struct {
		Slice []customString `form:"slice"`
	}

	d := NewDecoder()
	d.RegisterCustomTypeFunc(func(vals []string) (i interface{}, e error) {
		custom := make([]customString, 0, len(vals))
		for i := 0; i < len(vals); i++ {
			custom = append(custom, customString("custom"+vals[i]))
		}
		return custom, nil
	}, []customString{})

	var v TestStruct
	err := d.Decode(&v, url.Values{"slice": []string{"v1", "v2"}})
	Equal(t, err, nil)
	Equal(t, v.Slice, []customString{"customv1", "customv2"})
}

func TestDecoder_RegisterCustomTypeFunc(t *testing.T) {
	type customString string

	type TestStruct struct {
		Slice []customString `form:"slice"`
	}

	d := NewDecoder()
	d.RegisterCustomTypeFunc(func(vals []string) (i interface{}, e error) {
		return customString("custom" + vals[0]), nil
	}, customString(""))

	var v TestStruct
	err := d.Decode(&v, url.Values{"slice": []string{"v1", "v2"}})
	Equal(t, err, nil)

	Equal(t, v.Slice, []customString{"customv1", "customv2"})
}

func TestDecoder_EmptyArrayString(t *testing.T) {
	type T1 struct {
		F1 string `form:"F1"`
	}
	in := url.Values{
		"F1": []string{},
	}

	v := new(T1)

	d := NewDecoder()
	err := d.Decode(v, in)
	Equal(t, err, nil)
}

func TestDecoder_EmptyArrayBool(t *testing.T) {
	type T1 struct {
		F1 bool `form:"F1"`
	}
	in := url.Values{
		"F1": []string{},
	}

	v := new(T1)
	d := NewDecoder()
	err := d.Decode(v, in)
	Equal(t, err, nil)
}

func TestDecoder_InvalidSliceIndex(t *testing.T) {
	type PostsRequest struct {
		PostIds []string
	}
	in := url.Values{
		"PostIds[]": []string{"1", "2"},
	}

	v := new(PostsRequest)
	d := NewDecoder()
	err := d.Decode(v, in)
	NotEqual(t, err, nil)
	Equal(t, err.Error(), "Field Namespace:PostIds ERROR:invalid slice index ''")

	// No error with proper name
	type PostsRequest2 struct {
		PostIds []string `form:"PostIds[]"`
	}

	v2 := new(PostsRequest2)
	err = d.Decode(v2, in)
	Equal(t, err, nil)
	Equal(t, len(v2.PostIds), 2)
	Equal(t, v2.PostIds[0], "1")
	Equal(t, v2.PostIds[1], "2")
}

package benchmarks

import (
	"net/url"
	"testing"

	"github.com/go-playground/form/v4"
)

// Simple Benchmarks

type User struct {
	FirstName string `form:"fname" schema:"fname" formam:"fname"`
	LastName  string `form:"lname" schema:"lname" formam:"lname"`
	Email     string `form:"email" schema:"email" formam:"email"`
	Age       uint8  `form:"age"   schema:"age"   formam:"age"`
}

func getUserStructValues() url.Values {
	return url.Values{
		"fname": []string{"Joey"},
		"lname": []string{"Bloggs"},
		"email": []string{"joeybloggs@gmail.com"},
		"age":   []string{"32"},
	}
}

func getUserStruct() *User {
	return &User{
		FirstName: "Joey",
		LastName:  "Bloggs",
		Email:     "joeybloggs@gmail.com",
		Age:       32,
	}
}

func BenchmarkSimpleUserDecodeStruct(b *testing.B) {

	values := getUserStructValues()
	decoder := form.NewDecoder()

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		var test User
		if err := decoder.Decode(&test, values); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkSimpleUserDecodeStructParallel(b *testing.B) {

	values := getUserStructValues()
	decoder := form.NewDecoder()

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var test User
			if err := decoder.Decode(&test, values); err != nil {
				b.Error(err)
			}
		}
	})
}
func BenchmarkSimpleUserEncodeStruct(b *testing.B) {

	test := getUserStruct()
	encoder := form.NewEncoder()

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		if _, err := encoder.Encode(&test); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkSimpleUserEncodeStructParallel(b *testing.B) {

	test := getUserStruct()
	encoder := form.NewEncoder()

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := encoder.Encode(&test); err != nil {
				b.Error(err)
			}
		}
	})
}

// Primitives ALL types

type PrimitivesStruct struct {
	String  string
	Int     int
	Int8    int8
	Int16   int16
	Int32   int32
	Int64   int64
	Uint    uint
	Uint8   uint8
	Uint16  uint16
	Uint32  uint32
	Uint64  uint64
	Float32 float32
	Float64 float64
	Bool    bool
}

func getPrimitivesStructValues() url.Values {
	return url.Values{
		"String":  []string{"joeybloggs"},
		"Int":     []string{"1"},
		"Int8":    []string{"2"},
		"Int16":   []string{"3"},
		"Int32":   []string{"4"},
		"Int64":   []string{"5"},
		"Uint":    []string{"1"},
		"Uint8":   []string{"2"},
		"Uint16":  []string{"3"},
		"Uint32":  []string{"4"},
		"Uint64":  []string{"5"},
		"Float32": []string{"1.1"},
		"Float64": []string{"5.0"},
		"Bool":    []string{"true"},
	}
}

func getPrimitivesStruct() *PrimitivesStruct {
	return &PrimitivesStruct{
		String:  "joeybloggs",
		Int:     1,
		Int8:    2,
		Int16:   3,
		Int32:   4,
		Int64:   5,
		Uint:    1,
		Uint8:   2,
		Uint16:  3,
		Uint32:  4,
		Uint64:  5,
		Float32: 1.1,
		Float64: 5.0,
		Bool:    true,
	}
}

func BenchmarkPrimitivesDecodeStructAllPrimitivesTypes(b *testing.B) {
	values := getPrimitivesStructValues()
	decoder := form.NewDecoder()

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		var test PrimitivesStruct
		if err := decoder.Decode(&test, values); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkPrimitivesDecodeStructAllPrimitivesTypesParallel(b *testing.B) {
	values := getPrimitivesStructValues()
	decoder := form.NewDecoder()

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var test PrimitivesStruct
			if err := decoder.Decode(&test, values); err != nil {
				b.Error(err)
			}
		}
	})
}

func BenchmarkPrimitivesEncodeStructAllPrimitivesTypes(b *testing.B) {
	test := getPrimitivesStruct()
	encoder := form.NewEncoder()

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		if _, err := encoder.Encode(&test); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkPrimitivesEncodeStructAllPrimitivesTypesParallel(b *testing.B) {
	test := getPrimitivesStruct()
	encoder := form.NewEncoder()

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := encoder.Encode(&test); err != nil {
				b.Error(err)
			}
		}
	})
}

// Complex Array ALL types

type ComplexArrayStruct struct {
	String       []string
	StringPtr    []*string
	Int          []int
	IntPtr       []*int
	Int8         []int8
	Int8Ptr      []*int8
	Int16        []int16
	Int16Ptr     []*int16
	Int32        []int32
	Int32Ptr     []*int32
	Int64        []int64
	Int64Ptr     []*int64
	Uint         []uint
	UintPtr      []*uint
	Uint8        []uint8
	Uint8Ptr     []*uint8
	Uint16       []uint16
	Uint16Ptr    []*uint16
	Uint32       []uint32
	Uint32Ptr    []*uint32
	Uint64       []uint64
	Uint64Ptr    []*uint64
	NestedInt    [][]int
	NestedIntPtr [][]*int
}

func getComplexArrayStructValues() url.Values {
	return url.Values{
		"String":             []string{"joeybloggs"},
		"StringPtr":          []string{"joeybloggs"},
		"Int":                []string{"1", "2"},
		"IntPtr":             []string{"1", "2"},
		"Int8[0]":            []string{"1"},
		"Int8[1]":            []string{"2"},
		"Int8Ptr[0]":         []string{"1"},
		"Int8Ptr[1]":         []string{"2"},
		"Int16":              []string{"1", "2"},
		"Int16Ptr":           []string{"1", "2"},
		"Int32":              []string{"1", "2"},
		"Int32Ptr":           []string{"1", "2"},
		"Int64":              []string{"1", "2"},
		"Int64Ptr":           []string{"1", "2"},
		"Uint":               []string{"1", "2"},
		"UintPtr":            []string{"1", "2"},
		"Uint8[0]":           []string{"1"},
		"Uint8[1]":           []string{"2"},
		"Uint8Ptr[0]":        []string{"1"},
		"Uint8Ptr[1]":        []string{"2"},
		"Uint16":             []string{"1", "2"},
		"Uint16Ptr":          []string{"1", "2"},
		"Uint32":             []string{"1", "2"},
		"Uint32Ptr":          []string{"1", "2"},
		"Uint64":             []string{"1", "2"},
		"Uint64Ptr":          []string{"1", "2"},
		"NestedInt[0][0]":    []string{"1"},
		"NestedIntPtr[0][1]": []string{"1"},
	}
}

func getComplexArrayStruct() *ComplexArrayStruct {

	s := "joeybloggs"

	i1 := int(1)
	i2 := int(2)
	i81 := int8(1)
	i82 := int8(2)
	i161 := int16(1)
	i162 := int16(2)
	i321 := int32(1)
	i322 := int32(2)
	i641 := int64(1)
	i642 := int64(2)

	ui1 := uint(1)
	ui2 := uint(2)
	ui81 := uint8(1)
	ui82 := uint8(2)
	ui161 := uint16(1)
	ui162 := uint16(2)
	ui321 := uint32(1)
	ui322 := uint32(2)
	ui641 := uint64(1)
	ui642 := uint64(2)

	return &ComplexArrayStruct{
		String:       []string{s},
		StringPtr:    []*string{&s},
		Int:          []int{i1, i2},
		IntPtr:       []*int{&i1, &i2},
		Int8:         []int8{i81, i82},
		Int8Ptr:      []*int8{&i81, &i82},
		Int16:        []int16{i161, i162},
		Int16Ptr:     []*int16{&i161, &i162},
		Int32:        []int32{i321, i322},
		Int32Ptr:     []*int32{&i321, &i322},
		Int64:        []int64{i641, i642},
		Int64Ptr:     []*int64{&i641, &i642},
		Uint:         []uint{ui1, ui2},
		UintPtr:      []*uint{&ui1, &ui2},
		Uint8:        []uint8{ui81, ui82},
		Uint8Ptr:     []*uint8{&ui81, &ui82},
		Uint16:       []uint16{ui161, ui162},
		Uint16Ptr:    []*uint16{&ui161, &ui162},
		Uint32:       []uint32{ui321, ui322},
		Uint32Ptr:    []*uint32{&ui321, &ui322},
		Uint64:       []uint64{ui641, ui642},
		Uint64Ptr:    []*uint64{&ui641, &ui642},
		NestedInt:    [][]int{{i1}},
		NestedIntPtr: [][]*int{nil, {&i1}},
	}
}

func BenchmarkComplexArrayDecodeStructAllTypes(b *testing.B) {
	values := getComplexArrayStructValues()
	decoder := form.NewDecoder()

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		var test ComplexArrayStruct
		if err := decoder.Decode(&test, values); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkComplexArrayDecodeStructAllTypesParallel(b *testing.B) {
	values := getComplexArrayStructValues()
	decoder := form.NewDecoder()

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var test ComplexArrayStruct
			if err := decoder.Decode(&test, values); err != nil {
				b.Error(err)
			}
		}
	})
}

func BenchmarkComplexArrayEncodeStructAllTypes(b *testing.B) {
	test := getComplexArrayStruct()
	encoder := form.NewEncoder()

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		if _, err := encoder.Encode(&test); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkComplexArrayEncodeStructAllTypesParallel(b *testing.B) {
	test := getComplexArrayStruct()
	encoder := form.NewEncoder()

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := encoder.Encode(&test); err != nil {
				b.Error(err)
			}
		}
	})
}

// Complex Map ALL types

type ComplexMapStruct struct {
	String       map[string]string
	StringPtr    map[*string]*string
	Int          map[int]int
	IntPtr       map[*int]*int
	Int8         map[int8]int8
	Int8Ptr      map[*int8]*int8
	Int16        map[int16]int16
	Int16Ptr     map[*int16]*int16
	Int32        map[int32]int32
	Int32Ptr     map[*int32]*int32
	Int64        map[int64]int64
	Int64Ptr     map[*int64]*int64
	Uint         map[uint]uint
	UintPtr      map[*uint]*uint
	Uint8        map[uint8]uint8
	Uint8Ptr     map[*uint8]*uint8
	Uint16       map[uint16]uint16
	Uint16Ptr    map[*uint16]*uint16
	Uint32       map[uint32]uint32
	Uint32Ptr    map[*uint32]*uint32
	Uint64       map[uint64]uint64
	Uint64Ptr    map[*uint64]*uint64
	NestedInt    map[int]map[int]int
	NestedIntPtr map[*int]map[*int]*int
}

func getComplexMapStructValues() url.Values {
	return url.Values{
		"String[key]":        []string{"value"},
		"StringPtr[key]":     []string{"value"},
		"Int[0]":             []string{"1"},
		"IntPtr[0]":          []string{"1"},
		"Int8[0]":            []string{"1"},
		"Int8Ptr[0]":         []string{"1"},
		"Int16[0]":           []string{"1"},
		"Int16Ptr[0]":        []string{"1"},
		"Int32[0]":           []string{"1"},
		"Int32Ptr[0]":        []string{"1"},
		"Int64[0]":           []string{"1"},
		"Int64Ptr[0]":        []string{"1"},
		"Uint[0]":            []string{"1"},
		"UintPtr[0]":         []string{"1"},
		"Uint8[0]":           []string{"1"},
		"Uint8Ptr[0]":        []string{"1"},
		"Uint16[0]":          []string{"1"},
		"Uint16Ptr[0]":       []string{"1"},
		"Uint32[0]":          []string{"1"},
		"Uint32Ptr[0]":       []string{"1"},
		"Uint64[0]":          []string{"1"},
		"Uint64Ptr[0]":       []string{"1"},
		"NestedInt[1][2]":    []string{"3"},
		"NestedIntPtr[1][2]": []string{"3"},
	}
}

func getComplexMapStruct() *ComplexMapStruct {

	key := "key"
	val := "value"

	i0 := int(0)
	i1 := int(1)
	i2 := int(2)
	i3 := int(3)
	i80 := int8(0)
	i81 := int8(1)
	i160 := int16(0)
	i161 := int16(1)
	i320 := int32(0)
	i321 := int32(1)
	i640 := int64(0)
	i641 := int64(1)

	ui0 := uint(0)
	ui1 := uint(1)
	ui80 := uint8(0)
	ui81 := uint8(1)
	ui160 := uint16(0)
	ui161 := uint16(1)
	ui320 := uint32(0)
	ui321 := uint32(1)
	ui640 := uint64(0)
	ui641 := uint64(1)

	return &ComplexMapStruct{
		String:       map[string]string{key: val},
		StringPtr:    map[*string]*string{&key: &val},
		Int:          map[int]int{i0: i1},
		IntPtr:       map[*int]*int{&i0: &i1},
		Int8:         map[int8]int8{i80: i81},
		Int8Ptr:      map[*int8]*int8{&i80: &i81},
		Int16:        map[int16]int16{i160: i161},
		Int16Ptr:     map[*int16]*int16{&i160: &i161},
		Int32:        map[int32]int32{i320: i321},
		Int32Ptr:     map[*int32]*int32{&i320: &i321},
		Int64:        map[int64]int64{i640: i641},
		Int64Ptr:     map[*int64]*int64{&i640: &i641},
		Uint:         map[uint]uint{ui0: ui1},
		UintPtr:      map[*uint]*uint{&ui0: &ui1},
		Uint8:        map[uint8]uint8{ui80: ui81},
		Uint8Ptr:     map[*uint8]*uint8{&ui80: &ui81},
		Uint16:       map[uint16]uint16{ui160: ui161},
		Uint16Ptr:    map[*uint16]*uint16{&ui160: &ui161},
		Uint32:       map[uint32]uint32{ui320: ui321},
		Uint32Ptr:    map[*uint32]*uint32{&ui320: &ui321},
		Uint64:       map[uint64]uint64{ui640: ui641},
		Uint64Ptr:    map[*uint64]*uint64{&ui640: &ui641},
		NestedInt:    map[int]map[int]int{i1: {i2: i3}},
		NestedIntPtr: map[*int]map[*int]*int{&i1: {&i2: &i3}},
	}
}

func BenchmarkComplexMapDecodeStructAllTypes(b *testing.B) {
	values := getComplexMapStructValues()
	decoder := form.NewDecoder()

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		var test ComplexMapStruct
		if err := decoder.Decode(&test, values); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkComplexMapDecodeStructAllTypesParallel(b *testing.B) {
	values := getComplexMapStructValues()
	decoder := form.NewDecoder()

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var test ComplexMapStruct
			if err := decoder.Decode(&test, values); err != nil {
				b.Error(err)
			}
		}
	})
}

func BenchmarkComplexMapEncodeStructAllTypes(b *testing.B) {
	test := getComplexMapStruct()
	encoder := form.NewEncoder()

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		if _, err := encoder.Encode(&test); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkComplexMapEncodeStructAllTypesParallel(b *testing.B) {
	test := getComplexMapStruct()
	encoder := form.NewEncoder()

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := encoder.Encode(&test); err != nil {
				b.Error(err)
			}
		}
	})
}

// NestedStruct Benchmarks

type Nested2 struct {
	Value   string
	Nested2 *Nested2
}

type Nested struct {
	Value string
}

type NestedStruct struct {
	Nested
	NestedArray    []Nested
	NestedPtrArray []*Nested
	Nested2        Nested2
}

func getNestedStructValues() url.Values {
	return url.Values{
		// Nested Field
		"Value": []string{"value"},
		// Nested Array
		"NestedArray[0].Value": []string{"value"},
		"NestedArray[1].Value": []string{"value"},
		// Nested Array Ptr
		"NestedPtrArray[0].Value": []string{"value"},
		"NestedPtrArray[1].Value": []string{"value"},
		// Nested 2
		"Nested2.Value":         []string{"value"},
		"Nested2.Nested2.Value": []string{"value"},
	}
}

func getNestedStruct() *NestedStruct {

	nested := Nested{
		Value: "value",
	}

	nested2 := Nested2{
		Value:   "value",
		Nested2: &Nested2{Value: "value"},
	}

	return &NestedStruct{
		Nested:         nested,
		NestedArray:    []Nested{nested, nested},
		NestedPtrArray: []*Nested{&nested, &nested},
		Nested2:        nested2,
	}
}

func BenchmarkDecodeNestedStruct(b *testing.B) {

	values := getNestedStructValues()
	decoder := form.NewDecoder()

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		var test NestedStruct
		if err := decoder.Decode(&test, values); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkDecodeNestedStructParallel(b *testing.B) {

	values := getNestedStructValues()
	decoder := form.NewDecoder()

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var test NestedStruct
			if err := decoder.Decode(&test, values); err != nil {
				b.Error(err)
			}
		}
	})
}

func BenchmarkEncodeNestedStruct(b *testing.B) {

	test := getNestedStruct()
	encoder := form.NewEncoder()

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		if _, err := encoder.Encode(&test); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkEncodeNestedStructParallel(b *testing.B) {

	test := getNestedStruct()
	encoder := form.NewEncoder()

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := encoder.Encode(&test); err != nil {
				b.Error(err)
			}
		}
	})
}

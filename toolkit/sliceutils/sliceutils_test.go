package sliceutils

import (
	"github.com/stretchr/testify/assert"
	"github.com/unionj-cloud/go-doudou/toolkit/cast"
	"reflect"
	"testing"
)

func TestContains(t *testing.T) {
	type args struct {
		src  []interface{}
		test interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "1",
			args: args{
				src:  []interface{}{"a", "2", "c"},
				test: "2",
			},
			want: true,
		},
		{
			name: "1",
			args: args{
				src:  []interface{}{"a", "2", "c"},
				test: "3",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.args.src, tt.args.test); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterfaceSlice2StringSlice(t *testing.T) {
	type args struct {
		strSlice []interface{}
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "2",
			args: args{
				strSlice: []interface{}{"a", "n"},
			},
			want: []string{"a", "n"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InterfaceSlice2StringSlice(tt.args.strSlice); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InterfaceSlice2StringSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndexOfAny(t *testing.T) {
	type args struct {
		target   interface{}
		anySlice interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "",
			args: args{
				target: "a",
				anySlice: []string{
					"b", "m", "a", "K",
				},
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "",
			args: args{
				target:   "a",
				anySlice: "abc",
			},
			want:    -1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IndexOfAny(tt.args.target, tt.args.anySlice)
			if (err != nil) != tt.wantErr {
				t.Errorf("IndexOfAny() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IndexOfAny() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndexOfAnyInt(t *testing.T) {
	type args struct {
		target   interface{}
		anySlice interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "4",
			args: args{
				target: 3,
				anySlice: []int{
					2, 5, 1, 6, 3, 8, 9,
				},
			},
			want:    4,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IndexOfAny(tt.args.target, tt.args.anySlice)
			if (err != nil) != tt.wantErr {
				t.Errorf("IndexOfAny() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IndexOfAny() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndexOfAnyNotContain(t *testing.T) {
	type args struct {
		target   interface{}
		anySlice interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "4",
			args: args{
				target: 3,
				anySlice: []int{
					2, 5, 1, 6, 11, 8, 9,
				},
			},
			want:    -1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IndexOfAny(tt.args.target, tt.args.anySlice)
			if (err != nil) != tt.wantErr {
				t.Errorf("IndexOfAny() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IndexOfAny() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertAny2Interface(t *testing.T) {
	type args struct {
		src interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []interface{}
		wantErr bool
	}{
		{
			name: "",
			args: args{
				src: []int{1, 2, 3},
			},
			want:    []interface{}{1, 2, 3},
			wantErr: false,
		},
		{
			name: "",
			args: args{
				src: "abc",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertAny2Interface(tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertAny2Interface() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertAny2Interface() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// 此为 TestConvertAny2Interface 的单元测试，同时也在演示要如何使用此函数
// This is a unit test for TestConvertAny2Interface, and also demonstrates how to use this function
func TestConvertAny2InterfaceForExplainUsage(t *testing.T) {
	// 说明 explain
	// 在 go-dou-dou 社区中，有人提问，如何把 slice 传入有动态类型的函数？
	// In the go-dou-dou community, there are people who ask questions, how to pass a slice to a function with dynamic parameter?
	t.Run("To insert a dynamically sized slice as a variadic parameter.", func(t *testing.T) {
		// DstFunc 函数是一个接受一个可变参数的函数
		// DstFunc is an example of a function with a variadic parameter.
		dstFunc := func(values ...interface{}) int {
			return len(values)
		}
		// 先测试 dstFunc 的功能
		// test dstFunc in the following
		count := dstFunc(1, 2, 3)
		if dstFunc(1, 2, 3) != 3 {
			t.Errorf("dstFunc() error, got = %d, want 3", count)
			return
		}
		count = dstFunc(1, 2, 3, 4, 5)
		if dstFunc(1, 2, 3) != 3 {
			t.Errorf("dstFunc() error, got = %d, want 5", count)
			return
		}
		// 再测试 ConvertAny2Interface 的功能
		// test ConvertAny2Interface in the following
		input, err := ConvertAny2Interface([]int{1, 2, 3})
		if err != nil {
			t.Errorf("ConvertAny2Interface() error = %v", err)
			return
		}
		// 错误的输入，因为没有加 ...
		// wrong input, because there is no ...
		if dstFunc(input) != 1 {
			t.Errorf("dstFunc() error, got = %d, want 1", count)
			return
		}
		// 正确的输入，因为加了 ...
		// correct input, because there is ...
		if dstFunc(input...) != 3 {
			t.Errorf("dstFunc() error, got = %d, want 3", count)
			return
		}
		input, err = ConvertAny2Interface([]int{1, 2, 3, 4, 5})
		if err != nil {
			t.Errorf("ConvertAny2Interface() error = %v", err)
			return
		}
		if dstFunc(input...) != 5 {
			t.Errorf("dstFunc() error, got = %d, want 5", count)
			return
		}
	})
}

func TestStringSlice2InterfaceSlice(t *testing.T) {
	type args struct {
		strSlice []string
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{
			name: "",
			args: args{
				strSlice: []string{"a", "b", "c"},
			},
			want: []interface{}{"a", "b", "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringSlice2InterfaceSlice(tt.args.strSlice); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringSlice2InterfaceSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

type containsDeepStruct struct {
	Name string
	Age  int
}

func TestContainsDeep(t *testing.T) {
	type args struct {
		src  []interface{}
		test interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "1",
			args: args{
				src: []interface{}{
					containsDeepStruct{
						Name: "Jack",
						Age:  10,
					},
					containsDeepStruct{
						Name: "Rose",
						Age:  18,
					},
					containsDeepStruct{
						Name: "David",
						Age:  14,
					},
				},
				test: containsDeepStruct{
					Name: "David",
					Age:  14,
				},
			},
			want: true,
		},
		{
			name: "2",
			args: args{
				src: []interface{}{
					containsDeepStruct{
						Name: "Jack",
						Age:  10,
					},
					containsDeepStruct{
						Name: "Rose",
						Age:  18,
					},
					containsDeepStruct{
						Name: "David",
						Age:  14,
					},
				},
				test: containsDeepStruct{
					Name: "Lily",
					Age:  14,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsDeep(tt.args.src, tt.args.test); got != tt.want {
				t.Errorf("ContainsDeep() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringContains(t *testing.T) {
	type args struct {
		src  []string
		test string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "1",
			args: args{
				src:  []string{"a", "b", "c"},
				test: "b",
			},
			want: true,
		},
		{
			name: "2",
			args: args{
				src:  []string{"a", "b", "c"},
				test: "d",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringContains(tt.args.src, tt.args.test); got != tt.want {
				t.Errorf("StringContains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndexOf(t *testing.T) {
	type args struct {
		element string
		data    []string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "",
			args: args{
				element: "c",
				data:    []string{"a", "b", "c"},
			},
			want: 2,
		},
		{
			name: "",
			args: args{
				element: "d",
				data:    []string{"a", "b", "c"},
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IndexOf(tt.args.element, tt.args.data); got != tt.want {
				t.Errorf("IndexOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsEmpty(t *testing.T) {
	type args struct {
		src interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "",
			args: args{
				src: []interface{}{},
			},
			want: true,
		},
		{
			name: "",
			args: args{
				src: []interface{}{1},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmpty(tt.args.src); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsEmptyP(t *testing.T) {
	assert.Panics(t, func() {
		IsEmpty(1)
	})
}

func TestStringFilter(t *testing.T) {
	type args struct {
		src []string
		fn  func(item string) bool
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "",
			args: args{
				src: []string{"a", "2", "c"},
				fn: func(item string) bool {
					_, err := cast.ToIntE(item)
					return err == nil
				},
			},
			want: []string{"2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, StringFilter(tt.args.src, tt.args.fn), "StringFilter(%v, %v)", tt.args.src, tt.args.fn)
		})
	}
}

package sliceutils

import (
	"github.com/stretchr/testify/assert"
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
			name: "3",
			args: args{
				target: "a",
				anySlice: []string{
					"b", "m", "a", "K",
				},
			},
			want:    2,
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

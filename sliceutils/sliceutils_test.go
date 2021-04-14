package sliceutils

import (
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

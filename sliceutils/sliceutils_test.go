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

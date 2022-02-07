package v3

import "testing"

func Test_isSupport(t *testing.T) {
	type args struct {
		t string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "1",
			args: args{
				t: "float32",
			},
			want: true,
		},
		{
			name: "2",
			args: args{
				t: "[]int64",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSupport(tt.args.t); got != tt.want {
				t.Errorf("isSupport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_castFunc(t *testing.T) {
	type args struct {
		t string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				t: "uint64",
			},
			want: "ToUint64",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CastFunc(tt.args.t); got != tt.want {
				t.Errorf("castFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

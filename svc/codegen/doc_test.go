package codegen

import (
	"github.com/unionj-cloud/go-doudou/pathutils"
	"testing"
)

func Test_schemasOf(t *testing.T) {
	type args struct {
		vofile string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Test_schemasOf",
			args: args{
				vofile: pathutils.Abs("testfiles") + "/vo.go",
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := schemasOf(tt.args.vofile); len(got) != tt.want {
				t.Errorf("schemasOf() = %v, want %v", len(got), tt.want)
			}
		})
	}
}

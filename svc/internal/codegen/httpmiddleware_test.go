package codegen

import (
	"os"
	"testing"
)

func TestGenRouterMiddleware(t *testing.T) {
	dir := testDir + "middleware1"
	InitSvc(dir)
	defer os.RemoveAll(dir)
	type args struct {
		dir string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "1",
			args: args{
				dir: dir,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GenHttpMiddleware(tt.args.dir)
		})
	}
}
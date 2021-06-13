package tests

import (
	"github.com/unionj-cloud/go-doudou/svc"
	. "github.com/unionj-cloud/go-doudou/svc/codegen"
	"os"
	"testing"
)

func TestGenRouterMiddleware(t *testing.T) {
	dir := testDir + "middleware1"
	receiver := svc.Svc{
		Dir: dir,
	}
	receiver.Init()
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

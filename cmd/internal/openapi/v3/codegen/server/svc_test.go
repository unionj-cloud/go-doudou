package server

import (
	"path/filepath"
	"testing"
)

var dir = "../testdata"

func TestGenSvcGo(t *testing.T) {
	type args struct {
		dir     string
		docPath string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				dir:     filepath.Join(dir, "testgensvcgo"),
				docPath: filepath.Join(dir, "prometheus_openapi3.json"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GenSvcGo(tt.args.dir, tt.args.docPath)
		})
	}
}

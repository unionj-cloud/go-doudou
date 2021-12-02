package codegen

import (
	"github.com/unionj-cloud/go-doudou/pathutils"
	"os"
	"path/filepath"
	"testing"
)

func TestGenK8sStatefulset(t *testing.T) {
	type args struct {
		dir     string
		svcname string
		image   string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				dir:     pathutils.Abs("./testdata"),
				svcname: "corpus",
				image:   "google.com/corpus:v2.0.0",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GenK8sStatefulset(tt.args.dir, tt.args.svcname, tt.args.image)
		})
	}
}

func TestGenK8sStatefulset2(t *testing.T) {
	defer os.Remove(filepath.Join("testdata", "nostatefulset", "corpus_statefulset.yaml"))
	type args struct {
		dir     string
		svcname string
		image   string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				dir:     filepath.Join("testdata", "nostatefulset"),
				svcname: "corpus",
				image:   "google.com/corpus:v2.0.0",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GenK8sStatefulset(tt.args.dir, tt.args.svcname, tt.args.image)
		})
	}
}

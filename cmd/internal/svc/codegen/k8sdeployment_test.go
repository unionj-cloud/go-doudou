package codegen

import (
	"fmt"
	"github.com/unionj-cloud/go-doudou/toolkit/pathutils"
	"os"
	"path/filepath"
	"testing"
)

func TestModifyVersion(t *testing.T) {
	type args struct {
		yfile string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				yfile: pathutils.Abs("./testdata/k8s.yaml"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := modifyVersion(tt.args.yfile, "v1.0.0")
			fmt.Println(string(result))
		})
	}
}

func TestModifyVersion2(t *testing.T) {
	type args struct {
		yfile string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				yfile: pathutils.Abs("./testdata/corpus_statefulset.yaml"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := modifyVersion(tt.args.yfile, "v1.0.0")
			fmt.Println(string(result))
		})
	}
}

func TestGenK8sDeployment(t *testing.T) {
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
			GenK8sDeployment(tt.args.dir, tt.args.svcname, tt.args.image)
		})
	}
}

func TestGenK8sDeployment2(t *testing.T) {
	os.MkdirAll(filepath.Join("testdata", "nodeployment"), os.ModePerm)
	defer os.RemoveAll(filepath.Join("testdata", "nodeployment"))
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
				dir:     filepath.Join("testdata", "nodeployment"),
				svcname: "corpus",
				image:   "google.com/corpus:v2.0.0",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GenK8sDeployment(tt.args.dir, tt.args.svcname, tt.args.image)
		})
	}
}

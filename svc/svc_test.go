package svc

import (
	"github.com/stretchr/testify/assert"
	"github.com/unionj-cloud/go-doudou/astutils"
	"github.com/unionj-cloud/go-doudou/pathutils"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var testDir string

func init() {
	testDir = pathutils.Abs("testdata")
}

func TestSvc_Create(t *testing.T) {
	type fields struct {
		Dir string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "1",
			fields: fields{
				Dir: testDir + "1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := Svc{
				Dir: tt.fields.Dir,
			}
			receiver.Init()
			defer os.RemoveAll(tt.fields.Dir)
		})
	}
}

func TestSvc_Http(t *testing.T) {
	type fields struct {
		Dir          string
		Handler      bool
		Client       string
		Omitempty    bool
		Doc          bool
		Jsonattrcase string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "",
			fields: fields{
				Dir:          testDir + "2",
				Handler:      true,
				Client:       "go",
				Omitempty:    true,
				Doc:          true,
				Jsonattrcase: "snake",
			},
		},
		{
			name: "",
			fields: fields{
				Dir:       testDir + "3",
				Handler:   true,
				Client:    "go",
				Omitempty: false,
				Doc:       false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := Svc{
				Dir:          tt.fields.Dir,
				Handler:      tt.fields.Handler,
				Client:       tt.fields.Client,
				Omitempty:    tt.fields.Omitempty,
				Doc:          tt.fields.Doc,
				Jsonattrcase: tt.fields.Jsonattrcase,
			}
			assert.NotPanics(t, func() {
				receiver.Init()
			})
			defer os.RemoveAll(tt.fields.Dir)
			assert.NotPanics(t, func() {
				receiver.Http()
			})
		})
	}
}

func Test_checkIc(t *testing.T) {
	svcfile := testDir + "/svc.go"
	ic := astutils.BuildInterfaceCollector(svcfile, astutils.ExprString)
	type args struct {
		ic astutils.InterfaceCollector
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				ic: ic,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				validateRestApi(ic)
			})
		})
	}
}

func Test_checkIc1(t *testing.T) {
	svcfile := testDir + "/svcp.go"
	ic := astutils.BuildInterfaceCollector(svcfile, astutils.ExprString)
	type args struct {
		ic astutils.InterfaceCollector
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "",
			args: args{
				ic: ic,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Panics(t, func() {
				validateRestApi(ic)
			})
		})
	}
}

func TestSvc_Deploy(t *testing.T) {
	dir := testDir + "deploy"
	receiver := Svc{
		Dir: dir,
	}
	receiver.Init()
	defer os.RemoveAll(dir)
	err := os.Chdir(dir)
	if err != nil {
		t.Fatal(err)
	}
	assert.Panics(t, func() {
		receiver.Deploy()
	})
}

func TestSvc_Shutdown(t *testing.T) {
	dir := testDir + "shutdown"
	receiver := Svc{
		Dir: dir,
	}
	receiver.Init()
	defer os.RemoveAll(dir)
	err := os.Chdir(dir)
	if err != nil {
		t.Fatal(err)
	}
	assert.Panics(t, func() {
		receiver.Shutdown()
	})
}

func Test_validateDataType(t *testing.T) {
	assert.NotPanics(t, func() {
		validateDataType(testDir)
	})
}

func Test_validateDataType_shouldpanic(t *testing.T) {
	assert.Panics(t, func() {
		validateDataType(pathutils.Abs("testdata1"))
	})
}

func Test_GenClient(t *testing.T) {
	defer os.RemoveAll(filepath.Join(testDir, "client"))
	s := Svc{
		Dir:       testDir,
		DocPath:   filepath.Join(testDir, "testfilesdoc1_openapi3.json"),
		Omitempty: true,
		Client:    "go",
		ClientPkg: "client",
	}
	assert.NotPanics(t, func() {
		s.GenClient()
	})
}

func TestSvc_Seed(t *testing.T) {
	assert.NotPanics(t, func() {
		s := Svc{}
		go s.Seed()
		time.Sleep(2 * time.Second)
	})
}

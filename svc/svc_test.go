package svc

import (
	"fmt"
	"github.com/radovskyb/watcher"
	"github.com/stretchr/testify/assert"
	"github.com/unionj-cloud/go-doudou/astutils"
	"github.com/unionj-cloud/go-doudou/executils"
	"github.com/unionj-cloud/go-doudou/pathutils"
	"os"
	"os/exec"
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
				dir: tt.fields.Dir,
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
		{
			name: "",
			fields: fields{
				Dir:       testDir + "4",
				Handler:   false,
				Client:    "go",
				Omitempty: false,
				Doc:       false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := Svc{
				dir:          tt.fields.Dir,
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

func Test_checkIc_no_interface(t *testing.T) {
	svcfile := filepath.Join(testDir, "nosvc", "svc.go")
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

func Test_checkIc_input_anonystruct(t *testing.T) {
	svcfile := filepath.Join(testDir, "inputanonystruct", "svc.go")
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

func Test_checkIc_output_anonystruct(t *testing.T) {
	svcfile := filepath.Join(testDir, "outputanonystruct", "svc.go")
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
	dir := testDir + "/deploy"
	receiver := NewMockSvc(dir)
	receiver.Init()
	defer os.RemoveAll(dir)
	os.Chdir(dir)
	receiver.Deploy("")
}

func TestSvc_Shutdown(t *testing.T) {
	dir := testDir + "/shutdown"
	receiver := NewMockSvc(dir)
	receiver.Init()
	defer os.RemoveAll(dir)
	os.Chdir(dir)
	receiver.Shutdown("")
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
		dir:       testDir,
		DocPath:   filepath.Join(testDir, "testfilesdoc1_openapi3.json"),
		Omitempty: true,
		Client:    "go",
		ClientPkg: "client",
	}
	assert.NotPanics(t, func() {
		s.GenClient()
	})
}

func TestSvc_Push(t *testing.T) {
	s := NewMockSvc(pathutils.Abs("./testdata"))
	s.Push("wubin1989")
}

func TestHelperProcess(*testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	defer os.Exit(0)
	fmt.Println("testing helper process")
}

func TestSvc_run(t *testing.T) {
	s := NewMockSvc(pathutils.Abs("./testdata"))
	s.run()
}

func TestSvc_restart(t *testing.T) {
	s := NewMockSvc(pathutils.Abs("./testdata"))
	s.cmd = s.run()
	s.restart()
}

func TestSvc_watch(t *testing.T) {
	s := NewMockSvc(pathutils.Abs("./testdata/change"))
	s.w = watcher.New()
	go s.watch()
	time.Sleep(1 * time.Second)
	f, _ := os.Create(filepath.Join(s.dir, "change.go"))
	defer f.Close()
	f.WriteString("test")
	time.Sleep(6 * time.Second)
	s.w.Close()
}

func TestSvc_Run(t *testing.T) {
	s := NewMockSvc(pathutils.Abs("./testdata/change"))
	s.w = watcher.New()
	defer s.w.Close()
	go s.Run(true)
	time.Sleep(1 * time.Second)
	f, _ := os.Create(filepath.Join(s.dir, "change.go"))
	defer f.Close()
	f.WriteString("test")
	time.Sleep(6 * time.Second)
}

func TestSvc_Run_unwatch(t *testing.T) {
	s := NewMockSvc("")
	s.Run(false)
}

func TestSvc_GenClient_DocPathEmpty2(t *testing.T) {
	type fields struct {
		dir                  string
		Handler              bool
		Client               string
		Omitempty            bool
		Doc                  bool
		Jsonattrcase         string
		DocPath              string
		Env                  string
		ClientPkg            string
		cmd                  *exec.Cmd
		restartSig           chan int
		RoutePatternStrategy int
		runner               executils.Runner
		w                    *watcher.Watcher
		ModName              string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "",
			fields: fields{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := Svc{}
			assert.Panics(t, func() {
				receiver.GenClient()
			})
		})
	}
}

func TestSvc_GenClient_DocPathEmpty1(t *testing.T) {
	os.Chdir("testdata")
	type fields struct {
		dir                  string
		Handler              bool
		Client               string
		Omitempty            bool
		Doc                  bool
		Jsonattrcase         string
		DocPath              string
		Env                  string
		ClientPkg            string
		cmd                  *exec.Cmd
		restartSig           chan int
		RoutePatternStrategy int
		runner               executils.Runner
		w                    *watcher.Watcher
		ModName              string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "",
			fields: fields{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := Svc{}
			assert.NotPanics(t, func() {
				receiver.GenClient()
			})
		})
	}
}

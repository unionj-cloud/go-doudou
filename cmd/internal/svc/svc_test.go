package svc_test

import (
	"fmt"
	"github.com/radovskyb/watcher"
	"github.com/stretchr/testify/assert"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/svc"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/svc/validate"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/astutils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/executils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/pathutils"
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

// NewMockSvc new Svc instance for unit test purpose
func NewMockSvc(dir string, opts ...svc.SvcOption) svc.ISvc {
	return svc.NewSvc(dir, svc.WithRunner(mockRunner{}))
}

type mockRunner struct {
}

func (r mockRunner) Output(command string, args ...string) ([]byte, error) {
	return []byte("go version go1.17.8 darwin/amd64"), nil
}

func (r mockRunner) Run(command string, args ...string) error {
	cs := []string{"-test.run=TestHelperProcess", "--"}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	return nil
}

func (r mockRunner) Start(command string, args ...string) (*exec.Cmd, error) {
	cs := []string{"-test.run=TestHelperProcess", "--"}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	return cmd, nil
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
			receiver := svc.NewSvc(tt.fields.Dir)
			receiver.Init()
			defer os.RemoveAll(tt.fields.Dir)
		})
	}
}

func TestSvc_Http(t *testing.T) {
	type fields struct {
		Dir          string
		Handler      bool
		Client       bool
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
				Client:       true,
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
				Client:    true,
				Omitempty: false,
				Doc:       false,
			},
		},
		{
			name: "",
			fields: fields{
				Dir:       testDir + "4",
				Handler:   false,
				Client:    true,
				Omitempty: false,
				Doc:       false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := svc.NewSvc(tt.fields.Dir)
			s := receiver.(*svc.Svc)
			s.Handler = tt.fields.Handler
			s.Client = tt.fields.Client
			s.Omitempty = tt.fields.Omitempty
			s.Doc = tt.fields.Doc
			s.JsonCase = tt.fields.Jsonattrcase
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
				validate.RestApi(testDir, ic)
			})
		})
	}
}

func Test_checkIc2(t *testing.T) {
	svcfile := filepath.Join(testDir, "checkIc2", "svc.go")
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
				validate.RestApi(testDir, ic)
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
				validate.RestApi(testDir, ic)
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
				validate.RestApi(testDir, ic)
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
				validate.RestApi(testDir, ic)
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
				validate.RestApi(testDir, ic)
			})
		})
	}
}

func TestSvc_Deploy(t *testing.T) {
	dir := testDir + "/deploy"
	receiver := NewMockSvc(dir)
	receiver.Init()
	defer os.RemoveAll(dir)
	receiver.Deploy("")
}

func TestSvc_Shutdown(t *testing.T) {
	dir := testDir + "/shutdown"
	receiver := NewMockSvc(dir)
	receiver.Init()
	defer os.RemoveAll(dir)
	receiver.Shutdown("")
}

func Test_validateDataType(t *testing.T) {
	assert.NotPanics(t, func() {
		validate.DataType(testDir)
	})
}

func Test_validateDataType_shouldpanic(t *testing.T) {
	assert.Panics(t, func() {
		validate.DataType(pathutils.Abs("testdata1"))
	})
}

func Test_GenClient(t *testing.T) {
	defer os.RemoveAll(filepath.Join(testDir, "client"))
	receiver := svc.NewSvc(testDir)
	s := receiver.(*svc.Svc)
	s.DocPath = filepath.Join(testDir, "testfilesdoc1_openapi3.json")
	s.ClientPkg = "client"
	s.Omitempty = true
	assert.NotPanics(t, func() {
		s.GenClient()
	})
}

func TestSvc_Push(t *testing.T) {
	receiver := NewMockSvc(pathutils.Abs("./testdata"))
	s := receiver.(*svc.Svc)
	s.Push(svc.PushConfig{
		Repo:   "wubin1989",
		Prefix: "go-doudou-",
	})
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
	s.DoRun()
}

func TestSvc_restart(t *testing.T) {
	s := NewMockSvc(pathutils.Abs("./testdata"))
	s.DoRun()
	s.DoRestart()
}

func TestSvc_watch(t *testing.T) {
	s := NewMockSvc(pathutils.Abs("./testdata/change"))
	s.SetWatcher(watcher.New())
	go s.DoWatch()
	time.Sleep(1 * time.Second)
	f, _ := os.Create(filepath.Join(s.GetDir(), "change.go"))
	defer f.Close()
	f.WriteString("test")
	time.Sleep(6 * time.Second)
	s.GetWatcher().Close()
}

func TestSvc_Run(t *testing.T) {
	s := NewMockSvc(pathutils.Abs("./testdata/change"))
	s.SetWatcher(watcher.New())
	defer s.GetWatcher().Close()
	go s.Run(true)
	time.Sleep(1 * time.Second)
	f, _ := os.Create(filepath.Join(s.GetDir(), "change.go"))
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
			receiver := svc.Svc{}
			assert.Panics(t, func() {
				receiver.GenClient()
			})
		})
	}
}

func TestSvc_GenClient_DocPathEmpty1(t *testing.T) {
	defer os.RemoveAll(filepath.Join(testDir, "openapi", "client"))
	receiver := svc.NewSvc(filepath.Join(testDir, "openapi"))
	s := receiver.(*svc.Svc)
	s.ClientPkg = "client"
	assert.NotPanics(t, func() {
		receiver.GenClient()
	})
}

func TestNewSvc(t *testing.T) {
	assert.NotPanics(t, func() {
		svc.NewSvc("")
	})
}

package svc

import (
	"fmt"
	"github.com/radovskyb/watcher"
	"github.com/stretchr/testify/assert"
	"github.com/unionj-cloud/go-doudou/astutils"
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
	dir := testDir + "/deploy"
	receiver := NewSvc()
	receiver.Init()
	defer os.RemoveAll(dir)
	err := os.Chdir(dir)
	if err != nil {
		t.Fatal(err)
	}
	receiver.Deploy("")
}

func ExampleSvc_Shutdown() {
	dir := testDir + "/shutdown"
	receiver := Svc{
		Dir:    dir,
		runner: MockRunner{},
	}
	receiver.Init()
	defer os.RemoveAll(dir)
	os.Chdir(dir)
	receiver.Shutdown("")
	// Output:
	// 1.16
	// shutdown
	// testing helper process
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

func ExampleSvc_Push() {
	s := Svc{
		runner: MockRunner{},
		Dir:    pathutils.Abs("./testdata"),
	}
	s.Push("wubin1989")
	// Output:
	// testing helper process
	// testing helper process
	// testing helper process
	// testing helper process
}

type MockRunner struct {
}

func (r MockRunner) Run(command string, args ...string) error {
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

func (r MockRunner) Start(command string, args ...string) (*exec.Cmd, error) {
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

func TestHelperProcess(*testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	defer os.Exit(0)
	fmt.Println("testing helper process")
}

func ExampleSvc_run() {
	s := Svc{
		runner: MockRunner{},
		Dir:    pathutils.Abs("./testdata"),
	}
	s.run()
	// Output:
	// testing helper process
	// testing helper process
}

func ExampleSvc_restart() {
	s := Svc{
		runner: MockRunner{},
		Dir:    pathutils.Abs("./testdata"),
	}
	s.cmd = s.run()
	s.restart()
	// Output:
	// testing helper process
	// testing helper process
	// testing helper process
}

func ExampleSvc_watch() {
	s := Svc{
		runner: MockRunner{},
		w:      watcher.New(),
		Dir:    pathutils.Abs("./testdata/change"),
	}
	go s.watch()
	time.Sleep(1 * time.Second)
	f, _ := os.Create(filepath.Join(s.Dir, "change.go"))
	defer f.Close()
	f.WriteString("test")
	time.Sleep(6 * time.Second)
	s.w.Close()
	// Output:
	// FILE "change.go" WRITE [/Users/wubin1989/workspace/cloud/go-doudou/svc/testdata/change/change.go]
	// testing helper process
}

func ExampleSvc_Run() {
	s := Svc{
		runner: MockRunner{},
		w:      watcher.New(),
		Dir:    pathutils.Abs("./testdata/change"),
	}
	defer s.w.Close()
	go s.Run(true)
	time.Sleep(1 * time.Second)
	f, _ := os.Create(filepath.Join(s.Dir, "change.go"))
	defer f.Close()
	f.WriteString("test")
	time.Sleep(6 * time.Second)
	// Output:
	// testing helper process
	// testing helper process
	// FILE "change.go" WRITE [/Users/wubin1989/workspace/cloud/go-doudou/svc/testdata/change/change.go]
	// testing helper process
}

func ExampleSvc_Run_unwatch() {
	s := Svc{
		runner: MockRunner{},
	}
	s.Run(false)
	// Output:
	// testing helper process
	// testing helper process
}

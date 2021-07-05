package svc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic"
	"github.com/stretchr/testify/assert"
	"github.com/unionj-cloud/go-doudou/astutils"
	"github.com/unionj-cloud/go-doudou/esutils"
	"github.com/unionj-cloud/go-doudou/logutils"
	"github.com/unionj-cloud/go-doudou/pathutils"
	"github.com/unionj-cloud/go-doudou/test"
	"os"
	"testing"
)

var testDir string

func init() {
	testDir = pathutils.Abs("testfiles")
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

func TestSvc_Publish(t *testing.T) {
	terminator, host, port := test.PrepareTestEnvironment()
	defer terminator()

	esclient, err := elastic.NewSimpleClient(
		elastic.SetErrorLog(logutils.NewLogger()),
		elastic.SetURL([]string{fmt.Sprintf("http://%s:%d", host, port)}...),
		elastic.SetGzip(true),
	)
	if err != nil {
		t.Errorf("call NewSimpleClient() error: %+v\n", err)
	}
	es := esutils.NewEs("doc", "", esutils.WithClient(esclient))

	svc := Svc{
		DocPath: testDir + "/testfilesdoc1_openapi3.json",
		Es:      es,
	}

	doc, err := es.GetByID(context.Background(), svc.Publish())
	if err != nil {
		t.Error(err)
	}
	jsonrs, _ := json.MarshalIndent(doc, "", "  ")
	fmt.Println(string(jsonrs))
}

func TestSvc_Push(t *testing.T) {
	dir := testDir + "push"
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
		receiver.Push()
	})
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

func TestSvc_Scale(t *testing.T) {
	dir := testDir + "scale"
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
		receiver.Scale()
	})
}

func Test_validateDataType(t *testing.T) {
	assert.NotPanics(t, func() {
		validateDataType(testDir)
	})
}

func Test_validateDataType_shouldpanic(t *testing.T) {
	assert.Panics(t, func() {
		validateDataType(pathutils.Abs("testfiles1"))
	})
}

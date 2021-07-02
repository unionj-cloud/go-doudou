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
			defer os.RemoveAll(testDir + "1")
		})
	}
}

func TestSvc_Http(t *testing.T) {
	type fields struct {
		Dir string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "2",
			fields: fields{
				Dir: testDir + "2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := Svc{
				Dir: tt.fields.Dir,
			}
			receiver.Init()
			defer os.RemoveAll(testDir + "2")
			receiver.Http()
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
	esaddr, terminator := prepareTestEnvironment()
	defer terminator()

	esclient, err := elastic.NewSimpleClient(
		elastic.SetErrorLog(logutils.NewLogger()),
		elastic.SetURL([]string{esaddr}...),
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

package svc

import (
	"github.com/unionj-cloud/go-doudou/astutils"
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
	receiver := Svc{
		Dir: testDir + "3",
	}
	receiver.Init()
	defer os.RemoveAll(testDir + "3")
	svcfile := testDir + "3" + "/svc.go"
	ic := buildIc(svcfile)
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
			checkIc(ic)
		})
	}
}

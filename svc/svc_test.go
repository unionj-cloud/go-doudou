package svc

import (
	"fmt"
	"github.com/unionj-cloud/go-doudou/astutils"
	"testing"
)

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
				Dir: "/Users/wubin1989/workspace/cloud/comment-svc",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := Svc{
				Dir: tt.fields.Dir,
			}
			receiver.Init()
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
				Dir: "/Users/wubin1989/workspace/cloud/comment-svc",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := Svc{
				Dir: tt.fields.Dir,
			}
			receiver.Http()
		})
	}
}

func ExampleParseInterface() {
	svcfile := "/Users/wubin1989/workspace/cloud/comment-svc/svc.go"
	ic := buildIc(svcfile)
	fmt.Printf("%+v\n", ic)
	// Output:

}

func Test_checkIc(t *testing.T) {
	svcfile := "/Users/wubin1989/workspace/cloud/ordersvc/svc.go"
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

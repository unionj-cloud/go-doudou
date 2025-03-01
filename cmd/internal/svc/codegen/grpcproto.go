package codegen

import (
	"fmt"
	"github.com/unionj-cloud/go-doudou/v2/version"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/toolkit/astutils"
	protov3 "github.com/unionj-cloud/toolkit/protobuf/v3"
	"github.com/unionj-cloud/toolkit/templateutils"
)

var protoTmpl = `/**
* Generated by go-doudou {{.Version}}.
* Don't edit!
*
* Version No.: {{.ProtoVer}}
*/
syntax = "proto3";

package {{.Package}};
option go_package = "{{.GoPackage}}";
{{ range $i := .Imports }}
import "{{$i}}";
{{- end }}
{{ range $e := .Enums }}
enum {{$e.Name}} {
  {{- range $f := $e.Fields }}
  {{$f.Name}} = {{$f.Number}};
  {{- end }}
}
{{- end }}

{{- range $m := .Messages }}
{{ Eval "Message" $m }}
{{- end }}

{{- define "Message"}}
{{- if .Comments }}
{{ toComment .Comments }}
{{- end }}
message {{.Name}} {
  {{- range $f := .Fields }}
  {{- if $f.Type.Inner}}
  {{ Eval "Message" $f.Type }}
  {{- end }}
  {{- if $f.Comments }}
  {{ toComment $f.Comments }}
  {{- end }}
  {{$f.Type.GetName | label}} {{$f.Name}} = {{$f.Number}}{{if $f.JsonName}} [json_name="{{$f.JsonName}}"]{{end}};
  {{- end }}
}
{{- end}}

service {{.Name}} {
  {{- range $r := .Rpcs }}
  {{- if $r.Comments }}
  {{ toComment $r.Comments }}
  {{- end }}
  rpc {{$r.Name}}({{$r.Request.Name}}) returns ({{$r.Response.Name}});
  {{- end}}
}
`

func toComment(comments []string) string {
	if len(comments) == 0 {
		return ""
	}
	var b strings.Builder
	for i := range comments {
		b.WriteString(fmt.Sprintf("// %s\n", comments[i]))
	}
	return strings.TrimSuffix(b.String(), "\n")
}

func GenGrpcProto(dir string, ic astutils.InterfaceCollector, p protov3.ProtoGenerator) (service protov3.Service, protoFile string) {
	var (
		err     error
		svcname string
		fi      os.FileInfo
		tpl     *template.Template
		f       *os.File
		grpcDir string
	)
	grpcDir = filepath.Join(dir, "transport", "grpc")
	if err = os.MkdirAll(grpcDir, os.ModePerm); err != nil {
		panic(err)
	}
	svcname = ic.Interfaces[0].Name
	protoFile = filepath.Join(grpcDir, strings.ToLower(svcname)+".proto")
	if fi, err = os.Stat(protoFile); err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	if fi != nil {
		logrus.Warningln("file " + protoFile + " will be overwritten")
	}
	if f, err = os.Create(protoFile); err != nil {
		panic(err)
	}
	defer f.Close()
	servicePkg := astutils.GetPkgPath(dir)
	service = p.NewService(svcname, servicePkg+"/transport/grpc", version.Release)
	service.Comments = ic.Interfaces[0].Comments
	for _, method := range ic.Interfaces[0].Methods {
		rpc := p.NewRpc(method)
		if rpc == nil {
			continue
		}
		service.Rpcs = append(service.Rpcs, *rpc)
	}
	for k := range protov3.ImportStore {
		service.Imports = append(service.Imports, k)
	}
	sort.Strings(service.Imports)
	for _, v := range protov3.MessageStore {
		service.Messages = append(service.Messages, v)
	}
	sort.SliceStable(service.Messages, func(i, j int) bool {
		return service.Messages[i].Name < service.Messages[j].Name
	})
	for _, v := range protov3.EnumStore {
		service.Enums = append(service.Enums, v)
	}
	sort.SliceStable(service.Enums, func(i, j int) bool {
		return service.Enums[i].Name < service.Enums[j].Name
	})
	tpl = template.New("proto.tmpl")
	funcMap := make(map[string]interface{})
	funcMap["toComment"] = toComment
	funcMap["Eval"] = templateutils.Eval(tpl)
	funcMap["label"] = func(input string) string {
		if !strings.Contains(input, "repeated ") && !strings.Contains(input, "map<") &&
			!strings.Contains(input, "required ") && !strings.Contains(input, "optional ") {
			return "optional " + input
		}
		return input
	}
	if tpl, err = tpl.Funcs(funcMap).Parse(protoTmpl); err != nil {
		panic(err)
	}
	if err = tpl.Execute(f, service); err != nil {
		panic(err)
	}
	return
}

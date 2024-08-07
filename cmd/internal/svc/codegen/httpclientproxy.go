package codegen

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/templates"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/astutils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/copier"
	v3helper "github.com/unionj-cloud/go-doudou/v2/toolkit/openapi/v3"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/zlogger"
	"github.com/unionj-cloud/go-doudou/v2/version"
)

var cpimportTmpl = `
	"context"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/zlogger"
	"github.com/slok/goresilience"
	"github.com/go-resty/resty/v2"
	"github.com/slok/goresilience/circuitbreaker"
	rerrors "github.com/slok/goresilience/errors"
	"github.com/slok/goresilience/metrics"
	"github.com/slok/goresilience/retry"
	"github.com/slok/goresilience/timeout"
	v3 "github.com/unionj-cloud/go-doudou/v2/toolkit/openapi/v3"
	"os"
	"time"
	"{{.DtoPackage}}"
`

var appendTmpl = `
{{- range $m := .Meta.Methods }}
	func (receiver *{{$.SvcName}}ClientProxy) {{$m.Name}}(ctx context.Context, _headers map[string]string, {{- range $i, $p := $m.Params}}
	{{- if ne $p.Type "context.Context" }}
	{{- $p.Name}} {{$p.Type}},
	{{- end }}
    {{- end }} options Options) (_resp *resty.Response, {{- range $i, $r := $m.Results}}
                     {{- if $i}},{{end}}
                     {{- $r.Name}} {{$r.Type}}
                     {{- end }}) {
		if _err := receiver.runner.Run(ctx, func(ctx context.Context) error {
			_resp, {{ range $i, $r := $m.Results }}{{- if $i}},{{- end}}{{- $r.Name }}{{- end }} = receiver.client.{{$m.Name}}(
				ctx,
				_headers,
				{{- range $p := $m.Params }}
				{{- if ne $p.Type "context.Context" }}
				{{- if isVarargs $p.Type }}
				{{ $p.Name }}...,
				{{- else }}
				{{ $p.Name }},
				{{- end }}
				{{- end }}
				{{- end }}
				options,
			)
			{{- range $r := $m.Results }}
				{{- if eq $r.Type "error" }}
					if {{ $r.Name }} != nil {
						return errors.Wrap({{ $r.Name }}, "call {{$m.Name}} fail")
					}
				{{- end }}
			{{- end }}
			return nil
		}); _err != nil {
			// you can implement your fallback logic here
			if errors.Is(_err, rerrors.ErrCircuitOpen) {
				receiver.logger.Error().Err(_err).Msg("")
			}
			{{- range $r := $m.Results }}
				{{- if eq $r.Type "error" }}
					{{ $r.Name }} = errors.Wrap(_err, "call {{$m.Name}} fail")
				{{- end }}
			{{- end }}
		}
		return
	}
{{- end }}
`

var baseTmpl = templates.EditableHeaderTmpl + `package client

import ()

type {{.SvcName}}ClientProxy struct {
	client *{{.SvcName}}Client
	logger zerolog.Logger
	runner goresilience.Runner
}

` + appendTmpl + `

type ProxyOption func(*{{.SvcName}}ClientProxy)

func WithRunner(runner goresilience.Runner) ProxyOption {
	return func(proxy *{{.SvcName}}ClientProxy) {
		proxy.runner = runner
	}
}

func WithLogger(logger zerolog.Logger) ProxyOption {
	return func(proxy *{{.SvcName}}ClientProxy) {
		proxy.logger = logger
	}
}

func New{{.SvcName}}ClientProxy(client *{{.SvcName}}Client, opts ...ProxyOption) *{{.SvcName}}ClientProxy {
	cp := &{{.SvcName}}ClientProxy{
		client: client,
		logger: zlogger.Logger,
	}

	for _, opt := range opts {
		opt(cp)
	}

	if cp.runner == nil {
		var mid []goresilience.Middleware
		mid = append(mid, metrics.NewMiddleware("{{.ServicePackage}}_client", metrics.NewPrometheusRecorder(prometheus.DefaultRegisterer)))
		mid = append(mid, circuitbreaker.NewMiddleware(circuitbreaker.Config{
			ErrorPercentThresholdToOpen:        50,
			MinimumRequestToOpen:               6,
			SuccessfulRequiredOnHalfOpen:       1,
			WaitDurationInOpenState:            5 * time.Second,
			MetricsSlidingWindowBucketQuantity: 10,
			MetricsBucketDuration:              1 * time.Second,
		}),
			timeout.NewMiddleware(timeout.Config{
				Timeout: 3 * time.Minute,
			}),
			retry.NewMiddleware(retry.Config{
				Times: 3,
			}))

		cp.runner = goresilience.RunnerChain(mid...)
	}

	return cp
}
`

func unimplementedSvcMethods(meta *astutils.InterfaceMeta, clientfile string) {
	fset := token.NewFileSet()
	root, err := parser.ParseFile(fset, clientfile, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	sc := astutils.NewStructCollector(astutils.ExprString)
	ast.Walk(sc, root)
	if handlers, exists := sc.Methods[meta.Name+"ClientProxy"]; exists {
		var notimplemented []astutils.MethodMeta
		for _, item := range meta.Methods {
			for _, handler := range handlers {
				if item.Name == handler.Name {
					goto L
				}
			}
			notimplemented = append(notimplemented, item)

		L:
		}

		meta.Methods = notimplemented
	}
}

// GenGoClientProxy wraps client with resiliency features
func GenGoClientProxy(dir string, ic astutils.InterfaceCollector) {
	var (
		err             error
		clientfile      string
		f               *os.File
		tpl             *template.Template
		buf             bytes.Buffer
		clientDir       string
		fi              os.FileInfo
		meta            astutils.InterfaceMeta
		clientProxyTmpl string
		importBuf       bytes.Buffer
	)
	clientDir = filepath.Join(dir, "client")
	if err = os.MkdirAll(clientDir, os.ModePerm); err != nil {
		panic(err)
	}

	clientfile = filepath.Join(clientDir, "clientproxy.go")
	fi, err = os.Stat(clientfile)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	err = copier.DeepCopy(ic.Interfaces[0], &meta)
	if err != nil {
		panic(err)
	}
	if fi != nil {
		zlogger.Warn().Msg("New content will be append to clientproxy.go file")
		if f, err = os.OpenFile(clientfile, os.O_APPEND, os.ModePerm); err != nil {
			panic(err)
		}
		defer f.Close()
		clientProxyTmpl = appendTmpl

		unimplementedSvcMethods(&meta, clientfile)
	} else {
		if f, err = os.Create(clientfile); err != nil {
			panic(err)
		}
		defer f.Close()
		clientProxyTmpl = baseTmpl
	}

	servicePkg := astutils.GetPkgPath(dir)
	cfgPkg := astutils.GetPkgPath(filepath.Join(dir, "config"))
	dtoPkg := astutils.GetPkgPath(filepath.Join(dir, "dto"))

	funcMap := make(map[string]interface{})
	funcMap["isVarargs"] = v3helper.IsVarargs
	if tpl, err = template.New("clientproxy.go.tmpl").Funcs(funcMap).Parse(clientProxyTmpl); err != nil {
		panic(err)
	}
	if err = tpl.Execute(&buf, struct {
		VoPackage      string
		DtoPackage     string
		Meta           astutils.InterfaceMeta
		ServicePackage string
		ServiceAlias   string
		SvcName        string
		Version        string
	}{
		DtoPackage:     dtoPkg,
		Meta:           meta,
		ServicePackage: servicePkg,
		ServiceAlias:   ic.Package.Name,
		SvcName:        ic.Interfaces[0].Name,
		Version:        version.Release,
	}); err != nil {
		panic(err)
	}

	original, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	original = append(original, buf.Bytes()...)
	if tpl, err = template.New("cpimportimpl.go.tmpl").Parse(cpimportTmpl); err != nil {
		panic(err)
	}
	if err = tpl.Execute(&importBuf, struct {
		ConfigPackage string
		VoPackage     string
		DtoPackage    string
	}{
		DtoPackage:    dtoPkg,
		ConfigPackage: cfgPkg,
	}); err != nil {
		panic(err)
	}
	original = astutils.AppendImportStatements(original, importBuf.Bytes())
	astutils.FixImport(original, clientfile)
}

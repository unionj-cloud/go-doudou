package rest

import (
	"bytes"
	"encoding/json"
	"github.com/hako/durafmt"
	registry "github.com/unionj-cloud/go-doudou/v2/framework/registry/memberlist"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/memberlist"
	"net/http"
	"text/template"
	"time"
)

type Row struct {
	Index     int                    `json:"index"`
	SvcName   string                 `json:"svcName"`
	Hostname  string                 `json:"hostname"`
	BaseUrl   string                 `json:"baseUrl"`
	Status    string                 `json:"status"`
	Uptime    string                 `json:"uptime"`
	GoVer     string                 `json:"goVer"`
	GddVer    string                 `json:"gddVer"`
	BuildUser string                 `json:"buildUser"`
	BuildTime string                 `json:"buildTime"`
	Data      map[string]interface{} `json:"data"`
	Host      string                 `json:"host"`
	SvcPort   int                    `json:"svcPort"`
	MemPort   int                    `json:"memPort"`
}

func NewRow(index int, service registry.Service, uptime string, meta registry.NodeMeta, node *memberlist.Node) Row {
	status := "up"
	if node.State == memberlist.StateSuspect {
		status = "suspect"
	}
	return Row{
		Index:     index,
		SvcName:   service.Name,
		Hostname:  node.Name,
		BaseUrl:   service.BaseUrl(),
		Status:    status,
		Uptime:    uptime,
		GoVer:     meta.GoVer,
		GddVer:    meta.GddVer,
		BuildUser: meta.BuildUser,
		BuildTime: meta.BuildTime,
		Data:      service.Data,
		Host:      service.Host,
		SvcPort:   service.Port,
		MemPort:   int(node.Port),
	}
}

func MemberlistUIRoutes() []Route {
	return []Route{
		{
			Name:    "GetRegistry",
			Method:  "GET",
			Pattern: "/go-doudou/registry",
			HandlerFunc: func(writer http.ResponseWriter, request *http.Request) {
				var (
					tpl   *template.Template
					err   error
					buf   bytes.Buffer
					nodes []*memberlist.Node
					rows  []Row
					ret   []byte
				)
				if nodes, err = registry.AllNodes(); err != nil {
					http.Error(writer, err.Error(), http.StatusInternalServerError)
					return
				}
				var i int
				for _, node := range nodes {
					meta, _ := registry.ParseMeta(node)
					var uptime string
					if meta.RegisterAt != nil {
						uptime = time.Since(*meta.RegisterAt).String()
						if duration, err := durafmt.ParseString(uptime); err == nil {
							uptime = duration.LimitFirstN(2).String()
						}
					}
					for _, service := range meta.Services {
						i++
						rows = append(rows, NewRow(i, service, uptime, meta, node))
					}
				}
				ret, _ = json.Marshal(rows)
				if tpl, err = template.New("registry.tmpl").Parse(memberlistUITmpl); err != nil {
					panic(err)
				}
				if err = tpl.Execute(&buf, struct {
					Rows string
				}{
					Rows: string(ret),
				}); err != nil {
					panic(err)
				}
				writer.Header().Set("Content-Type", "text/html; charset=utf-8")
				writer.Write(buf.Bytes())
			},
		},
	}
}

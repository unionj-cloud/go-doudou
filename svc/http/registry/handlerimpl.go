package registry

import (
	"bytes"
	"encoding/json"
	"github.com/unionj-cloud/go-doudou/svc/registry"
	"net/http"
	"text/template"
)

type RegistryHandlerImpl struct {
}

type row struct {
	Index int `json:"index"`
	registry.NodeInfo
}

func (receiver *RegistryHandlerImpl) GetRegistry(_writer http.ResponseWriter, _req *http.Request) {
	var (
		tpl   *template.Template
		err   error
		buf   bytes.Buffer
		nodes []*registry.Node
		rows  []row
		ret   []byte
	)
	if registry.LocalNode != nil {
		nodes, _ = registry.LocalNode.Discover("")
	}
	for i, node := range nodes {
		rows = append(rows, row{
			Index:    i + 1,
			NodeInfo: node.Info(),
		})
	}
	ret, _ = json.Marshal(rows)
	if tpl, err = template.New("registry.tmpl").Parse(indexTmpl); err != nil {
		panic(err)
	}
	if err = tpl.Execute(&buf, struct {
		Rows string
	}{
		Rows: string(ret),
	}); err != nil {
		panic(err)
	}
	_writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	_writer.Write(buf.Bytes())
}

func NewRegistryHandler() RegistryHandler {
	return &RegistryHandlerImpl{}
}

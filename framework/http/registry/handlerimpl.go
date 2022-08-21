package registry

import (
	"bytes"
	"encoding/json"
	"github.com/unionj-cloud/go-doudou/framework/memberlist"
	"github.com/unionj-cloud/go-doudou/framework/registry"
	"net/http"
	"text/template"
)

// RegistryHandlerImpl define implementation for RegistryHandler
type RegistryHandlerImpl struct {
}

type row struct {
	Index int `json:"index"`
	registry.NodeInfo
}

// GetRegistry returns registry UI
func (receiver *RegistryHandlerImpl) GetRegistry(_writer http.ResponseWriter, _req *http.Request) {
	var (
		tpl   *template.Template
		err   error
		buf   bytes.Buffer
		nodes []*memberlist.Node
		rows  []row
		ret   []byte
	)
	if nodes, err = registry.AllNodes(); err != nil {
		http.Error(_writer, err.Error(), http.StatusInternalServerError)
		return
	}
	for i, node := range nodes {
		rows = append(rows, row{
			Index:    i + 1,
			NodeInfo: registry.Info(node),
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

// NewRegistryHandler creates new RegistryHandlerImpl
func NewRegistryHandler() RegistryHandler {
	return &RegistryHandlerImpl{}
}

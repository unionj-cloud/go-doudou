package v3

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/unionj-cloud/go-doudou/cmd/internal/astutils"
)

// service SearchService {
//   rpc Search(SearchRequest) returns (SearchResponse);
// }

const Syntax = "proto3"

type Service struct {
	Name      string
	Package   string
	GoPackage string
	Rpcs      []*Rpc
	Messages  []*Message
	Comments  []string
	Imports   []string
	Syntax    string
}

func NewService(name, goPackage string) *Service {
	return &Service{
		Name:      name,
		Package:   strcase.ToSnake(name),
		GoPackage: goPackage,
		Syntax:    Syntax,
	}
}

type Rpc struct {
	Name     string
	Request  *Message
	Response *Message
	Comments []string
}

func NewRpc(method astutils.MethodMeta) *Rpc {
	rpcName := strcase.ToCamel(method.Name)
	rpcRequest := NewRequest(fmt.Sprintf("%sRequest", rpcName), method.Params)
	rpcResponse := NewResponse(fmt.Sprintf("%sResponse", rpcName), method.Results)
	return &Rpc{
		Name:     rpcName,
		Request:  rpcRequest,
		Response: rpcResponse,
		Comments: method.Comments,
	}
}

func NewRequest(name string, params []astutils.FieldMeta) *Message {
	return &Message{
		Name:   name,
		Fields: nil,
	}
}

func NewResponse(name string, params []astutils.FieldMeta) *Message {
	return &Message{
		Name:   name,
		Fields: nil,
	}
}

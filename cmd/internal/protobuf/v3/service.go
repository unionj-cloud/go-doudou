package v3

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/unionj-cloud/go-doudou/cmd/internal/astutils"
	"github.com/unionj-cloud/go-doudou/toolkit/constants"
	"github.com/unionj-cloud/go-doudou/version"
	"time"
)

// service SearchService {
//   rpc Search(SearchRequest) returns (SearchResponse);
// }

const Syntax = "proto3"

type Service struct {
	Package   string
	GoPackage string
	Syntax    string
	// go-doudou version
	Version  string
	ProtoVer string
	Rpcs     []*Rpc
	Messages []*Message
	Enums    []*Enum
	Comments []string
	Imports  []string
}

func NewService(name, goPackage string) *Service {
	return &Service{
		Package:   strcase.ToSnake(name),
		GoPackage: goPackage,
		Syntax:    Syntax,
		Version:   version.Release,
		ProtoVer:  fmt.Sprintf("v%s", time.Now().Local().Format(constants.FORMAT10)),
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
	if len(params) == 0 {
		return Empty
	}
	if len(params) == 1 && params[0].Type == "context.Context" {
		return Empty
	}
	if len(params) > 0 {
		if params[0].Type == "context.Context" {
			params = params[1:]
		}
		if len(params) == 1 {
			if m, ok := MessageOf(params[0].Type).(*Message); ok && m.IsTopLevel {
				return m
			}
		}
	}
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

package v3

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/unionj-cloud/go-doudou/cmd/internal/astutils"
	"github.com/unionj-cloud/go-doudou/toolkit/constants"
	"github.com/unionj-cloud/go-doudou/version"
	"reflect"
	"strings"
	"time"
)

// service SearchService {
//   rpc Search(SearchRequest) returns (SearchResponse);
// }

const Syntax = "proto3"

type Service struct {
	Name      string
	Package   string
	GoPackage string
	Syntax    string
	// go-doudou version
	Version  string
	ProtoVer string
	Rpcs     []Rpc
	Messages []Message
	Enums    []Enum
	Comments []string
	Imports  []string
}

func NewService(name, goPackage string) Service {
	return Service{
		Name:      strcase.ToCamel(name) + "Service",
		Package:   strcase.ToSnake(name),
		GoPackage: goPackage,
		Syntax:    Syntax,
		Version:   version.Release,
		ProtoVer:  fmt.Sprintf("v%s", time.Now().Local().Format(constants.FORMAT10)),
	}
}

type StreamType int

const (
	biStream = iota + 1
	clientStream
	serverStream
)

type Rpc struct {
	Name       string
	Request    Message
	Response   Message
	Comments   []string
	StreamType StreamType
}

func NewRpc(method astutils.MethodMeta) Rpc {
	rpcName := strcase.ToCamel(method.Name) + "Rpc"
	rpcRequest := newRequest(rpcName, method.Params)
	if reflect.DeepEqual(rpcRequest, Empty) {
		ImportStore["google/protobuf/empty.proto"] = struct{}{}
	}
	if !strings.HasPrefix(rpcRequest.Name, "stream ") {
		if _, ok := MessageStore[rpcRequest.Name]; !ok {
			MessageStore[rpcRequest.Name] = rpcRequest
		}
	}
	rpcResponse := newResponse(rpcName, method.Results)
	if reflect.DeepEqual(rpcResponse, Empty) {
		ImportStore["google/protobuf/empty.proto"] = struct{}{}
	}
	if !strings.HasPrefix(rpcResponse.Name, "stream ") {
		if _, ok := MessageStore[rpcResponse.Name]; !ok {
			MessageStore[rpcResponse.Name] = rpcResponse
		}
	}
	var st StreamType
	if strings.HasPrefix(rpcRequest.Name, "stream ") && strings.HasPrefix(rpcResponse.Name, "stream ") {
		st = biStream
	} else if strings.HasPrefix(rpcRequest.Name, "stream ") {
		st = clientStream
	} else if strings.HasPrefix(rpcResponse.Name, "stream ") {
		st = serverStream
	}
	return Rpc{
		Name:       rpcName,
		Request:    rpcRequest,
		Response:   rpcResponse,
		Comments:   method.Comments,
		StreamType: st,
	}
}

func newRequest(rpcName string, params []astutils.FieldMeta) Message {
	if len(params) == 0 {
		return Empty
	}
	if len(params) == 1 && params[0].Type == "context.Context" {
		return Empty
	}
	if params[0].Type == "context.Context" {
		params = params[1:]
	}
	if len(params) == 1 {
		if m, ok := MessageOf(params[0].Type).(Message); ok && m.IsTopLevel {
			if strings.HasPrefix(params[0].Name, "stream") {
				m.Name = "stream " + m.Name
			}
			return m
		}
	}
	var fields []Field
	for i, field := range params {
		fields = append(fields, newField(field, i+1))
	}
	return Message{
		Name:       strcase.ToCamel(rpcName + "Request"),
		Fields:     fields,
		IsTopLevel: true,
	}
}

func newResponse(rpcName string, params []astutils.FieldMeta) Message {
	if len(params) == 0 {
		return Empty
	}
	if len(params) == 1 && params[0].Type == "error" {
		return Empty
	}
	if params[len(params)-1].Type == "error" {
		params = params[:len(params)-1]
	}
	if len(params) == 1 {
		if m, ok := MessageOf(params[0].Type).(Message); ok && m.IsTopLevel {
			if strings.HasPrefix(params[0].Name, "stream") {
				m.Name = "stream " + m.Name
			}
			return m
		}
	}
	var fields []Field
	for i, field := range params {
		fields = append(fields, newField(field, i+1))
	}
	return Message{
		Name:       strcase.ToCamel(rpcName + "Response"),
		Fields:     fields,
		IsTopLevel: true,
	}
}

package template

const Svc = EditMarkForGDD + `
package service

import (
	"context"
	"{{.DtoPackage}}"
	v3 "github.com/unionj-cloud/go-doudou/v2/toolkit/openapi/v3"
)

//go:generate go-doudou svc http
//go:generate go-doudou svc grpc

type {{.InterfaceName}} interface {
}
`

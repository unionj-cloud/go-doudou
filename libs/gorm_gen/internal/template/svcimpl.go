package template

const SvcImpl = EditMarkForGDD + `
package service

import ()

var _ {{.InterfaceName.Name}} = (*{{.InterfaceName.Name}}Impl)(nil)

type {{.InterfaceName.Name}}Impl struct {
	conf *config.Config
}

func New{{.InterfaceName.Name}}(conf *config.Config) *{{.InterfaceName.Name}}Impl {
	return &{{.InterfaceName.Name}}Impl{
		conf: conf,
	}
}
`

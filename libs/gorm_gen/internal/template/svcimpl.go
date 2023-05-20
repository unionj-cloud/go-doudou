package template

const SvcImpl = EditMarkForGDD + `
package service

import ()

func init() {
	query.SetDefault(database.Db)
}

var _ {{.InterfaceName}} = (*{{.InterfaceName}}Impl)(nil)

type {{.InterfaceName}}Impl struct {
	conf *config.Config
	pg   *paginate.Pagination
}

func New{{.InterfaceName}}(conf *config.Config) *{{.InterfaceName}}Impl {
	pg := paginate.New(&paginate.Config{
		FieldSelectorEnabled: true,
	})
	return &{{.InterfaceName}}Impl{
		conf: conf,
		pg:   pg,
	}
}
`

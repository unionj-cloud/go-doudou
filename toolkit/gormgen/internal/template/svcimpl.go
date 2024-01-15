package template

const SvcImpl = EditMarkForGDD + `
package service

import ()

var _ {{.InterfaceName}} = (*{{.InterfaceName}}Impl)(nil)

type {{.InterfaceName}}Impl struct {
	conf *config.Config
	pg   *paginate.Pagination
	q    *query.Query
}

func New{{.InterfaceName}}(conf *config.Config) *{{.InterfaceName}}Impl {
	pg := paginate.New(&paginate.Config{
		FieldSelectorEnabled: true,
	})
	return &{{.InterfaceName}}Impl{
		conf: conf,
		pg:   pg,
		q: query.Q,
	}
}

func (receiver {{.InterfaceName}}Impl) clone(q *query.Query) *{{.InterfaceName}}Impl {
	receiver.q = q
	return &receiver
}
`

package template

const AppendSvcImpl = `
// Post{{.ModelStructName}} {{.StructComment}}
` + NotEditMarkForGDDShort + `
func (receiver *{{.InterfaceName}}Impl) Post{{.ModelStructName}}(ctx context.Context, body dto.{{.ModelStructName}}) (data {{.PriKeyType}}, err error) {
	m := &model.{{.ModelStructName}}{}
	copier.Copy(m, body)
	u := query.{{.ModelStructName}}
	err = errorx.Wrap(u.WithContext(ctx).Create(m))
	data = m.ID
	return
}

// Get{{.ModelStructName}}_Id {{.StructComment}}
` + NotEditMarkForGDDShort + `
func (receiver *{{.InterfaceName}}Impl) Get{{.ModelStructName}}_Id(ctx context.Context, id {{.PriKeyType}}) (data dto.{{.ModelStructName}}, err error) {
	u := query.{{.ModelStructName}}
	m, err := u.WithContext(ctx).Where(u.ID.Eq(id)).First()
	err = errorx.Wrap(err)
	copier.Copy(&data, m)
	return
}

// Put{{.ModelStructName}} {{.StructComment}}
` + NotEditMarkForGDDShort + `
func (receiver *{{.InterfaceName}}Impl) Put{{.ModelStructName}}(ctx context.Context, body dto.{{.ModelStructName}}) (err error) {
	m := &model.{{.ModelStructName}}{}
	copier.Copy(m, body)
	u := query.{{.ModelStructName}}
	_, err = u.WithContext(ctx).Where(u.ID.Eq(body.ID)).Updates(m)
	err = errorx.Wrap(err)
	return
}

// Delete{{.ModelStructName}}_Id {{.StructComment}}
` + NotEditMarkForGDDShort + `
func (receiver *{{.InterfaceName}}Impl) Delete{{.ModelStructName}}_Id(ctx context.Context, id {{.PriKeyType}}) (err error) {
	u := query.{{.ModelStructName}}
	_, err = u.WithContext(ctx).Where(u.ID.Eq(id)).Delete()
	err = errorx.Wrap(err)
	return
}

// Get{{.ModelStructName}}s {{.StructComment}}
` + NotEditMarkForGDDShort + `
func (receiver *{{.InterfaceName}}Impl) Get{{.ModelStructName}}s(ctx context.Context, parameter dto.Parameter) (data dto.Page, err error) {
	paginated := receiver.pg.With(database.Db.Model(&model.{{.ModelStructName}}{})).Request(paginate.Parameter(parameter)).Response(&[]model.{{.ModelStructName}}{})
	data = dto.Page(paginated)
	return
}

`

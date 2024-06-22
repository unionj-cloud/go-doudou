package template

const AppendSvcImpl = `
// PostGen{{.ModelStructName}} {{.StructComment}}
` + NotEditMarkForGDDShort + `
func (receiver *{{.InterfaceName}}Impl) PostGen{{.ModelStructName}}(ctx context.Context, body dto.{{.ModelStructName}}) (data {{.PriKeyType}}, err error) {
	var m model.{{.ModelStructName}}
	copier.DeepCopy(body, &m)
	u := receiver.q.{{.ModelStructName}}
	err = errors.WithStack(u.WithContext(ctx).Create(&m))
	data = m.ID
	return
}

// PostGen{{.ModelStructName}}s {{.StructComment}}
` + NotEditMarkForGDDShort + `
func (receiver *{{.InterfaceName}}Impl) PostGen{{.ModelStructName}}s(ctx context.Context, body []dto.{{.ModelStructName}}) (data []{{.PriKeyType}}, err error) {
	list := make([]*model.{{.ModelStructName}}, 0, len(body))
	for _, item := range body {
		var m model.{{.ModelStructName}}
		copier.DeepCopy(item, &m)
		list = append(list, &m)
	}
	u := receiver.q.{{.ModelStructName}}
	if err = errors.WithStack(u.WithContext(ctx).Create(list...)); err != nil {
		return
	}
	data = make([]{{.PriKeyType}}, 0, len(list))
	for _, item := range list {
		data = append(data, item.ID)
	}
	return
}

// GetGen{{.ModelStructName}}_Id {{.StructComment}}
` + NotEditMarkForGDDShort + `
func (receiver *{{.InterfaceName}}Impl) GetGen{{.ModelStructName}}_Id(ctx context.Context, id {{.PriKeyType}}) (data dto.{{.ModelStructName}}, err error) {
	u := receiver.q.{{.ModelStructName}}
	m, err := u.WithContext(ctx).Where(u.ID.Eq(id)).First()
	if err != nil {
		return dto.{{.ModelStructName}}{}, errors.WithStack(err)
	}
	copier.DeepCopy(m, &data)
	return
}

// PutGen{{.ModelStructName}} {{.StructComment}}
` + NotEditMarkForGDDShort + `
func (receiver *{{.InterfaceName}}Impl) PutGen{{.ModelStructName}}(ctx context.Context, body dto.{{.ModelStructName}}) (err error) {
	var m model.{{.ModelStructName}}
	copier.DeepCopy(body, &m)
	u := receiver.q.{{.ModelStructName}}
	_, err = u.WithContext(ctx).Where(u.ID.Eq(body.ID)).Updates(m)
	return errors.WithStack(err)
}

// DeleteGen{{.ModelStructName}}_Id {{.StructComment}}
` + NotEditMarkForGDDShort + `
func (receiver *{{.InterfaceName}}Impl) DeleteGen{{.ModelStructName}}_Id(ctx context.Context, id {{.PriKeyType}}) (err error) {
	u := receiver.q.{{.ModelStructName}}
	_, err = u.WithContext(ctx).Where(u.ID.Eq(id)).Delete()
	return errors.WithStack(err)
}

// GetGen{{.ModelStructName}}s {{.StructComment}}
` + NotEditMarkForGDDShort + `
func (receiver *{{.InterfaceName}}Impl) GetGen{{.ModelStructName}}s(ctx context.Context, parameter dto.Parameter) (data dto.Page, err error) {
	paginated := receiver.pg.With(database.Db.Model(&model.{{.ModelStructName}}{})).Request(parameter).Response(&[]model.{{.ModelStructName}}{})
	data = dto.Page(paginated)
	return
}

`

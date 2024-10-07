package template

const AppendSvcImpl = `
// PostGen{{.ModelStructName}} {{.StructComment}}
` + NotEditMarkForGDDShort + `
func (receiver *{{.InterfaceName}}Impl) PostGen{{.ModelStructName}}Rpc(ctx context.Context, body *pb.{{.ModelStructName}}) (data *pb.{{.ModelStructName}}, err error) {
	var m model.{{.ModelStructName}}
	copier.DeepCopy(body, &m)
	u := receiver.q.{{.ModelStructName}}
	if err = u.WithContext(ctx).Create(&m); err != nil {
		return nil, errors.WithStack(err)
	}
	data = new(pb.{{.ModelStructName}})
	copier.DeepCopy(m, data)
	return
}

// GetGen{{.ModelStructName}} {{.StructComment}}
` + NotEditMarkForGDDShort + `
func (receiver *{{.InterfaceName}}Impl) GetGen{{.ModelStructName}}Rpc(ctx context.Context, body *pb.{{.ModelStructName}}) (data *pb.{{.ModelStructName}}, err error) {
	u := receiver.q.{{.ModelStructName}}
	m, err := u.WithContext(ctx).Where(u.ID.Eq(body.ID)).First()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	data = new(pb.{{.ModelStructName}})
	copier.DeepCopy(m, data)
	return
}

// PutGen{{.ModelStructName}} {{.StructComment}}
` + NotEditMarkForGDDShort + `
func (receiver *{{.InterfaceName}}Impl) PutGen{{.ModelStructName}}Rpc(ctx context.Context, body *pb.{{.ModelStructName}}) (*emptypb.Empty, error) {
	var m model.{{.ModelStructName}}
	copier.DeepCopy(body, &m)
	u := receiver.q.{{.ModelStructName}}
	_, err := u.WithContext(ctx).Where(u.ID.Eq(body.ID)).Updates(m)
	return &emptypb.Empty{}, errors.WithStack(err)
}

// DeleteGen{{.ModelStructName}} {{.StructComment}}
` + NotEditMarkForGDDShort + `
func (receiver *{{.InterfaceName}}Impl) DeleteGen{{.ModelStructName}}Rpc(ctx context.Context, body *pb.{{.ModelStructName}}) (*emptypb.Empty, error) {
	u := receiver.q.{{.ModelStructName}}
	_, err := u.WithContext(ctx).Where(u.ID.Eq(body.ID)).Delete()
	return &emptypb.Empty{}, errors.WithStack(err)
}

// GetGen{{.ModelStructName}}s {{.StructComment}}
` + NotEditMarkForGDDShort + `
func (receiver *{{.InterfaceName}}Impl) GetGen{{.ModelStructName}}sRpc(ctx context.Context, request *pb.Parameter) (data *pb.Page, err error) {
	var body dto.Parameter
	copier.DeepCopy(request, &body)
	resCxt := receiver.pg.With(receiver.q.Db().Model(&model.{{.ModelStructName}}{})).Request(body)
	paginated := resCxt.Response(&[]model.{{.ModelStructName}}{})
	if resCxt.Error() != nil {
		return nil, errors.WithStack(resCxt.Error())
	}
	data = new(pb.Page)
	copier.DeepCopy(paginated, data)
	return
}

`

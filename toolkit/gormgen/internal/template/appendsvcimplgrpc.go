package template

const AppendSvcImplGrpc = `
// PostGen{{.ModelStructName}}Rpc {{.StructComment}}
` + NotEditMarkForGDDShort + `
func (receiver *{{.InterfaceName}}Impl) PostGen{{.ModelStructName}}Rpc(ctx context.Context, request *pb.{{.ModelStructName}}) (*pb.PostGen{{.ModelStructName}}RpcResponse, error) {
	var body dto.{{.ModelStructName}}
	copier.DeepCopy(request, &body)
	data, err := receiver.PostGen{{.ModelStructName}}(ctx, body)
	return &pb.PostGen{{.ModelStructName}}RpcResponse{
		Data: data,
	}, errors.WithStack(err)
}

// PostGen{{.ModelStructName}}sRpc {{.StructComment}}
` + NotEditMarkForGDDShort + `
func (receiver *{{.InterfaceName}}Impl) PostGen{{.ModelStructName}}sRpc(ctx context.Context, request *pb.PostGen{{.ModelStructName}}sRpcRequest) (*pb.PostGen{{.ModelStructName}}sRpcResponse, error) {
	list := make([]dto.{{.ModelStructName}}, 0, len(request.Body))
	for _, item := range request.Body {
		var d dto.{{.ModelStructName}}
		copier.DeepCopy(item, &d)
		list = append(list, d)
	}
	data, err := receiver.PostGen{{.ModelStructName}}s(ctx, list)
	return &pb.PostGen{{.ModelStructName}}sRpcResponse{
		Data: data,
	}, errors.WithStack(err)
}

// GetGen{{.ModelStructName}}IdRpc {{.StructComment}}
` + NotEditMarkForGDDShort + `
func (receiver *{{.InterfaceName}}Impl) GetGen{{.ModelStructName}}IdRpc(ctx context.Context, request *pb.GetGen{{.ModelStructName}}IdRpcRequest) (*pb.{{.ModelStructName}}, error) {
	data, err := receiver.GetGen{{.ModelStructName}}_Id(ctx, request.Id)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var ret pb.{{.ModelStructName}}
	copier.DeepCopy(data, &ret)
	return &ret, nil
}

// PutGen{{.ModelStructName}}Rpc {{.StructComment}}
` + NotEditMarkForGDDShort + `
func (receiver *{{.InterfaceName}}Impl) PutGen{{.ModelStructName}}Rpc(ctx context.Context, request *pb.{{.ModelStructName}}) (*emptypb.Empty, error) {
	var body dto.{{.ModelStructName}}
	copier.DeepCopy(request, &body)
	return &emptypb.Empty{}, errors.WithStack(receiver.PutGen{{.ModelStructName}}(ctx, body))
}

// DeleteGen{{.ModelStructName}}IdRpc {{.StructComment}}
` + NotEditMarkForGDDShort + `
func (receiver *{{.InterfaceName}}Impl) DeleteGen{{.ModelStructName}}IdRpc(ctx context.Context, request *pb.DeleteGen{{.ModelStructName}}IdRpcRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, errors.WithStack(receiver.DeleteGen{{.ModelStructName}}_Id(ctx, request.Id))
}

// GetGen{{.ModelStructName}}sRpc {{.StructComment}}
` + NotEditMarkForGDDShort + `
func (receiver *{{.InterfaceName}}Impl) GetGen{{.ModelStructName}}sRpc(ctx context.Context, request *pb.Parameter) (*pb.Page, error) {
	filters := make([]interface{}, 0, len(request.Filters))
	for _, item := range request.Filters {
		str := wrappers.StringValue{}
		if err := anypb.UnmarshalTo(item, &str, proto.UnmarshalOptions{}); err != nil {
			return nil, errors.WithStack(err)
		}
		filters = append(filters, str.Value)
	}
	var parameter dto.Parameter
	copier.DeepCopy(request, &parameter)
	parameter.Filters = filters
	data, err := receiver.GetGen{{.ModelStructName}}s(ctx, parameter)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	items := make([]*anypb.Any, 0, len(data.Items))
	for _, item := range data.Items {
		d := dto.{{.ModelStructName}}(item.(model.{{.ModelStructName}}))
		var msg pb.{{.ModelStructName}}
		copier.DeepCopy(d, &msg)
		a, err := anypb.New(&msg)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		items = append(items, a)
	}
	var ret pb.Page
	copier.DeepCopy(data, &ret)
	ret.Items = items
	return &ret, nil
}

`

package template

const AppendSvcImplGrpc = `
// Post{{.ModelStructName}}Rpc {{.StructComment}}
` + NotEditMarkForGDDShort + `
func (receiver *{{.InterfaceName}}Impl) Post{{.ModelStructName}}Rpc(ctx context.Context, request *pb.{{.ModelStructName}}) (*pb.Post{{.ModelStructName}}RpcResponse, error) {
	var body dto.{{.ModelStructName}}
	jsoncopier.DeepCopy(request, &body)
	data, err := receiver.Post{{.ModelStructName}}(ctx, body)
	return &pb.Post{{.ModelStructName}}RpcResponse{
		Data: data,
	}, errors.WithStack(err)
}

// Post{{.ModelStructName}}sRpc {{.StructComment}}
` + NotEditMarkForGDDShort + `
func (receiver *{{.InterfaceName}}Impl) Post{{.ModelStructName}}sRpc(ctx context.Context, request *pb.Post{{.ModelStructName}}sRpcRequest) (*pb.Post{{.ModelStructName}}sRpcResponse, error) {
	list := make([]dto.{{.ModelStructName}}, 0, len(request.Body))
	for _, item := range request.Body {
		var d dto.{{.ModelStructName}}
		jsoncopier.DeepCopy(item, &d)
		list = append(list, d)
	}
	data, err := receiver.Post{{.ModelStructName}}s(ctx, list)
	return &pb.Post{{.ModelStructName}}sRpcResponse{
		Data: data,
	}, errors.WithStack(err)
}

// Get{{.ModelStructName}}IdRpc {{.StructComment}}
` + NotEditMarkForGDDShort + `
func (receiver *{{.InterfaceName}}Impl) Get{{.ModelStructName}}IdRpc(ctx context.Context, request *pb.Get{{.ModelStructName}}IdRpcRequest) (*pb.{{.ModelStructName}}, error) {
	data, err := receiver.Get{{.ModelStructName}}_Id(ctx, request.Id)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var ret pb.{{.ModelStructName}}
	jsoncopier.DeepCopy(data, &ret)
	return &ret, nil
}

// Put{{.ModelStructName}}Rpc {{.StructComment}}
` + NotEditMarkForGDDShort + `
func (receiver *{{.InterfaceName}}Impl) Put{{.ModelStructName}}Rpc(ctx context.Context, request *pb.{{.ModelStructName}}) (*emptypb.Empty, error) {
	var body dto.{{.ModelStructName}}
	jsoncopier.DeepCopy(request, &body)
	return &emptypb.Empty{}, errors.WithStack(receiver.Put{{.ModelStructName}}(ctx, body))
}

// Delete{{.ModelStructName}}IdRpc {{.StructComment}}
` + NotEditMarkForGDDShort + `
func (receiver *{{.InterfaceName}}Impl) Delete{{.ModelStructName}}IdRpc(ctx context.Context, request *pb.Delete{{.ModelStructName}}IdRpcRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, errors.WithStack(receiver.Delete{{.ModelStructName}}_Id(ctx, request.Id))
}

// Get{{.ModelStructName}}sRpc {{.StructComment}}
` + NotEditMarkForGDDShort + `
func (receiver *{{.InterfaceName}}Impl) Get{{.ModelStructName}}sRpc(ctx context.Context, request *pb.Parameter) (*pb.Page, error) {
	filters := make([]interface{}, 0, len(request.Filters))
	for _, item := range request.Filters {
		str := wrappers.StringValue{}
		if err := anypb.UnmarshalTo(item, &str, proto.UnmarshalOptions{}); err != nil {
			return nil, errors.WithStack(err)
		}
		filters = append(filters, str.Value)
	}
	var parameter dto.Parameter
	jsoncopier.DeepCopy(request, &parameter)
	parameter.Filters = filters
	data, err := receiver.Get{{.ModelStructName}}s(ctx, parameter)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	items := make([]*anypb.Any, 0, len(data.Items))
	for _, item := range data.Items {
		d := dto.{{.ModelStructName}}(item.(model.{{.ModelStructName}}))
		var msg pb.{{.ModelStructName}}
		jsoncopier.DeepCopy(d, &msg)
		a, err := anypb.New(&msg)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		items = append(items, a)
	}
	var ret pb.Page
	jsoncopier.DeepCopy(data, &ret)
	ret.Items = items
	return &ret, nil
}

`

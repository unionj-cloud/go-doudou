package template

const AppendSvc = `
// PostGen{{.ModelStructName}} {{.StructComment}}
` + NotEditMarkForGDDShort + `
PostGen{{.ModelStructName}}(ctx context.Context, body model.{{.ModelStructName}}) (model.{{.ModelStructName}}, error)

// GetGen{{.ModelStructName}} {{.StructComment}}
` + NotEditMarkForGDDShort + `
GetGen{{.ModelStructName}}(ctx context.Context, body model.{{.ModelStructName}}) (model.{{.ModelStructName}}, error)

// PutGen{{.ModelStructName}} {{.StructComment}}
` + NotEditMarkForGDDShort + `
PutGen{{.ModelStructName}}(ctx context.Context, body model.{{.ModelStructName}}) error

// DeleteGen{{.ModelStructName}} {{.StructComment}}
` + NotEditMarkForGDDShort + `
DeleteGen{{.ModelStructName}}(ctx context.Context, body model.{{.ModelStructName}}) error

// GetGen{{.ModelStructName}}s {{.StructComment}}
` + NotEditMarkForGDDShort + `
GetGen{{.ModelStructName}}s(ctx context.Context, parameter dto.Parameter) (data dto.Page, err error)

`

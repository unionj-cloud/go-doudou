package template

const AppendSvc = `
// PostGen{{.ModelStructName}} {{.StructComment}}
` + NotEditMarkForGDDShort + `
PostGen{{.ModelStructName}}(ctx context.Context, body dto.{{.ModelStructName}}) (data {{.PriKeyType}}, err error)

// PostGen{{.ModelStructName}}s {{.StructComment}}
` + NotEditMarkForGDDShort + `
PostGen{{.ModelStructName}}s(ctx context.Context, body []dto.{{.ModelStructName}}) (data []{{.PriKeyType}}, err error)

// GetGen{{.ModelStructName}}_Id {{.StructComment}}
` + NotEditMarkForGDDShort + `
GetGen{{.ModelStructName}}_Id(ctx context.Context, id {{.PriKeyType}}) (data dto.{{.ModelStructName}}, err error)

// PutGen{{.ModelStructName}} {{.StructComment}}
` + NotEditMarkForGDDShort + `
PutGen{{.ModelStructName}}(ctx context.Context, body dto.{{.ModelStructName}}) error

// DeleteGen{{.ModelStructName}}_Id {{.StructComment}}
` + NotEditMarkForGDDShort + `
DeleteGen{{.ModelStructName}}_Id(ctx context.Context, id {{.PriKeyType}}) error

// GetGen{{.ModelStructName}}s {{.StructComment}}
` + NotEditMarkForGDDShort + `
GetGen{{.ModelStructName}}s(ctx context.Context, parameter dto.Parameter) (data dto.Page, err error)

`

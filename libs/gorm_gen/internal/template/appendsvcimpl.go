package template

const AppendSvcImpl = `
// Post{{.ModelStructName}} {{.StructComment}}
` + NotEditMarkForGDDShort + `
Post{{.ModelStructName}}(ctx context.Context, body dto.{{.ModelStructName}}) (data {{.PriKeyType}}, err error)

// Get{{.ModelStructName}}_Id {{.StructComment}}
` + NotEditMarkForGDDShort + `
Get{{.ModelStructName}}_Id(ctx context.Context, id {{.PriKeyType}}) (data dto.{{.ModelStructName}}, err error)

// Put{{.ModelStructName}} {{.StructComment}}
` + NotEditMarkForGDDShort + `
Put{{.ModelStructName}}(ctx context.Context, body dto.{{.ModelStructName}}) error

// Delete{{.ModelStructName}}_Id {{.StructComment}}
` + NotEditMarkForGDDShort + `
Delete{{.ModelStructName}}_Id(ctx context.Context, id {{.PriKeyType}}) error

// Get{{.ModelStructName}}s {{.StructComment}}
` + NotEditMarkForGDDShort + `
Get{{.ModelStructName}}s(ctx context.Context, parameter dto.Parameter) (data dto.Page, err error)

`

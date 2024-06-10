package template

// Dto used as a variable because it cannot load template file after packed, params still can pass file
const Dto = EditMark + `
package dto

import (
	"encoding/json"
	"time"

	"github.com/wubin1989/datatypes"
	"github.com/wubin1989/gorm"
	"github.com/wubin1989/gorm/schema"
	{{range .ImportPkgPaths}}{{.}} ` + "\n" + `{{end}}
)

// {{.ModelStructName}} {{.StructComment}}
type {{.ModelStructName}} struct {
    {{range .Fields}}
    {{if .MultilineComment -}}
	/*
{{.ColumnComment}}
    */
	{{end -}}
    {{.Name}} {{.Type | convert}} ` + "`{{.Tags}}` " +
	"{{if not .MultilineComment}}{{if .ColumnComment}}// {{.ColumnComment}}{{end}}{{end}}" +
	`{{end}}
}

`

package codegen

import (
	"bytes"
	"github.com/iancoleman/strcase"
	log "github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/astutils"
	"github.com/unionj-cloud/go-doudou/v2/cmd/internal/ddl/table"
	"github.com/unionj-cloud/go-doudou/v2/version"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var daosqltmpl = `{{` + "`" + `{{` + "`" + `}}define "NoneZeroSet"{{` + "`" + `}}` + "`" + `}}
	{{- range $i, $co := .UpdateColumns}}
	{{` + "`" + `{{` + "`" + `}}- if .{{$co.Meta.Name}}{{` + "`" + `}}` + "`" + `}}
	` + "`" + `{{$co.Name}}` + "`" + `=:{{$co.Name}},
	{{` + "`" + `{{` + "`" + `}}- end{{` + "`" + `}}` + "`" + `}}
	{{- end}}
{{` + "`" + `{{` + "`" + `}}end{{` + "`" + `}}` + "`" + `}}

{{` + "`" + `{{` + "`" + `}}define "Insert{{.EntityName}}"{{` + "`" + `}}` + "`" + `}}
INSERT INTO ` + "`" + `{{.Schema}}` + "`" + `.` + "`" + `{{.TableName}}` + "`" + `
({{- range $i, $co := .InsertColumns}}
{{- if $i}},{{end}}
` + "`" + `{{$co.Name}}` + "`" + `
{{- end }})
VALUES ({{- range $i, $co := .InsertColumns}}
	   {{- if $i}},{{end}}
	   :{{$co.Name}}
	   {{- end }})
{{` + "`" + `{{` + "`" + `}}end{{` + "`" + `}}` + "`" + `}}

{{` + "`" + `{{` + "`" + `}}define "Update{{.EntityName}}"{{` + "`" + `}}` + "`" + `}}
UPDATE ` + "`" + `{{.Schema}}` + "`" + `.` + "`" + `{{.TableName}}` + "`" + `
SET
	{{- range $i, $co := .UpdateColumns}}
	{{- if $i}},{{end}}
	` + "`" + `{{$co.Name}}` + "`" + `=:{{$co.Name}}
	{{- end }}
WHERE
    ` + "`" + `{{.Pk.Name}}` + "`" + ` =:{{.Pk.Name}}
{{` + "`" + `{{` + "`" + `}}end{{` + "`" + `}}` + "`" + `}}

{{` + "`" + `{{` + "`" + `}}define "Update{{.EntityName}}NoneZero"{{` + "`" + `}}` + "`" + `}}
UPDATE ` + "`" + `{{.Schema}}` + "`" + `.` + "`" + `{{.TableName}}` + "`" + `
SET
    {{` + "`" + `{{` + "`" + `}}Eval "NoneZeroSet" . | TrimSuffix ","{{` + "`" + `}}` + "`" + `}}
WHERE
    ` + "`" + `{{.Pk.Name}}` + "`" + `=:{{.Pk.Name}}
{{` + "`" + `{{` + "`" + `}}end{{` + "`" + `}}` + "`" + `}}

{{` + "`" + `{{` + "`" + `}}define "Upsert{{.EntityName}}"{{` + "`" + `}}` + "`" + `}}
INSERT INTO ` + "`" + `{{.Schema}}` + "`" + `.` + "`" + `{{.TableName}}` + "`" + `
({{- range $i, $co := .InsertColumns}}
{{- if $i}},{{end}}
` + "`" + `{{$co.Name}}` + "`" + `
{{- end }})
VALUES ({{- range $i, $co := .InsertColumns}}
        {{- if $i}},{{end}}
        :{{$co.Name}}
        {{- end }}) ON DUPLICATE KEY
UPDATE
		{{- range $i, $co := .UpdateColumns}}
		{{- if $i}},{{end}}
		` + "`" + `{{$co.Name}}` + "`" + `=:{{$co.Name}}
		{{- end }}
{{` + "`" + `{{` + "`" + `}}end{{` + "`" + `}}` + "`" + `}}

{{` + "`" + `{{` + "`" + `}}define "Upsert{{.EntityName}}NoneZero"{{` + "`" + `}}` + "`" + `}}
INSERT INTO ` + "`" + `{{.Schema}}` + "`" + `.` + "`" + `{{.TableName}}` + "`" + `
({{- range $i, $co := .InsertColumns}}
{{- if $i}},{{end}}
` + "`" + `{{$co.Name}}` + "`" + `
{{- end }})
VALUES ({{- range $i, $co := .InsertColumns}}
        {{- if $i}},{{end}}
        :{{$co.Name}}
        {{- end }}) ON DUPLICATE KEY
UPDATE
		{{` + "`" + `{{` + "`" + `}}Eval "NoneZeroSet" . | TrimSuffix ","{{` + "`" + `}}` + "`" + `}}
{{` + "`" + `{{` + "`" + `}}end{{` + "`" + `}}` + "`" + `}}

{{` + "`" + `{{` + "`" + `}}define "Get{{.EntityName}}"{{` + "`" + `}}` + "`" + `}}
select *
from ` + "`" + `{{.Schema}}` + "`" + `.` + "`" + `{{.TableName}}` + "`" + `
where ` + "`" + `{{.Pk.Name}}` + "`" + ` = ?
{{` + "`" + `{{` + "`" + `}}end{{` + "`" + `}}` + "`" + `}}

{{` + "`" + `{{` + "`" + `}}define "Update{{.EntityName}}s"{{` + "`" + `}}` + "`" + `}}
UPDATE ` + "`" + `{{.Schema}}` + "`" + `.` + "`" + `{{.TableName}}` + "`" + `
SET
	{{- range $i, $co := .UpdateColumns}}
	{{- if $i}},{{end}}
	` + "`" + `{{$co.Name}}` + "`" + `=:{{$co.Name}}
	{{- end }}
{{` + "`" + `{{` + "`" + `}}end{{` + "`" + `}}` + "`" + `}}

{{` + "`" + `{{` + "`" + `}}define "Update{{.EntityName}}sNoneZero"{{` + "`" + `}}` + "`" + `}}
UPDATE ` + "`" + `{{.Schema}}` + "`" + `.` + "`" + `{{.TableName}}` + "`" + `
SET
    {{` + "`" + `{{` + "`" + `}}Eval "NoneZeroSet" . | TrimSuffix ","{{` + "`" + `}}` + "`" + `}}
{{` + "`" + `{{` + "`" + `}}end{{` + "`" + `}}` + "`" + `}}

{{` + "`" + `{{` + "`" + `}}define "InsertIgnore{{.EntityName}}"{{` + "`" + `}}` + "`" + `}}
INSERT IGNORE INTO ` + "`" + `{{.Schema}}` + "`" + `.` + "`" + `{{.TableName}}` + "`" + `
({{- range $i, $co := .InsertColumns}}
{{- if $i}},{{end}}
` + "`" + `{{$co.Name}}` + "`" + `
{{- end }})
VALUES ({{- range $i, $co := .InsertColumns}}
	   {{- if $i}},{{end}}
	   :{{$co.Name}}
	   {{- end }})
{{` + "`" + `{{` + "`" + `}}end{{` + "`" + `}}` + "`" + `}}

{{` + "`" + `{{` + "`" + `}}define "UpdateClause{{.EntityName}}"{{` + "`" + `}}` + "`" + `}}
ON DUPLICATE KEY
UPDATE
		{{- range $i, $co := .UpdateColumns}}
		{{- if $i}},{{end}}
		` + "`" + `{{$co.Name}}` + "`" + `=VALUES({{$co.Name}})
		{{- end }}
{{` + "`" + `{{` + "`" + `}}end{{` + "`" + `}}` + "`" + `}}

{{` + "`" + `{{` + "`" + `}}define "UpdateClauseSelect{{.EntityName}}"{{` + "`" + `}}` + "`" + `}}
ON DUPLICATE KEY
UPDATE
		{{` + "`" + `{{` + "`" + `}}- range $i, $co := .Columns{{` + "`" + `}}` + "`" + `}}
		{{` + "`" + `{{` + "`" + `}}- if $i{{` + "`" + `}}` + "`" + `}},{{` + "`" + `{{` + "`" + `}}end{{` + "`" + `}}` + "`" + `}}
		` + "`" + `{{` + "`" + `{{` + "`" + `}}$co{{` + "`" + `}}` + "`" + `}}` + "`" + `=VALUES({{` + "`" + `{{` + "`" + `}}$co{{` + "`" + `}}` + "`" + `}})
		{{` + "`" + `{{` + "`" + `}}- end {{` + "`" + `}}` + "`" + `}}
{{` + "`" + `{{` + "`" + `}}end{{` + "`" + `}}` + "`" + `}}`

// GenDaoSQL generates sql statements used by dao layer
func GenDaoSQL(entityPath string, t table.Table, folder ...string) error {
	var (
		err      error
		daopath  string
		f        *os.File
		funcMap  map[string]interface{}
		tpl      *template.Template
		iColumns []table.Column
		uColumns []table.Column
		df       string
		sqlBuf   bytes.Buffer
	)
	df = "dao"
	if len(folder) > 0 {
		df = folder[0]
	}
	daopath = filepath.Join(filepath.Dir(entityPath), df)
	_ = os.MkdirAll(daopath, os.ModePerm)
	daofile := filepath.Join(daopath, strings.ToLower(t.Meta.Name)+"daosql.go")
	if _, err = os.Stat(daofile); os.IsNotExist(err) {
		f, _ = os.Create(daofile)
		defer f.Close()

		funcMap = make(map[string]interface{})
		funcMap["ToSnake"] = strcase.ToSnake
		tpl, _ = template.New("daosql.tmpl").Funcs(funcMap).Parse(daosqltmpl)

		for _, co := range t.Columns {
			if !co.AutoSet {
				iColumns = append(iColumns, co)
			}
			if !co.AutoSet && !co.Pk {
				uColumns = append(uColumns, co)
			}
		}

		var pkColumn table.Column
		for _, co := range t.Columns {
			if co.Pk {
				pkColumn = co
				break
			}
		}
		_ = tpl.Execute(&sqlBuf, struct {
			Schema        string
			TableName     string
			EntityName    string
			InsertColumns []table.Column
			UpdateColumns []table.Column
			Pk            table.Column
		}{
			Schema:        os.Getenv("DB_SCHEMA"),
			TableName:     t.Name,
			EntityName:    t.Meta.Name,
			InsertColumns: iColumns,
			UpdateColumns: uColumns,
			Pk:            pkColumn,
		})
		sqlStr := strings.TrimSpace(sqlBuf.String())
		sqlStr = strings.ReplaceAll(sqlStr, "`", "`"+" + "+`"`+"`"+`"`+" + "+"`")
		sqlBuf.Reset()
		sqlBuf.WriteString(`/**
* Generated by go-doudou ` + version.Release + `.
* You can edit it as your need.
*/
package dao
`)
		sqlBuf.WriteString("var " + strings.ToLower(t.Meta.Name) + "daosql=`" + sqlStr + "`\n")
		astutils.FixImport(sqlBuf.Bytes(), daofile)
	} else {
		log.Warnf("file %s already exists", daofile)
	}
	return nil
}

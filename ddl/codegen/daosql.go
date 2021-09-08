package codegen

import (
	"bytes"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/astutils"
	"github.com/unionj-cloud/go-doudou/ddl/table"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var daosqltmpl = `{{` + "`" + `{{` + "`" + `}}define "NoneZeroSet"{{` + "`" + `}}` + "`" + `}}
	{{- range $i, $co := .UpdateColumns}}
	{{- if or (eq $co.Meta.Type "time.Time") (eq $co.Meta.Type "*time.Time")}}
	{{` + "`" + `{{` + "`" + `}}- if .{{$co.Meta.Name}}{{` + "`" + `}}` + "`" + `}}
	` + "`" + `{{$co.Name}}` + "`" + `='{{` + "`" + `{{` + "`" + `}}.{{$co.Meta.Name}} | FormatTime{{` + "`" + `}}` + "`" + `}}',
	{{` + "`" + `{{` + "`" + `}}- end{{` + "`" + `}}` + "`" + `}}
	{{- else if or (eq $co.Meta.Type "bool") (eq $co.Meta.Type "*bool")}}
	{{` + "`" + `{{` + "`" + `}}- if .{{$co.Meta.Name}}{{` + "`" + `}}` + "`" + `}}
	` + "`" + `{{$co.Name}}` + "`" + `='{{` + "`" + `{{` + "`" + `}}.{{$co.Meta.Name}} | BoolToInt{{` + "`" + `}}` + "`" + `}}',
	{{` + "`" + `{{` + "`" + `}}- end{{` + "`" + `}}` + "`" + `}}
	{{- else}}
	{{` + "`" + `{{` + "`" + `}}- if .{{$co.Meta.Name}}{{` + "`" + `}}` + "`" + `}}
	` + "`" + `{{$co.Name}}` + "`" + `='{{` + "`" + `{{` + "`" + `}}.{{$co.Meta.Name}}{{` + "`" + `}}` + "`" + `}}',
	{{` + "`" + `{{` + "`" + `}}- end{{` + "`" + `}}` + "`" + `}}
	{{- end}}
	{{- end}}
{{` + "`" + `{{` + "`" + `}}end{{` + "`" + `}}` + "`" + `}}

{{` + "`" + `{{` + "`" + `}}define "InsertClause"{{` + "`" + `}}` + "`" + `}}
	{{- range $i, $co := .InsertColumns}}
	{{- if $i}},{{end}}
	{{- if eq $co.Meta.Type "time.Time" }}
	'{{` + "`" + `{{` + "`" + `}}.{{$co.Meta.Name}} | FormatTime{{` + "`" + `}}` + "`" + `}}'
	{{- else if eq $co.Meta.Type "*time.Time" }}
	{{` + "`" + `{{` + "`" + `}}- if .{{$co.Meta.Name}}{{` + "`" + `}}` + "`" + `}}
	'{{` + "`" + `{{` + "`" + `}}.{{$co.Meta.Name}} | FormatTime{{` + "`" + `}}` + "`" + `}}'
	{{` + "`" + `{{` + "`" + `}}- else{{` + "`" + `}}` + "`" + `}}
	null
	{{` + "`" + `{{` + "`" + `}}- end{{` + "`" + `}}` + "`" + `}}
	{{- else if eq $co.Meta.Type "bool" }}
	'{{` + "`" + `{{` + "`" + `}}.{{$co.Meta.Name}} | BoolToInt{{` + "`" + `}}` + "`" + `}}'
	{{- else if eq $co.Meta.Type "*bool" }}
	{{` + "`" + `{{` + "`" + `}}- if .{{$co.Meta.Name}}{{` + "`" + `}}` + "`" + `}}
	'{{` + "`" + `{{` + "`" + `}}.{{$co.Meta.Name}} | BoolToInt{{` + "`" + `}}` + "`" + `}}'
	{{` + "`" + `{{` + "`" + `}}- else{{` + "`" + `}}` + "`" + `}}
	null
	{{` + "`" + `{{` + "`" + `}}- end{{` + "`" + `}}` + "`" + `}}
	{{- else}}
	{{` + "`" + `{{` + "`" + `}}- if isNil .{{$co.Meta.Name}} {{` + "`" + `}}` + "`" + `}}
	null
	{{` + "`" + `{{` + "`" + `}}- else{{` + "`" + `}}` + "`" + `}}
	'{{` + "`" + `{{` + "`" + `}}.{{$co.Meta.Name}}{{` + "`" + `}}` + "`" + `}}'
	{{` + "`" + `{{` + "`" + `}}- end{{` + "`" + `}}` + "`" + `}}
	{{- end}}
	{{- end }}
{{` + "`" + `{{` + "`" + `}}end{{` + "`" + `}}` + "`" + `}}

{{` + "`" + `{{` + "`" + `}}define "Insert{{.DomainName}}"{{` + "`" + `}}` + "`" + `}}
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

{{` + "`" + `{{` + "`" + `}}define "Update{{.DomainName}}"{{` + "`" + `}}` + "`" + `}}
UPDATE ` + "`" + `{{.Schema}}` + "`" + `.` + "`" + `{{.TableName}}` + "`" + `
SET
	{{- range $i, $co := .UpdateColumns}}
	{{- if $i}},{{end}}
	` + "`" + `{{$co.Name}}` + "`" + `=:{{$co.Name}}
	{{- end }}
WHERE
    ` + "`" + `{{.Pk.Name}}` + "`" + ` =:{{.Pk.Name}}
{{` + "`" + `{{` + "`" + `}}end{{` + "`" + `}}` + "`" + `}}

{{` + "`" + `{{` + "`" + `}}define "Update{{.DomainName}}NoneZero"{{` + "`" + `}}` + "`" + `}}
UPDATE ` + "`" + `{{.Schema}}` + "`" + `.` + "`" + `{{.TableName}}` + "`" + `
SET
    {{` + "`" + `{{` + "`" + `}}Eval "NoneZeroSet" . | TrimSuffix ","{{` + "`" + `}}` + "`" + `}}
WHERE
    ` + "`" + `{{.Pk.Name}}` + "`" + `='{{` + "`" + `{{` + "`" + `}}.{{.Pk.Meta.Name}}{{` + "`" + `}}` + "`" + `}}'
{{` + "`" + `{{` + "`" + `}}end{{` + "`" + `}}` + "`" + `}}

{{` + "`" + `{{` + "`" + `}}define "Upsert{{.DomainName}}"{{` + "`" + `}}` + "`" + `}}
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

{{` + "`" + `{{` + "`" + `}}define "Upsert{{.DomainName}}NoneZero"{{` + "`" + `}}` + "`" + `}}
INSERT INTO ` + "`" + `{{.Schema}}` + "`" + `.` + "`" + `{{.TableName}}` + "`" + `
({{- range $i, $co := .InsertColumns}}
{{- if $i}},{{end}}
` + "`" + `{{$co.Name}}` + "`" + `
{{- end }})
VALUES ({{` + "`" + `{{` + "`" + `}}Eval "InsertClause" . | TrimSuffix ","{{` + "`" + `}}` + "`" + `}}) ON DUPLICATE KEY
UPDATE
		{{` + "`" + `{{` + "`" + `}}Eval "NoneZeroSet" . | TrimSuffix ","{{` + "`" + `}}` + "`" + `}}
{{` + "`" + `{{` + "`" + `}}end{{` + "`" + `}}` + "`" + `}}

{{` + "`" + `{{` + "`" + `}}define "Get{{.DomainName}}"{{` + "`" + `}}` + "`" + `}}
select *
from ` + "`" + `{{.Schema}}` + "`" + `.` + "`" + `{{.TableName}}` + "`" + `
where ` + "`" + `{{.Pk.Name}}` + "`" + ` = ?
{{` + "`" + `{{` + "`" + `}}end{{` + "`" + `}}` + "`" + `}}

{{` + "`" + `{{` + "`" + `}}define "Update{{.DomainName}}s"{{` + "`" + `}}` + "`" + `}}
UPDATE ` + "`" + `{{.Schema}}` + "`" + `.` + "`" + `{{.TableName}}` + "`" + `
SET
    {{- range $i, $co := .UpdateColumns}}
	{{- if $i}},{{end}}
	{{- if eq $co.Meta.Type "time.Time" }}
	` + "`" + `{{$co.Name}}` + "`" + `='{{` + "`" + `{{` + "`" + `}}.{{$co.Meta.Name}} | FormatTime{{` + "`" + `}}` + "`" + `}}'
	{{- else if eq $co.Meta.Type "*time.Time" }}
	{{` + "`" + `{{` + "`" + `}}- if .{{$co.Meta.Name}}{{` + "`" + `}}` + "`" + `}}
	` + "`" + `{{$co.Name}}` + "`" + `='{{` + "`" + `{{` + "`" + `}}.{{$co.Meta.Name}} | FormatTime{{` + "`" + `}}` + "`" + `}}'
	{{` + "`" + `{{` + "`" + `}}- else{{` + "`" + `}}` + "`" + `}}
	` + "`" + `{{$co.Name}}` + "`" + `=null
	{{` + "`" + `{{` + "`" + `}}- end{{` + "`" + `}}` + "`" + `}}
	{{- else if eq $co.Meta.Type "bool" }}
	` + "`" + `{{$co.Name}}` + "`" + `='{{` + "`" + `{{` + "`" + `}}.{{$co.Meta.Name}} | BoolToInt{{` + "`" + `}}` + "`" + `}}'
	{{- else if eq $co.Meta.Type "*bool" }}
	{{` + "`" + `{{` + "`" + `}}- if .{{$co.Meta.Name}}{{` + "`" + `}}` + "`" + `}}
	` + "`" + `{{$co.Name}}` + "`" + `='{{` + "`" + `{{` + "`" + `}}.{{$co.Meta.Name}} | BoolToInt{{` + "`" + `}}` + "`" + `}}'
	{{` + "`" + `{{` + "`" + `}}- else{{` + "`" + `}}` + "`" + `}}
	` + "`" + `{{$co.Name}}` + "`" + `=null
	{{` + "`" + `{{` + "`" + `}}- end{{` + "`" + `}}` + "`" + `}}
	{{- else}}
	{{` + "`" + `{{` + "`" + `}}- if isNil .{{$co.Meta.Name}} {{` + "`" + `}}` + "`" + `}}
	` + "`" + `{{$co.Name}}` + "`" + `=null
	{{` + "`" + `{{` + "`" + `}}- else{{` + "`" + `}}` + "`" + `}}
	` + "`" + `{{$co.Name}}` + "`" + `='{{` + "`" + `{{` + "`" + `}}.{{$co.Meta.Name}}{{` + "`" + `}}` + "`" + `}}'
	{{` + "`" + `{{` + "`" + `}}- end{{` + "`" + `}}` + "`" + `}}
	{{- end}}
	{{- end }}
WHERE
    {{` + "`" + `{{` + "`" + `}}.Where{{` + "`" + `}}` + "`" + `}}
{{` + "`" + `{{` + "`" + `}}end{{` + "`" + `}}` + "`" + `}}

{{` + "`" + `{{` + "`" + `}}define "Update{{.DomainName}}sNoneZero"{{` + "`" + `}}` + "`" + `}}
UPDATE ` + "`" + `{{.Schema}}` + "`" + `.` + "`" + `{{.TableName}}` + "`" + `
SET
    {{` + "`" + `{{` + "`" + `}}Eval "NoneZeroSet" . | TrimSuffix ","{{` + "`" + `}}` + "`" + `}}
WHERE
    {{` + "`" + `{{` + "`" + `}}.Where{{` + "`" + `}}` + "`" + `}}
{{` + "`" + `{{` + "`" + `}}end{{` + "`" + `}}` + "`" + `}}`

func GenDaoSql(domainpath string, t table.Table, folder ...string) error {
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
	daopath = filepath.Join(filepath.Dir(domainpath), df)
	if err = os.MkdirAll(daopath, os.ModePerm); err != nil {
		return errors.Wrap(err, "error")
	}

	daofile := filepath.Join(daopath, strings.ToLower(t.Meta.Name)+"daosql.go")
	if _, err = os.Stat(daofile); os.IsNotExist(err) {
		if f, err = os.Create(daofile); err != nil {
			return errors.Wrap(err, "error")
		}
		defer f.Close()

		funcMap = make(map[string]interface{})
		funcMap["ToSnake"] = strcase.ToSnake
		if tpl, err = template.New("daosql.tmpl").Funcs(funcMap).Parse(daosqltmpl); err != nil {
			return errors.Wrap(err, "error")
		}

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

		if err = tpl.Execute(&sqlBuf, struct {
			Schema        string
			TableName     string
			DomainName    string
			InsertColumns []table.Column
			UpdateColumns []table.Column
			Pk            table.Column
		}{
			Schema:        os.Getenv("DB_SCHEMA"),
			TableName:     t.Name,
			DomainName:    t.Meta.Name,
			InsertColumns: iColumns,
			UpdateColumns: uColumns,
			Pk:            pkColumn,
		}); err != nil {
			return errors.Wrap(err, "error")
		}
		sqlStr := strings.TrimSpace(sqlBuf.String())
		sqlStr = strings.ReplaceAll(sqlStr, "`", "`"+" + "+`"`+"`"+`"`+" + "+"`")
		sqlBuf.Reset()
		sqlBuf.WriteString("package dao\n")
		sqlBuf.WriteString("var " + strings.ToLower(t.Meta.Name) + "daosql=`" + sqlStr + "`\n")
		astutils.FixImport(sqlBuf.Bytes(), daofile)
	} else {
		log.Warnf("file %s already exists", daofile)
	}
	return nil
}

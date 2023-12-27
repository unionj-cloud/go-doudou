package mysql

var (
	createTable = `CREATE TABLE ` + "`" + `{{.Name}}` + "`" + ` (
{{- range $co := .Columns }}
` + "`" + `{{$co.Name}}` + "`" + ` {{$co.Type}} {{if $co.Nullable}}NULL{{else}}NOT NULL{{end}}{{if $co.Autoincrement}} AUTO_INCREMENT{{end}}{{if $co.Default}} DEFAULT {{$co.Default}}{{end}}{{if $co.Extra}} {{$co.Extra}}{{end}},
{{- end }}
PRIMARY KEY (` + "`" + `{{.Pk}}` + "`" + `))`

	dropTable = `DROP TABLE ` + "`" + `{{.Name}}` + "`;"

	alterTable = `{{define "change"}}
ALTER TABLE ` + "`" + `{{.Table}}` + "`" + `
CHANGE COLUMN ` + "`" + `{{if .OldName}}{{.OldName}}{{else}}{{.Name}}{{end}}` + "`" + ` ` + "`" + `{{.Name}}` + "`" + ` {{.Type}} {{if .Nullable}}NULL{{else}}NOT NULL{{end}}{{if .Autoincrement}} AUTO_INCREMENT{{end}}{{if .Default}} DEFAULT {{.Default}}{{end}}{{if .Extra}} {{.Extra}}{{end}};
{{end}}

{{define "add"}}
ALTER TABLE ` + "`" + `{{.Table}}` + "`" + `
ADD COLUMN ` + "`" + `{{.Name}}` + "`" + ` {{.Type}} {{if .Nullable}}NULL{{else}}NOT NULL{{end}}{{if .Autoincrement}} AUTO_INCREMENT{{end}}{{if .Default}} DEFAULT {{.Default}}{{end}}{{if .Extra}} {{.Extra}}{{end}};
{{end}}

{{define "drop"}}
ALTER TABLE ` + "`" + `{{.Table}}` + "`" + `
DROP COLUMN ` + "`" + `{{.Name}}` + "`" + `;
{{end}}
`

	insertInto = `INSERT INTO ` + "`" + `{{.Schema}}` + "`" + `.` + "`" + `{{.TableName}}` + "`" + `
({{- range $i, $co := .InsertColumns}}
{{- if $i}},{{end}}
` + "`" + `{{$co.Name}}` + "`" + `
{{- end }})
VALUES ({{- range $i, $co := .InsertColumns}}
	   {{- if $i}},{{end}}
	   ?
	   {{- end }});
`

	updateTable = `UPDATE ` + "`" + `{{.Schema}}` + "`" + `.` + "`" + `{{.TableName}}` + "`" + `
SET
	{{- range $i, $co := .UpdateColumns}}
	{{- if $i}},{{end}}
	` + "`" + `{{$co.Name}}` + "`" + `=?
	{{- end }}
WHERE
    ` + "`" + `{{.Pk.Name}}` + "`" + ` =?;
`

	deleteFrom = `DELETE FROM ` + "`" + `{{.Schema}}` + "`" + `.` + "`" + `{{.TableName}}` + "`" + `
WHERE
    ` + "`" + `{{.Pk.Name}}` + "`" + ` =?;
`

	selectFromById = `SELECT * FROM ` + "`" + `{{.Schema}}` + "`" + `.` + "`" + `{{.TableName}}` + "`" + `
WHERE
    ` + "`" + `{{.Pk.Name}}` + "`" + ` =?;
`
)

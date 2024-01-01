package postgres

var (
	createTable = `CREATE TABLE IF NOT EXISTS {{if .TablePrefix }}"{{.TablePrefix}}".{{end}}"{{.Name}}" (
{{- range $co := .Columns }}
"{{$co.Name}}" {{$co.Type}} {{if $co.Nullable}}NULL{{else}}NOT NULL{{end}}{{if $co.Default}} DEFAULT '{{$co.Default}}'{{end}},
{{- end }}
PRIMARY KEY ("{{.Pk}}"))
{{- if .Inherited }}
INHERITS ({{.Inherited}})
{{- end }};
    
{{- range $co := .Columns }}
COMMENT ON COLUMN {{if $.TablePrefix }}"{{$.TablePrefix}}".{{end}}"{{$.Name}}"."{{$co.Name}}" IS {{if $co.Comment}}$${{$co.Comment}}$${{else}}''{{end}};
{{- end }}
`

	dropTable = `DROP TABLE "{{.Name}}";`

	alterTable = `{{define "change"}}
{{if .OldName}}ALTER TABLE {{if .TablePrefix }}"{{.TablePrefix}}".{{end}}"{{.Table}}" RENAME COLUMN "{{.OldName}}" TO "{{.Name}}";{{end}}
ALTER TABLE {{if .TablePrefix }}"{{.TablePrefix}}".{{end}}"{{.Table}}" ALTER COLUMN "{{.Name}}" TYPE {{.Type}};
ALTER TABLE {{if .TablePrefix }}"{{.TablePrefix}}".{{end}}"{{.Table}}" ALTER COLUMN "{{.Name}}" {{if .Nullable}}DROP{{else}}SET{{end}} NOT NULL;
ALTER TABLE {{if .TablePrefix }}"{{.TablePrefix}}".{{end}}"{{.Table}}" ALTER COLUMN "{{.Name}}" {{if .Default}}SET DEFAULT {{.Default}}{{else}}DROP DEFAULT{{end}};
COMMENT ON COLUMN {{if .TablePrefix }}"{{.TablePrefix}}".{{end}}"{{.Table}}"."{{.Name}}" IS {{if .Comment}}$${{.Comment}}$${{else}}''{{end}};
{{end}}

{{define "add"}}
ALTER TABLE {{if .TablePrefix }}"{{.TablePrefix}}".{{end}}"{{.Table}}" ADD COLUMN "{{.Name}}" {{.Type}} {{if .Nullable}}NULL{{else}}NOT NULL{{end}} {{if .Default}}DEFAULT {{.Default}}{{end}};
COMMENT ON COLUMN {{if .TablePrefix }}"{{.TablePrefix}}".{{end}}"{{.Table}}"."{{.Name}}" IS {{if .Comment}}$${{.Comment}}$${{else}}''{{end}};
{{end}}

{{define "drop"}}
ALTER TABLE {{if .TablePrefix }}"{{.TablePrefix}}".{{end}}"{{.Table}}" DROP COLUMN "{{.Name}}";
{{end}}
`

	insertInto = `INSERT INTO {{if .TablePrefix }}"{{.TablePrefix}}".{{end}}"{{.TableName}}" 
({{- range $i, $co := .InsertColumns}}
{{- if $i}},{{end}}
"{{$co.Name}}"
{{- end }})
VALUES ({{- range $i, $co := .InsertColumns}}
	   {{- if $i}},{{end}}
	   ?
	   {{- end }}) RETURNING "{{.Pk.Name}}";
`

	insertIntoBatch = `INSERT INTO {{if .TablePrefix }}"{{.TablePrefix}}".{{end}}"{{.TableName}}" 
({{- range $i, $co := .InsertColumns}}
{{- if $i}},{{end}}
"{{$co.Name}}"
{{- end }})
VALUES {{- range $i, $ro := .Rows}}
{{- if $i}},{{end}}
({{- range $i, $co := $.InsertColumns}}
   {{- if $i}},{{end}}
   ?
{{- end }})
{{- end }};
`

	updateTable = `UPDATE {{if .TablePrefix }}"{{.TablePrefix}}".{{end}}"{{.TableName}}" 
SET
	{{- range $i, $co := .UpdateColumns}}
	{{- if $i}},{{end}}
	"{{$co.Name}}"=?
	{{- end }}
WHERE
    "{{.Pk.Name}}"=?;
`

	deleteFrom = `DELETE FROM {{if .TablePrefix }}"{{.TablePrefix}}".{{end}}"{{.TableName}}" 
WHERE
    "{{.Pk.Name}}"=?;
`

	selectFromById = `SELECT * FROM "{{.TableName}}" 
WHERE
    "{{.Pk.Name}}"=?;
`
)

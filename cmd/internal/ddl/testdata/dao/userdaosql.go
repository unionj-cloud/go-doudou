package dao

var userdaosql = `{{define "NoneZeroSet"}}
	{{- if .Name}}
	` + "`" + `name` + "`" + `=:name,
	{{- end}}
	{{- if .Phone}}
	` + "`" + `phone` + "`" + `=:phone,
	{{- end}}
	{{- if .Age}}
	` + "`" + `age` + "`" + `=:age,
	{{- end}}
	{{- if .No}}
	` + "`" + `no` + "`" + `=:no,
	{{- end}}
	{{- if .UniqueCol}}
	` + "`" + `unique_col` + "`" + `=:unique_col,
	{{- end}}
	{{- if .UniqueCol2}}
	` + "`" + `unique_col_2` + "`" + `=:unique_col_2,
	{{- end}}
	{{- if .School}}
	` + "`" + `school` + "`" + `=:school,
	{{- end}}
	{{- if .IsStudent}}
	` + "`" + `is_student` + "`" + `=:is_student,
	{{- end}}
	{{- if .Rule}}
	` + "`" + `rule` + "`" + `=:rule,
	{{- end}}
	{{- if .RuleType}}
	` + "`" + `rule_type` + "`" + `=:rule_type,
	{{- end}}
	{{- if .ArriveAt}}
	` + "`" + `arrive_at` + "`" + `=:arrive_at,
	{{- end}}
	{{- if .Status}}
	` + "`" + `status` + "`" + `=:status,
	{{- end}}
	{{- if .DeleteAt}}
	` + "`" + `delete_at` + "`" + `=:delete_at,
	{{- end}}
{{end}}

{{define "InsertUser"}}
INSERT INTO ` + "`" + `` + "`" + `.` + "`" + `user` + "`" + `
(
` + "`" + `id` + "`" + `,
` + "`" + `name` + "`" + `,
` + "`" + `phone` + "`" + `,
` + "`" + `age` + "`" + `,
` + "`" + `no` + "`" + `,
` + "`" + `unique_col` + "`" + `,
` + "`" + `unique_col_2` + "`" + `,
` + "`" + `school` + "`" + `,
` + "`" + `is_student` + "`" + `,
` + "`" + `rule` + "`" + `,
` + "`" + `rule_type` + "`" + `,
` + "`" + `arrive_at` + "`" + `,
` + "`" + `status` + "`" + `,
` + "`" + `delete_at` + "`" + `)
VALUES (
	   :id,
	   :name,
	   :phone,
	   :age,
	   :no,
	   :unique_col,
	   :unique_col_2,
	   :school,
	   :is_student,
	   :rule,
	   :rule_type,
	   :arrive_at,
	   :status,
	   :delete_at)
{{end}}

{{define "UpdateUser"}}
UPDATE ` + "`" + `` + "`" + `.` + "`" + `user` + "`" + `
SET
	` + "`" + `name` + "`" + `=:name,
	` + "`" + `phone` + "`" + `=:phone,
	` + "`" + `age` + "`" + `=:age,
	` + "`" + `no` + "`" + `=:no,
	` + "`" + `unique_col` + "`" + `=:unique_col,
	` + "`" + `unique_col_2` + "`" + `=:unique_col_2,
	` + "`" + `school` + "`" + `=:school,
	` + "`" + `is_student` + "`" + `=:is_student,
	` + "`" + `rule` + "`" + `=:rule,
	` + "`" + `rule_type` + "`" + `=:rule_type,
	` + "`" + `arrive_at` + "`" + `=:arrive_at,
	` + "`" + `status` + "`" + `=:status,
	` + "`" + `delete_at` + "`" + `=:delete_at
WHERE
    ` + "`" + `id` + "`" + ` =:id
{{end}}

{{define "UpdateUserNoneZero"}}
UPDATE ` + "`" + `` + "`" + `.` + "`" + `user` + "`" + `
SET
    {{Eval "NoneZeroSet" . | TrimSuffix ","}}
WHERE
    ` + "`" + `id` + "`" + `=:id
{{end}}

{{define "UpsertUser"}}
INSERT INTO ` + "`" + `` + "`" + `.` + "`" + `user` + "`" + `
(
` + "`" + `id` + "`" + `,
` + "`" + `name` + "`" + `,
` + "`" + `phone` + "`" + `,
` + "`" + `age` + "`" + `,
` + "`" + `no` + "`" + `,
` + "`" + `unique_col` + "`" + `,
` + "`" + `unique_col_2` + "`" + `,
` + "`" + `school` + "`" + `,
` + "`" + `is_student` + "`" + `,
` + "`" + `rule` + "`" + `,
` + "`" + `rule_type` + "`" + `,
` + "`" + `arrive_at` + "`" + `,
` + "`" + `status` + "`" + `,
` + "`" + `delete_at` + "`" + `)
VALUES (
        :id,
        :name,
        :phone,
        :age,
        :no,
        :unique_col,
        :unique_col_2,
        :school,
        :is_student,
        :rule,
        :rule_type,
        :arrive_at,
        :status,
        :delete_at) ON DUPLICATE KEY
UPDATE
		` + "`" + `name` + "`" + `=:name,
		` + "`" + `phone` + "`" + `=:phone,
		` + "`" + `age` + "`" + `=:age,
		` + "`" + `no` + "`" + `=:no,
		` + "`" + `unique_col` + "`" + `=:unique_col,
		` + "`" + `unique_col_2` + "`" + `=:unique_col_2,
		` + "`" + `school` + "`" + `=:school,
		` + "`" + `is_student` + "`" + `=:is_student,
		` + "`" + `rule` + "`" + `=:rule,
		` + "`" + `rule_type` + "`" + `=:rule_type,
		` + "`" + `arrive_at` + "`" + `=:arrive_at,
		` + "`" + `status` + "`" + `=:status,
		` + "`" + `delete_at` + "`" + `=:delete_at
{{end}}

{{define "UpsertUserNoneZero"}}
INSERT INTO ` + "`" + `` + "`" + `.` + "`" + `user` + "`" + `
(
` + "`" + `id` + "`" + `,
` + "`" + `name` + "`" + `,
` + "`" + `phone` + "`" + `,
` + "`" + `age` + "`" + `,
` + "`" + `no` + "`" + `,
` + "`" + `unique_col` + "`" + `,
` + "`" + `unique_col_2` + "`" + `,
` + "`" + `school` + "`" + `,
` + "`" + `is_student` + "`" + `,
` + "`" + `rule` + "`" + `,
` + "`" + `rule_type` + "`" + `,
` + "`" + `arrive_at` + "`" + `,
` + "`" + `status` + "`" + `,
` + "`" + `delete_at` + "`" + `)
VALUES (
        :id,
        :name,
        :phone,
        :age,
        :no,
        :unique_col,
        :unique_col_2,
        :school,
        :is_student,
        :rule,
        :rule_type,
        :arrive_at,
        :status,
        :delete_at) ON DUPLICATE KEY
UPDATE
		{{Eval "NoneZeroSet" . | TrimSuffix ","}}
{{end}}

{{define "GetUser"}}
select *
from ` + "`" + `` + "`" + `.` + "`" + `user` + "`" + `
where ` + "`" + `id` + "`" + ` = ?
{{end}}

{{define "UpdateUsers"}}
UPDATE ` + "`" + `` + "`" + `.` + "`" + `user` + "`" + `
SET
	` + "`" + `name` + "`" + `=:name,
	` + "`" + `phone` + "`" + `=:phone,
	` + "`" + `age` + "`" + `=:age,
	` + "`" + `no` + "`" + `=:no,
	` + "`" + `unique_col` + "`" + `=:unique_col,
	` + "`" + `unique_col_2` + "`" + `=:unique_col_2,
	` + "`" + `school` + "`" + `=:school,
	` + "`" + `is_student` + "`" + `=:is_student,
	` + "`" + `rule` + "`" + `=:rule,
	` + "`" + `rule_type` + "`" + `=:rule_type,
	` + "`" + `arrive_at` + "`" + `=:arrive_at,
	` + "`" + `status` + "`" + `=:status,
	` + "`" + `delete_at` + "`" + `=:delete_at
{{end}}

{{define "UpdateUsersNoneZero"}}
UPDATE ` + "`" + `` + "`" + `.` + "`" + `user` + "`" + `
SET
    {{Eval "NoneZeroSet" . | TrimSuffix ","}}
{{end}}`

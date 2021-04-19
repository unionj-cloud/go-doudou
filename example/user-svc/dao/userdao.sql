{{define "NoneZeroSet"}}
	{{- if .Name}}
	`name`='{{.Name}}',
	{{- end}}
	{{- if .Phone}}
	`phone`='{{.Phone}}',
	{{- end}}
	{{- if .Dept}}
	`dept`='{{.Dept}}',
	{{- end}}
	{{- if .DeleteAt}}
	`delete_at`='{{.DeleteAt | FormatTime}}',
	{{- end}}
{{end}}

{{define "InsertUser"}}
INSERT INTO `test`.`dd_user`
(
`id`,
`name`,
`phone`,
`dept`,
`delete_at`)
VALUES (
	   :id,
	   :name,
	   :phone,
	   :dept,
	   :delete_at)
{{end}}

{{define "UpdateUser"}}
UPDATE `test`.`dd_user`
SET
	`name`=:name,
	`phone`=:phone,
	`dept`=:dept,
	`delete_at`=:delete_at
WHERE
    `id` = ?
{{end}}

{{define "UpdateUserNoneZero"}}
UPDATE `test`.`dd_user`
SET
    {{Eval "NoneZeroSet" . | TrimSuffix ","}}
WHERE
    `id`='{{.Id}}'
{{end}}

{{define "UpsertUser"}}
INSERT INTO `test`.`dd_user`
(
`id`,
`name`,
`phone`,
`dept`,
`delete_at`)
VALUES (
        :id,
        :name,
        :phone,
        :dept,
        :delete_at) ON DUPLICATE KEY
UPDATE
		`name`=:name,
		`phone`=:phone,
		`dept`=:dept,
		`delete_at`=:delete_at
{{end}}

{{define "UpsertUserNoneZero"}}
INSERT INTO `test`.`dd_user`
(
`id`,
`name`,
`phone`,
`dept`,
`delete_at`)
VALUES (
		'{{.Id}}',
		'{{.Name}}',
		'{{.Phone}}',
		'{{.Dept}}',
		{{- if .DeleteAt}}
		'{{.DeleteAt | FormatTime}}'
		{{- else}}
		null
		{{- end}}) ON DUPLICATE KEY
UPDATE
		{{Eval "NoneZeroSet" . | TrimSuffix ","}}
{{end}}

{{define "GetUser"}}
select *
from `test`.`dd_user`
where `id` = ?
{{end}}

{{define "UpdateUsers"}}
UPDATE `test`.`dd_user`
SET
	`name`='{{.Name}}',
	`phone`='{{.Phone}}',
	`dept`='{{.Dept}}',
	{{- if .DeleteAt}}
	`delete_at`='{{.DeleteAt | FormatTime}}'
	{{- else}}
	`delete_at`=null
	{{- end}}
WHERE
    {{.Where}}
{{end}}

{{define "UpdateUsersNoneZero"}}
UPDATE `test`.`dd_user`
SET
    {{Eval "NoneZeroSet" . | TrimSuffix ","}}
WHERE
    {{.Where}}
{{end}}




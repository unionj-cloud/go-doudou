{{define "NoneZeroSet"}}
	{{- if .Name}}
	`name`='{{.Name}}',
	{{- end}}
	{{- if .Phone}}
	`phone`='{{.Phone}}',
	{{- end}}
	{{- if .DeleteAt}}
	`delete_at`='{{.DeleteAt | FormatTime}}',
	{{- end}}
{{end}}

{{define "InsertUser"}}
INSERT INTO `test`.`biz_user`
(
`id`,
`name`,
`phone`,
`delete_at`)
VALUES (
	   :id,
	   :name,
	   :phone,
	   :delete_at)
{{end}}

{{define "UpdateUser"}}
UPDATE `test`.`biz_user`
SET
	`name`=:name,
	`phone`=:phone,
	`delete_at`=:delete_at
WHERE
    `id` = ?
{{end}}

{{define "UpdateUserNoneZero"}}
UPDATE `test`.`biz_user`
SET
    {{Eval "NoneZeroSet" . | TrimSuffix ","}}
WHERE
    `id`='{{.Id}}'
{{end}}

{{define "UpsertUser"}}
INSERT INTO `test`.`biz_user`
(
`id`,
`name`,
`phone`,
`delete_at`)
VALUES (
        :id,
        :name,
        :phone,
        :delete_at) ON DUPLICATE KEY
UPDATE
		`name`=:name,
		`phone`=:phone,
		`delete_at`=:delete_at
{{end}}

{{define "UpsertUserNoneZero"}}
INSERT INTO `test`.`biz_user`
(
`id`,
`name`,
`phone`,
`delete_at`)
VALUES (
		'{{.Id}}',
		'{{.Name}}',
		'{{.Phone}}',
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
from `test`.`biz_user`
where `id` = ?
{{end}}

{{define "UpdateUsers"}}
UPDATE `test`.`biz_user`
SET
	`name`='{{.Name}}',
	`phone`='{{.Phone}}',
	{{- if .DeleteAt}}
	`delete_at`='{{.DeleteAt | FormatTime}}'
	{{- else}}
	`delete_at`=null
	{{- end}}
WHERE
    {{.Where}}
{{end}}

{{define "UpdateUsersNoneZero"}}
UPDATE `test`.`biz_user`
SET
    {{Eval "NoneZeroSet" . | TrimSuffix ","}}
WHERE
    {{.Where}}
{{end}}




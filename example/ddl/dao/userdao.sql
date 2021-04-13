{{define "NoneZeroSet"}}
	{{- if .Name}}
	`name`='{{.Name}}',
	{{- end}}
	{{- if .Phone}}
	`phone`='{{.Phone}}',
	{{- end}}
	{{- if .Age}}
	`age`='{{.Age}}',
	{{- end}}
	{{- if .No}}
	`no`='{{.No}}',
	{{- end}}
	{{- if .School}}
	`school`='{{.School}}',
	{{- end}}
	{{- if .IsStudent}}
	`is_student`='{{.IsStudent | BoolToInt}}',
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
`age`,
`no`,
`school`,
`is_student`,
`delete_at`)
VALUES (
	   :id,
	   :name,
	   :phone,
	   :age,
	   :no,
	   :school,
	   :is_student,
	   :delete_at)
{{end}}

{{define "UpdateUser"}}
UPDATE `test`.`biz_user`
SET
	`name`=:name,
	`phone`=:phone,
	`age`=:age,
	`no`=:no,
	`school`=:school,
	`is_student`=:is_student,
	`delete_at`=:delete_at
WHERE
    `id` = ?
{{end}}

{{define "UpdateUserNoneZero"}}
UPDATE `test`.`biz_user`
SET
    {{Eval "NoneZeroSet" . | TrimSuffix ","}}
WHERE
    `id`='{{.ID}}'
{{end}}

{{define "UpsertUser"}}
INSERT INTO `test`.`biz_user`
(
`id`,
`name`,
`phone`,
`age`,
`no`,
`school`,
`is_student`,
`delete_at`)
VALUES (
        :id,
        :name,
        :phone,
        :age,
        :no,
        :school,
        :is_student,
        :delete_at) ON DUPLICATE KEY
UPDATE
		`name`=:name,
		`phone`=:phone,
		`age`=:age,
		`no`=:no,
		`school`=:school,
		`is_student`=:is_student,
		`delete_at`=:delete_at
{{end}}

{{define "UpsertUserNoneZero"}}
INSERT INTO `test`.`biz_user`
(
`id`,
`name`,
`phone`,
`age`,
`no`,
`school`,
`is_student`,
`delete_at`)
VALUES (
		'{{.ID}}',
		'{{.Name}}',
		'{{.Phone}}',
		'{{.Age}}',
		'{{.No}}',
		'{{.School}}',
		'{{.IsStudent | BoolToInt}}',
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
	`age`='{{.Age}}',
	`no`='{{.No}}',
	`school`='{{.School}}',
	`is_student`='{{.IsStudent | BoolToInt}}',
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




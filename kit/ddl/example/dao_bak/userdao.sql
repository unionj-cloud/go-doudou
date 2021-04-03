{{define "NoneZeroSet"}}
    {{- if .Name}}
    name='{{.Name}}',
    {{- end}}
    {{- if .Phone}}
    phone='{{.Phone}}',
    {{- end}}
    {{- if .Age}}
    age='{{.Age}}',
    {{- end}}
    {{- if .No}}
    no='{{.No}}',
    {{- end}}
    {{- if .DeleteAt}}
    delete_at='{{.DeleteAt | FormatTime}}',
    {{- end}}
    {{- if .School}}
    school='{{.School}}',
    {{- end}}
{{end}}

{{define "InsertUser"}}
INSERT INTO users
(id,
 name,
 phone,
 age,
 no,
 delete_at,
 school)
VALUES (:id,
        :name,
        :phone,
        :age,
        :no,
        :delete_at,
        :school)
{{end}}

{{define "UpdateUser"}}
UPDATE users
SET
    name=:name,
    phone=:phone,
    age=:age,
    no=:no,
    delete_at=:delete_at,
    school=:school
WHERE
    id=:id
{{end}}

{{define "UpdateUserNoneZero"}}
UPDATE users
SET
    {{Eval "NoneZeroSet" . | TrimSuffix ","}}
WHERE
    id='{{.ID}}'
{{end}}

{{define "UpdateUsers"}}
UPDATE users
SET
    name='{{.Name}}',
    phone='{{.Phone}}',
    age='{{.Age}}',
    no='{{.No}}',
    {{if .DeleteAt}}
    delete_at='{{.DeleteAt | FormatTime}}',
    {{else}}
    delete_at=null,
    {{end}}
    school='{{.School}}'
WHERE
    {{.Where}}
{{end}}

{{define "UpdateUsersNoneZero"}}
UPDATE users
SET
    {{Eval "NoneZeroSet" . | TrimSuffix ","}}
WHERE
    {{.Where}}
{{end}}

{{define "UpsertUser"}}
INSERT INTO users
(id,
 name,
 phone,
 age,
 no,
 delete_at,
 school)
VALUES (:id,
        :name,
        :phone,
        :age,
        :no,
        :delete_at,
        :school) ON DUPLICATE KEY
UPDATE
    name=:name,
    phone=:phone,
    age=:age,
    no=:no,
    delete_at=:delete_at,
    school=:school
{{end}}

{{define "UpsertUserNoneZero"}}
INSERT INTO users
(id,
 name,
 phone,
 age,
 no,
 delete_at,
 school)
VALUES ('{{.ID}}',
        '{{.Name}}',
        '{{.Phone}}',
        '{{.Age}}',
        '{{.No}}',
        {{if .DeleteAt}}
        '{{.DeleteAt | FormatTime}}',
        {{else}}
        null,
        {{end}}
        '{{.School}}') ON DUPLICATE KEY
UPDATE
    {{Eval "NoneZeroSet" . | TrimSuffix ","}}
{{end}}

{{define "GetUser"}}
select * from users where id=?
{{end}}




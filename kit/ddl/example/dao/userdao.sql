-- {{define "UpsertUser"}}
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
    school=:school;
-- {{end}}

-- {{define "GetUser"}}
select * from users where id=?
-- {{end}}




-- {{define "UpsertUser"}}
INSERT INTO `test`.`users`
(
`id`,
`name`,
`phone`,
`age`,
`no`,
`school`,
`delete_at`)
VALUES (
        :id,
        :name,
        :phone,
        :age,
        :no,
        :school,
        :delete_at)
ON DUPLICATE KEY
    UPDATE
		   `name`=:name,
		   `phone`=:phone,
		   `age`=:age,
		   `no`=:no,
		   `school`=:school,
		   `delete_at`=:delete_at;
-- {{end}}

-- {{define "GetUser"}}
select *
from `test`.`users`
where id = ?
-- {{end}}




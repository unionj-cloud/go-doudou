package domain

import (
	"time"
)

//dd:table
type User struct {
	Id        int        `dd:"pk;auto;type:int"`
	Name      string     `dd:"type:varchar(255);default:'jack';index:name_phone_idx,2,asc"`
	Phone     string     `dd:"type:varchar(255);default:'13552053960';extra:comment 'mobile phone';index:name_phone_idx,1,asc"`
	Age       int        `dd:"type:int;index:age_idx,1,asc"`
	No        int        `dd:"type:int;unique:no_idx,1,asc"`
	School    *string    `dd:"type:varchar(255);default:'harvard';extra:comment 'school'"`
	IsStudent bool       `dd:"type:tinyint"`
	CreateAt  *time.Time `dd:"type:datetime;default:CURRENT_TIMESTAMP"`
	UpdateAt  *time.Time `dd:"type:datetime;default:CURRENT_TIMESTAMP;extra:on update CURRENT_TIMESTAMP"`
	DeleteAt  *time.Time `dd:"type:datetime"`
}

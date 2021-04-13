package domain

import "time"

//dd:table
type User struct {
	Id        int        `dd:"pk;auto;type:int(11)"`
	Name      string     `dd:"type:varchar(255);default:'jack';index:name_phone_idx,2,asc"`
	Phone     string     `dd:"type:varchar(255);default:'13552053960';extra:comment '手机号';index:name_phone_idx,1,asc"`
	Age       int        `dd:"type:int(11);index:age_idx,1,asc"`
	No        int        `dd:"type:int(11);unique:no_idx,1,asc"`
	School    *string    `dd:"type:varchar(255);default:'harvard';extra:comment '学校'"`
	IsStudent bool       `dd:"type:tinyint(4)"`
	DeleteAt  *time.Time `dd:"type:datetime"`
	CreateAt  *time.Time `dd:"type:datetime;default:CURRENT_TIMESTAMP"`
	UpdateAt  *time.Time `dd:"type:datetime;default:CURRENT_TIMESTAMP;extra:on update CURRENT_TIMESTAMP"`
}

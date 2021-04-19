package domain

import "time"

//dd:table
type User struct {
	Id       int        `dd:"pk;auto;type:int(11)"`
	Name     string     `dd:"type:varchar(255);extra:comment '真实姓名'"`
	Phone    string     `dd:"type:varchar(255);extra:comment '手机号';unique:phone_idx,1,asc"`
	Dept     string     `dd:"type:varchar(255);extra:comment '所属部门'"`
	UpdateAt *time.Time `dd:"type:datetime;default:CURRENT_TIMESTAMP;extra:on update CURRENT_TIMESTAMP"`
	DeleteAt *time.Time `dd:"type:datetime"`
	CreateAt *time.Time `dd:"type:datetime;default:CURRENT_TIMESTAMP"`
}

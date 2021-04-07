package domain

import "time"

//dd:table
type User struct {
	Id        int        `dd:"pk;auto;type:int(11)"`
	Name      string     `dd:"type:varchar(255);default:'jack';index:name_phone_idx,2,asc"`
	Phone     string     `dd:"type:varchar(255);default:'13893456784';extra:comment '手机号';index:name_phone_idx,1,asc"`
	Age       int64      `dd:"type:bigint(20);index:age_idx,1,asc"`
	No        int        `dd:"type:int(11);unique:no_idx,1,asc"`
	UpdateAt  *time.Time `dd:"type:datetime;default:CURRENT_TIMESTAMP;extra:on update CURRENT_TIMESTAMP"`
	DeleteAt  *time.Time `dd:"type:datetime"`
	CreateAt  *time.Time `dd:"type:datetime;default:CURRENT_TIMESTAMP"`
	School    *string    `dd:"type:varchar(255);default:'harvard';extra:comment '学校'"`
	IsStudent bool       `dd:"type:tinyint(4)"`
	Meters    *int64     `dd:"type:bigint(20) unsigned"`
	Height    *time.Time `dd:"type:datetime;default:CURRENT_TIMESTAMP"`
}

package entity

import "time"

//dd:table
type Material struct {
	Id       int        `dd:"pk;auto;type:int(11)"`
	Name     string     `dd:"type:varchar(45);extra:comment '原料名称'"`
	Amount   int        `dd:"type:int(11);extra:comment '原料单件克数'"`
	CreateAt *time.Time `dd:"type:datetime;default:CURRENT_TIMESTAMP"`
	UpdateAt *time.Time `dd:"type:datetime;default:CURRENT_TIMESTAMP;extra:on update CURRENT_TIMESTAMP"`
	DeleteAt *time.Time `dd:"type:datetime"`
	Price    float32    `dd:"type:decimal(10,2);extra:comment '原料单件进价'"`
}

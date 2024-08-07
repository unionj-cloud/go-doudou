package entity

import "time"

//dd:table
type Purchase struct {
	Id         int        `dd:"pk;auto;type:int(11)"`
	PurchaseAt *time.Time `dd:"type:datetime;extra:comment '采购时间'"`
	CreateAt   *time.Time `dd:"type:datetime;default:CURRENT_TIMESTAMP"`
	UpdateAt   *time.Time `dd:"type:datetime;default:CURRENT_TIMESTAMP;extra:on update CURRENT_TIMESTAMP"`
	DeleteAt   *time.Time `dd:"type:datetime"`
	ArriveAt   *time.Time `dd:"type:datetime;extra:comment '到货时间'"`
	Status     int8       `dd:"type:tinyint(4);extra:comment '0: 进行中
1: 完结
2: 取消'"`
	Note string `dd:"type:text;extra:comment '备注'"`
}

package domain

import "time"

//dd:table
type User struct {
	ID         int    `dd:"pk;auto"`
	Name       string `dd:"index:name_phone_idx,2;default:'jack'"`
	Phone      string `dd:"index:name_phone_idx,1;default:'13552053960';extra:comment '手机号'"`
	Age        int    `dd:"index;unsigned"`
	No         int    `dd:"type:int;unique"`
	UniqueCol  int    `dd:"type:int;unique:unique_col_idx,1"`
	UniqueCol2 int    `dd:"type:int;unique:unique_col_idx,2"`
	School     string `dd:"null;default:'harvard';extra:comment '学校'"`
	IsStudent  bool
	Rule       string `dd:"type:varchar(255);unique;extra:comment '链接匹配规则，匹配的链接采用该css规则来爬'"`
	RuleType   string `dd:"type:varchar(45);extra:comment '链接匹配规则类型，支持prefix前缀匹配和regex正则匹配'"`

	ArriveAt *time.Time `dd:"type:datetime;extra:comment '到货时间'"`
	Status   int8       `dd:"type:tinyint(4);extra:comment '0进行中
1完结
2取消'"`

	Base
}

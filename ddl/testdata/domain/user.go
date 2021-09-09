package domain

//dd:table
type User struct {
	ID         int    `dd:"pk;auto"`
	Name       string `dd:"index:name_phone_idx,2;default:'jack'"`
	Phone      string `dd:"index:name_phone_idx,1;default:'13552053960';extra:comment '手机号'"`
	Age        int    `dd:"index;unsigned"`
	No         int    `dd:"unique"`
	UniqueCol  int    `dd:"type:int;unique:unique_col_idx,1"`
	UniqueCol2 int    `dd:"type:int;unique:unique_col_idx,2"`
	School     string `dd:"null;default:'harvard';extra:comment '学校'"`
	IsStudent  bool

	Base
}

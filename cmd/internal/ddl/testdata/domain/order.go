package domain

//dd:table
type Order struct {
	ID     int `dd:"pk;auto"`
	Amount int64
	UserId int `dd:"type:int;fk:ddl_user,id,fk_ddl_user,ON DELETE CASCADE ON UPDATE NO ACTION"`

	Base
}

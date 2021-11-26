package domain

//dd:table
type Book struct {
	ID          int `dd:"pk;auto"`
	UserId      int `dd:"type:int;fk:ddl_publisher,id,fk_user,ON DELETE CASCADE ON UPDATE NO ACTION"`
	PublisherId int

	Base
}

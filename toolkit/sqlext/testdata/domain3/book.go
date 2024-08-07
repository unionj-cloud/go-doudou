package domain

//dd:table
type Book struct {
	ID          int `dd:"pk;auto"`
	UserId      int `dd:"type:int"`
	PublisherId int `dd:"fk:ddl_publisher,id,fk_publisher,ON DELETE CASCADE ON UPDATE NO ACTION"`

	Base
}

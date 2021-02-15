package models

//papi:table
type User struct {
	ID    int    `papi:"pk;auto"`
	Name  string `papi:"index:name_phone_idx,1,asc;default:wubin"`
	Phone string `papi:"index:name_phone_idx,2,desc;default:13552053960"`
	Age   int

	Base
}

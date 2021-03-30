package models

//dd:table
type User struct {
	ID    int    `dd:"pk;auto"`
	Name  string `dd:"index:name_phone_idx,2;default:wubin"`
	Phone string `dd:"index:name_phone_idx,1;default:13552053960"`
	Age   int    `dd:"index"`
	No    int    `dd:"unique"`

	Base
}

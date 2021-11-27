package domain

//dd:table
type Publisher struct {
	ID   int `dd:"pk;auto"`
	Name string

	Base
}

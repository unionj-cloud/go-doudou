package models

//papi:table
type User struct {
	ID   int `papi:pk;auto`
	Name string
	Age  int

	Base
}

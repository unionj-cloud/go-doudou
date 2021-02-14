package models

//go:generate ddl -file $GOFILE
type User struct {
	ID   int `papi:pk;auto`
	Name string
	Age  int

	Base
}

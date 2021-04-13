package vo

import "time"

//go:generate name -file $GOFILE
type Student struct {
	School
	Company

	Name string
	Age  int

	TestScore int

	IsPaid bool
}

type School struct {
	Name     string
	Address  string
	CreateAt time.Time
}

type Company struct {
	Name     string
	Biz      string
	CreateAt time.Time
}

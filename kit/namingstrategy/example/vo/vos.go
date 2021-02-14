package vo

import "time"

//go:generate namingstrategy -file $GOFILE
type Student struct {
	Name      string
	Age       int
	TestScore int

	School
}

type School struct {
	Name     string
	Address  string
	CreateAt time.Time
}

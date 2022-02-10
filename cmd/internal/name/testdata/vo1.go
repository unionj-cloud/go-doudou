package testdata

import "time"

// comment for alia age
type age int

type Event struct {
	Name      string
	EventType int
}

type TestAlias struct {
	Age    age
	School []struct {
		Name string
		Addr struct {
			Zip   string
			Block string
			Full  string
		}
	}
	EventChan chan Event
	SigChan   chan int
	Callback  func(string) bool
	CallbackN func(param string) bool
}

type ta TestAlias

type tt time.Time

type mm map[string]interface{}

type MyInter interface {
	Speak() error
}

type starM *time.Time

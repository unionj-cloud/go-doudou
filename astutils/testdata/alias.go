package main

import (
	"github.com/unionj-cloud/go-doudou/esutils"
	"time"
)

// comment for alia age
type age int

type Event struct {
	Name      string
	EventType int
}

type TestAlias struct {
	esutils.Base
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

type TT time.Time

type mm map[string]interface{}

type MyInter interface {
	Speak() error
}

type starM *time.Time

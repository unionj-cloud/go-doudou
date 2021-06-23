package vo

import (
	"github.com/unionj-cloud/go-doudou/esutils"
	"time"
)

type TestAlias1 struct {
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

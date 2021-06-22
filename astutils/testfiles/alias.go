package main

import "time"

// comment for alia age
type age int

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
}

type ta TestAlias

type tt time.Time

type mm map[string]interface{}

type MyInter interface {
	Speak() error
}

type starM *time.Time

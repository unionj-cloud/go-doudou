package main

import "fmt"

type Cat struct {
	Hobbies map[string]interface{}
	Sleep   func() bool
	Run     chan string
}

// eat execute eat behavior for Cat
func (c *Cat) eat(food string) (
	// not hungry
	full bool,
	// how feel
	mood string) {
	fmt.Println("eat " + food)
	return true, "happy"
}

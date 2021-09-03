package vo

import "fmt"

type Cat struct {
	Hobbies map[string]interface{}
	Sleep   func() bool
	Run     chan string
}

func (c *Cat) eat(food string) bool {
	fmt.Println("eat " + food)
	return true
}

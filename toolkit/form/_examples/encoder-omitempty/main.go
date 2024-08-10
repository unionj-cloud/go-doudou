package main

import (
	"fmt"
	"log"

	"github.com/go-playground/form/v4"
)

// Test ...
type Test struct {
	String string `form:",omitempty"`
	Array  []string
	Map    map[string]string
}

// use a single instance of Encoder, it caches struct info
var encoder *form.Encoder

func main() {
	var t Test

	encoder = form.NewEncoder()

	values, err := encoder.Encode(t)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("%#v\n", values)

	t.String = "String Val"
	t.Array = []string{"arr1"}

	values, err = encoder.Encode(t)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("%#v\n", values)
}

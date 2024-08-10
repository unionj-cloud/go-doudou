package main

import (
	"fmt"
	"log"

	"github.com/go-playground/form/v4"
)

// A ...
type A struct {
	Field string
}

// B ...
type B struct {
	A
	Field string
}

// use a single instance of Encoder, it caches struct info
var encoder *form.Encoder

func main() {

	type A struct {
		Field string
	}

	type B struct {
		A
		Field string
	}

	b := B{
		A: A{
			Field: "A Val",
		},
		Field: "B Val",
	}

	encoder = form.NewEncoder()

	v, err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("%#v\n", v)

	encoder.SetAnonymousMode(form.AnonymousSeparate)
	v, err = encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("%#v\n", v)
}

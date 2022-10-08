# Example of Router

These examples show you the usage of `router`. You can easily build a web application with it. Or you can make your own midwares such as custom logger, metrics, or any one you want.

### Basic example

This is just a quick introduction, view the [GoDoc](https://pkg.go.dev/github.com/fasthttp/router) for details.

Let's start with a trivial example:

```go
package main

import (
	"fmt"
	"log"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

// Index is the index handler
func Index(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, "Welcome!\n")
}

// Hello is the Hello handler
func Hello(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "hello, %s!\n", ctx.UserValue("name"))
}

// MultiParams is the multi params handler
func MultiParams(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "hi, %s, %s!\n", ctx.UserValue("name"), ctx.UserValue("word"))
}

// QueryArgs is used for uri query args test #11:
// if the req uri is /ping?name=foo, output: Pong! foo
// if the req uri is /piNg?name=foo, redirect to /ping, output: Pong!
func QueryArgs(ctx *fasthttp.RequestCtx) {
	name := ctx.QueryArgs().Peek("name")
	fmt.Fprintf(ctx, "Pong! %s\n", string(name))
}

func main() {
	r := router.New()
	r.GET("/", Index)
	r.GET("/hello/{name}", Hello)
	r.GET("/multi/{name}/{word}", MultiParams)
	r.GET("/ping", QueryArgs)

	log.Fatal(fasthttp.ListenAndServe(":8080", r.Handler))
}
```

package main

import (
	"fmt"
	"log"

	"github.com/unionj-cloud/go-doudou/framework/http/router"
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

// RegexParams is the params handler with regex validation
func RegexParams(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "hi, %s\n", ctx.UserValue("name"))
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
	r.GET("/regex/{name:[a-zA-Z]+}/test", RegexParams)
	r.GET("/optional/{name?:[a-zA-Z]+}/{word?}", MultiParams)
	r.GET("/ping", QueryArgs)

	log.Fatal(fasthttp.ListenAndServe(":8080", r.Handler))
}

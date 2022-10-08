/*
Package router is a trie based high performance HTTP request router.

A trivial example is:

 package main

 import (
     "fmt"
     "log"

     "github.com/unionj-cloud/go-doudou/framework/http/router"
 )

 func Index(ctx *fasthttp.RequestCtx) {
     fmt.Fprint(w, "Welcome!\n")
 }

 func Hello(ctx *fasthttp.RequestCtx) {
     fmt.Fprintf(w, "hello, %s!\n", ctx.UserValue("name"))
 }

 func main() {
     r := router.New()
     r.GET("/", Index)
     r.GET("/hello/{name}", Hello)

     log.Fatal(fasthttp.ListenAndServe(":8080", r.Handler))
 }

The router matches incoming requests by the request method and the path.
If a handler is registered for this path and method, the router delegates the
request to that function.
For the methods GET, POST, PUT, PATCH, DELETE and OPTIONS shortcut functions exist to
register handles, for all other methods router.Handle can be used.

The registered path, against which the router matches incoming requests, can
contain two types of parameters:
 Syntax    	Type
 {name}     	named parameter
 {name:*}	catch-all parameter

Named parameters are dynamic path segments. They match anything until the
next '/' or the path end:
 Path: /blog/{category}/{post}

 Requests:
  /blog/go/request-routers            match: category="go", post="request-routers"
  /blog/go/request-routers/           no match, but the router would redirect
  /blog/go/                           no match
  /blog/go/request-routers/comments   no match

Catch-all parameters match anything until the path end, including the
directory index (the '/' before the catch-all). Since they match anything
until the end, catch-all parameters must always be the final path element.
 Path: /files/{filepath:*}

 Requests:
  /files/                             match: filepath="/"
  /files/LICENSE                      match: filepath="/LICENSE"
  /files/templates/article.html       match: filepath="/templates/article.html"
  /files                              no match, but the router would redirect

The value of parameters is saved in ctx.UserValue(<key>), consisting
each of a key and a value. The slice is passed to the Handle func as a third
parameter.
To retrieve the value of a parameter,gets by the name of the parameter
 user := ctx.UserValue("user") // defined by {user} or {user:*}
*/
package router

# Example of Router

These examples show you the usage of `router`. You can easily build a web application with it. Or you can make your own midwares such as custom logger, metrics, or any one you want.

### Basic Authentication

This Basic Authentication example uses `simple-scrypt` for password hashing:

	go get -u github.com/elithrar/simple-scrypt

Password hashing is used so that if your data store is compromised, the attackers will only have access to hashed passwords, which (if the hash is not itself compromised) will not be able to revert to the original plain text password.

After you have hashed the password and stored the hash in a data store, you should throw away the original plain-text password.

The next time the user attempts to log in, their password will be safely hashed and compared to the saved hash. If the hashes match, then the user will be accepted.

Only use constant time comparison functions that are built into your hash library to compare secret strings like passwords or hashes to prevent timing attacks.

```go
package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/fasthttp/router"
	"github.com/elithrar/simple-scrypt"
	"github.com/valyala/fasthttp"
)

// basicAuth returns the username and password provided in the request's
// Authorization header, if the request uses HTTP Basic Authentication.
// See RFC 2617, Section 2.
func basicAuth(ctx *fasthttp.RequestCtx) (username, password string, ok bool) {
	auth := ctx.Request.Header.Peek("Authorization")
	if auth == nil {
		return
	}
	return parseBasicAuth(string(auth))
}

// parseBasicAuth parses an HTTP Basic Authentication string.
// "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==" returns ("Aladdin", "open sesame", true).
func parseBasicAuth(auth string) (username, password string, ok bool) {
	const prefix = "Basic "
	if !strings.HasPrefix(auth, prefix) {
		return
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}
	return cs[:s], cs[s+1:], true
}

// BasicAuth is the basic auth handler
func BasicAuth(h fasthttp.RequestHandler, requiredUser string, requiredPasswordHash []byte) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		// Get the Basic Authentication credentials
		user, password, hasAuth := basicAuth(ctx)
		
		// WARNING: 
		// DO NOT use plain-text passwords for real apps.
		// A simple string comparison using == is vulnerable to a timing attack.
		// Instead, use the hash comparison function found in your hash library.
		// This example uses scrypt, which is a solid choice for secure hashing:
		//   go get -u github.com/elithrar/simple-scrypt
		
		if hasAuth && user == requiredUser {
						
			// Uses the parameters from the existing derived key. Return an error if they don't match.
			err := scrypt.CompareHashAndPassword(requiredPasswordHash, []byte(password))

		    	if err != nil {
				
				// log error and request Basic Authentication again below.
				log.Fatal(err)
				
		    	} else {
				
				// Delegate request to the given handle
				h(ctx)
				return
				
		    	}
			
		}
		
		// Request Basic Authentication otherwise
		ctx.Error(fasthttp.StatusMessage(fasthttp.StatusUnauthorized), fasthttp.StatusUnauthorized)
		ctx.Response.Header.Set("WWW-Authenticate", "Basic realm=Restricted")
	})
}

// Index is the index handler
func Index(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, "Not protected!\n")
}

// Protected is the Protected handler
func Protected(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, "Protected!\n")
}

func main() {
	user := "gordon"
	pass := "secret!"
	
	// generate a hashed password from the password above:
	hashedPassword, err := scrypt.GenerateFromPassword([]byte(pass), scrypt.DefaultParams)
    	if err != nil {
        	log.Fatal(err)
    	}

	r := router.New()
	r.GET("/", Index)
	r.GET("/protected/", BasicAuth(Protected, user, hashedPassword))

	log.Fatal(fasthttp.ListenAndServe(":8080", r.Handler))
}
```

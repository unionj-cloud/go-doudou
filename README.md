## go-doudou

[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)
[![GoDoc](https://godoc.org/github.com/unionj-cloud/go-doudou?status.png)](https://godoc.org/github.com/unionj-cloud/go-doudou)
[![Build Status](https://travis-ci.com/unionj-cloud/go-doudou.svg?branch=main)](https://travis-ci.com/unionj-cloud/go-doudou)
[![Go](https://github.com/unionj-cloud/go-doudou/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/unionj-cloud/go-doudou/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/unionj-cloud/go-doudou/branch/main/graph/badge.svg?token=QRLPRAX885)](https://codecov.io/gh/unionj-cloud/go-doudou)
[![Go Report Card](https://goreportcard.com/badge/github.com/unionj-cloud/go-doudou)](https://goreportcard.com/report/github.com/unionj-cloud/go-doudou)
[![Release](https://img.shields.io/github/v/release/unionj-cloud/go-doudou?style=flat-square)](https://github.com/unionj-cloud/go-doudou)
[![Goproxy.cn](https://goproxy.cn/stats/github.com/unionj-cloud/go-doudou/badges/download-count.svg)](https://goproxy.cn)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Slack](https://img.shields.io/badge/Join%20Our%20Community-Slack-blue)](https://join.slack.com/t/go-doudou/shared_invite/zt-xzohc7ds-u7~aio6B8PELp5UtAdY~uw)
[![wakatime](https://wakatime.com/badge/user/852bcf22-8a37-460a-a8e2-115833174eba/project/57c830f7-e507-4cb1-9fd1-feedd96685f6.svg)](https://wakatime.com/badge/user/852bcf22-8a37-460a-a8e2-115833174eba/project/57c830f7-e507-4cb1-9fd1-feedd96685f6)

[EN](./README.md) [中文](./README_zh.md)  
go-doudou（doudou pronounce /dəudəu/）is a gossip protocol and OpenAPI 3.0 spec based decentralized microservice
framework. It supports monolith service application as well. Currently, it supports restful service only.



<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
### TOC

  - [Why?](#why)
    - [Background](#background)
    - [Reason](#reason)
    - [Result](#result)
  - [Philosophy](#philosophy)
  - [Features](#features)
  - [Overview](#overview)
  - [Recommend Architecture](#recommend-architecture)
  - [Golang Version Support](#golang-version-support)
  - [Install](#install)
  - [Upgrade](#upgrade)
  - [Usage](#usage)
  - [Hello World](#hello-world)
    - [Initialize project](#initialize-project)
    - [Define methods](#define-methods)
    - [Generate code](#generate-code)
    - [Run](#run)
    - [Deployment](#deployment)
      - [Build docker image and push to your repository](#build-docker-image-and-push-to-your-repository)
      - [Deploy](#deploy)
      - [Shutdown](#shutdown)
  - [Must Know](#must-know)
  - [Cors](#cors)
  - [Service register & discovery](#service-register--discovery)
  - [Client Load Balancing](#client-load-balancing)
    - [Simple Round-robin Load Balancing](#simple-round-robin-load-balancing)
    - [Smooth Weighted Round-robin Balancing](#smooth-weighted-round-robin-balancing)
  - [Rate Limit](#rate-limit)
    - [Usage](#usage-1)
    - [Memory based rate limiter Example](#memory-based-rate-limiter-example)
    - [Redis based rate limiter Example](#redis-based-rate-limiter-example)
  - [Bulkhead](#bulkhead)
    - [Usage](#usage-2)
    - [Example](#example)
  - [Circuit Breaker / Timeout / Retry](#circuit-breaker--timeout--retry)
    - [Usage](#usage-3)
    - [Example](#example-1)
  - [Log](#log)
    - [Usage](#usage-4)
    - [Example](#example-2)
    - [ELK stack](#elk-stack)
  - [Jaeger](#jaeger)
    - [Usage](#usage-5)
    - [Screenshot](#screenshot)
  - [Grafana / Prometheus](#grafana--prometheus)
    - [Usage](#usage-6)
    - [Screenshot](#screenshot-1)
  - [Configuration](#configuration)
  - [Example](#example-3)
  - [Notable tools](#notable-tools)
    - [name](#name)
    - [ddl](#ddl)
  - [TODO](#todo)
- [Credits](#credits)
- [Community](#community)
- [License](#license)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->


![screenshot-doc](./screenshot-doc.png)
![screenshot-registry](./screenshot-registry.png)


&nbsp;
### Must Know

There are some constraints or notable things when you define your methods as exposed apis for client in svc.go file.

1. Only support GET, POST, PUT, DELETE http methods. If method name starts with one of Get/Post/Put/Delete, http method
   will be one of GET/POST/PUT/DELETE. If method name doesn't start with any of them, default http method is POST.
2. First input parameter MUST be context.Context.
3. Only support golang [built-in types](https://golang.org/pkg/builtin/), map with string key, custom structs in vo
   package, corresponding slice and pointer types for input and output parameters. When go-doudou generate code and
   OpenAPI 3.0 spec, it will scan structs in vo package. If there is a struct from other package, the struct fields
   cannot be known by go-doudou.
4. As special cases, it supports multipart.FileHeader for uploading file as input parameter, supports os.File for
   downloading file as output parameter.
5. NOT support alias types as field of a struct.
6. NOT support func, channel, interface and anonymous struct type as input and output parameter.
7. When execute  `go-doudou svc http --handler` , existing code in handlerimpl.go won't be overwritten. If you added
   methods in svc.go, new code will be appended to handlerimpl.go.
8. When execute  `go-doudou svc http --handler` , existing code in handler.go will be overwritten, so don't modify
   handler.go file.
9. When execute  `go-doudou svc http`, only handler.go file will be overwritten and others will be checked if exists, if
   already exists, do nothing.  
&nbsp;
### Cors
Recommend to use [github.com/rs/cors](github.com/rs/cors) library. Here is example code.
```
corsOpts := cors.New(cors.Options{
    AllowedMethods: []string{
        http.MethodGet,
        http.MethodPost,
        http.MethodPut,
        http.MethodPatch,
        http.MethodDelete,
        http.MethodOptions,
        http.MethodHead,
    },

    AllowedHeaders: []string{
        "*",
    },
})

srv := ddhttp.NewDefaultHttpSrv()
srv.AddMiddleware(corsOpts.Handler, ddhttp.Tracing, ddhttp.Metrics, requestid.RequestIDHandler, handlers.CompressHandler, handlers.ProxyHeaders, ddhttp.Logger, ddhttp.Rest)
```  
&nbsp;
### Service register & discovery

Go-doudou supports monolith and microservices architecture.
Add below code to enable microservices architecture:
```go
err := registry.NewNode()
if err != nil {
    logrus.Panicln(fmt.Sprintf("%+v", err))
}
defer registry.Shutdown()
```  
&nbsp;
### Client Load Balancing

#### Simple Round-robin Load Balancing

```go
package main

import (
	"fmt"
	"github.com/ascarter/requestid"
	"github.com/gorilla/handlers"
	"github.com/sirupsen/logrus"
	ddconfig "github.com/unionj-cloud/go-doudou/framework/internal/config"
	ddhttp "github.com/unionj-cloud/go-doudou/framework/http"
	"github.com/unionj-cloud/go-doudou/framework/registry"
	service "ordersvc"
	"ordersvc/config"
	"ordersvc/transport/httpsrv"
	"usersvc/client"
)

func main() {
	ddconfig.InitEnv()
	conf := config.LoadFromEnv()

	err := registry.NewNode()
	if err != nil {
		logrus.Panicln(fmt.Sprintf("%+v", err))
	}
	defer registry.Shutdown()

	usersvcProvider := ddhttp.NewMemberlistServiceProvider("github.com/usersvc")
	usersvcClient := client.NewUsersvc(ddhttp.WithProvider(usersvcProvider))

	svc := service.NewOrdersvc(conf, nil, usersvcClient)

	handler := httpsrv.NewOrdersvcHandler(svc)
	srv := ddhttp.NewDefaultHttpSrv()
	srv.AddMiddleware(ddhttp.Tracing, ddhttp.Metrics, requestid.RequestIDHandler, handlers.CompressHandler, handlers.ProxyHeaders, ddhttp.Logger, ddhttp.Rest, ddhttp.Recover)
	srv.AddRoute(httpsrv.Routes(handler)...)
	srv.Run()
}
```

#### Smooth Weighted Round-robin Balancing

If environment variable GDD_MEM_WEIGHT is not set, local node weight will be calculated by health score and cpu idle
percent every GDD_MEM_WEIGHT_INTERVAL and gossip to remote nodes.

```go
package main

import (
	"fmt"
	"github.com/ascarter/requestid"
	"github.com/gorilla/handlers"
	"github.com/sirupsen/logrus"
	ddconfig "github.com/unionj-cloud/go-doudou/framework/internal/config"
	ddhttp "github.com/unionj-cloud/go-doudou/framework/http"
	"github.com/unionj-cloud/go-doudou/framework/registry"
	service "ordersvc"
	"ordersvc/config"
	"ordersvc/transport/httpsrv"
	"usersvc/client"
)

func main() {
	ddconfig.InitEnv()
	conf := config.LoadFromEnv()

	err := registry.NewNode()
	if err != nil {
		logrus.Panicln(fmt.Sprintf("%+v", err))
	}
	defer registry.Shutdown()

	usersvcProvider := ddhttp.NewSmoothWeightedRoundRobinProvider("github.com/usersvc")
	usersvcClient := client.NewUsersvc(ddhttp.WithProvider(usersvcProvider))

	svc := service.NewOrdersvc(conf, nil, usersvcClient)

	handler := httpsrv.NewOrdersvcHandler(svc)
	srv := ddhttp.NewDefaultHttpSrv()
	srv.AddMiddleware(ddhttp.Tracing, ddhttp.Metrics, requestid.RequestIDHandler, handlers.CompressHandler, handlers.ProxyHeaders, ddhttp.Logger, ddhttp.Rest, ddhttp.Recover)
	srv.AddRoute(httpsrv.Routes(handler)...)
	srv.Run()
}
```  
&nbsp;
### Rate Limit
#### Usage
There is a built-in [golang.org/x/time/rate](https://pkg.go.dev/golang.org/x/time/rate) based token-bucket rate limiter implementation
in `github.com/unionj-cloud/go-doudou/ratelimit/memrate` package with a `MemoryStore` struct for storing key and `Limiter` instance pairs.

If you don't like the built-in rate limiter implementation, you can implement `Limiter` interface by yourself.

You can pass an option function `memrate.WithTimer` to `memrate.NewLimiter` function to set a timer to each of 
`memrate.Limiter` instance returned for deleting the key in `keys` of the `MemoryStore` instance if it has been idle for `timeout` duration.

There is also a built-in [go-redis/redis_rate](https://github.com/go-redis/redis_rate) based redis GCRA rate limiter implementation.

#### Memory based rate limiter Example
Memory based rate limiter is stored in memory, only for single process.  

```go
package main

import (
	"context"
	"fmt"
	"github.com/ascarter/requestid"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/handlers"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/ratelimit"
	"github.com/unionj-cloud/go-doudou/ratelimit/redisrate"
	ddconfig "github.com/unionj-cloud/go-doudou/framework/internal/config"
	ddhttp "github.com/unionj-cloud/go-doudou/framework/http"
	"github.com/unionj-cloud/go-doudou/framework/logger"
	"github.com/unionj-cloud/go-doudou/framework/registry"
	"github.com/unionj-cloud/go-doudou/framework/tracing"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
	service "usersvc"
	"usersvc/config"
	"usersvc/transport/httpsrv"
)

func main() {
	ddconfig.InitEnv()
	conf := config.LoadFromEnv()

	logger.Init()

	err := registry.NewNode()
	if err != nil {
		logrus.Panicln(fmt.Sprintf("%+v", err))
	}
	defer registry.Shutdown()

	tracer, closer := tracing.Init()
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	svc := service.NewUsersvc(conf)

	handler := httpsrv.NewUsersvcHandler(svc)
	srv := ddhttp.NewDefaultHttpSrv()

	store := memrate.NewMemoryStore(func(_ context.Context, store *memrate.MemoryStore, key string) ratelimit.Limiter {
		return memrate.NewLimiter(10, 30, memrate.WithTimer(10*time.Second, func() {
			store.DeleteKey(key)
		}))
	})

	srv.AddMiddleware(ddhttp.Tracing, ddhttp.Metrics, requestid.RequestIDHandler, handlers.CompressHandler, handlers.ProxyHeaders, httpsrv.RateLimit(store), ddhttp.Logger, ddhttp.Rest, ddhttp.Recover)
	srv.AddRoute(httpsrv.Routes(handler)...)
	srv.Run()
}
```
Note: you need write your own http middleware to fit your needs. Here is an example below.
```go
// RateLimit limits rate based on memrate.MemoryStore
func RateLimit(store *memrate.MemoryStore) func(inner http.Handler) http.Handler {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.RemoteAddr[:strings.LastIndex(r.RemoteAddr, ":")]
			limiter := store.GetLimiter(key)
			if !limiter.Allow() {
				http.Error(w, "too many requests", http.StatusTooManyRequests)
				return
			}
			inner.ServeHTTP(w, r)
		})
	}
}
```

#### Redis based rate limiter Example
Redis based rate limiter is stored in redis, so it can be used for multiple processes to limit one key across cluster.  

```go
package main

import (
	"context"
	"fmt"
	"github.com/ascarter/requestid"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/handlers"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/ratelimit"
	"github.com/unionj-cloud/go-doudou/ratelimit/redisrate"
	ddconfig "github.com/unionj-cloud/go-doudou/framework/internal/config"
	ddhttp "github.com/unionj-cloud/go-doudou/framework/http"
	"github.com/unionj-cloud/go-doudou/framework/logger"
	"github.com/unionj-cloud/go-doudou/framework/registry"
	"github.com/unionj-cloud/go-doudou/framework/tracing"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
	service "usersvc"
	"usersvc/config"
	"usersvc/transport/httpsrv"
)

func main() {
	ddconfig.InitEnv()
	conf := config.LoadFromEnv()

	if logger.CheckDev() {
		logger.Init(logger.WithWritter(os.Stdout))
	} else {
		logger.Init(logger.WithWritter(io.MultiWriter(os.Stdout, &lumberjack.Logger{
			Filename:   filepath.Join(os.Getenv("LOG_PATH"), fmt.Sprintf("%s.log", ddconfig.GddServiceName.Load())),
			MaxSize:    5,  // Max megabytes before log is rotated
			MaxBackups: 10, // Max number of old log files to keep
			MaxAge:     7,  // Max number of days to retain log files
			Compress:   true,
		})))
	}

	if ddconfig.GddMode.Load() == "micro" {
		err := registry.NewNode()
		if err != nil {
			logrus.Panicln(fmt.Sprintf("%+v", err))
		}
		defer registry.Shutdown()
	}

	tracer, closer := tracing.Init()
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	svc := service.NewUsersvc(conf)

	handler := httpsrv.NewUsersvcHandler(svc)
	srv := ddhttp.NewDefaultHttpSrv()

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	fn := redisrate.LimitFn(func(ctx context.Context) ratelimit.Limit {
		return ratelimit.PerSecondBurst(10, 30)
	})

	srv.AddMiddleware(ddhttp.Tracing, ddhttp.Metrics,
		requestid.RequestIDHandler, handlers.CompressHandler, handlers.ProxyHeaders,
		httpsrv.RedisRateLimit(rdb, fn),
		ddhttp.Logger,
		ddhttp.Rest, ddhttp.Recover)
	srv.AddRoute(httpsrv.Routes(handler)...)
	srv.Run()
}
```
Note: you need write your own http middleware to fit your needs. Here is an example below.
```go
// RedisRateLimit limits rate based on redisrate.GcraLimiter
func RedisRateLimit(rdb redisrate.Rediser, fn redisrate.LimitFn) func(inner http.Handler) http.Handler {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.RemoteAddr[:strings.LastIndex(r.RemoteAddr, ":")]
			limiter := redisrate.NewGcraLimiterLimitFn(rdb, key, fn)
			if !limiter.Allow() {
				http.Error(w, "too many requests", http.StatusTooManyRequests)
				return
			}
			inner.ServeHTTP(w, r)
		})
	}
}
```  
&nbsp;
### Bulkhead
#### Usage
There is built-in [github.com/slok/goresilience](github.com/slok/goresilience) based bulkhead pattern support by BulkHead middleware in `github.com/unionj-cloud/go-doudou/framework/http` package.

```go
http.BulkHead(3, 10*time.Millisecond)
```

In above code, the first parameter `3` means the number of workers in the execution pool, the second parameter `10*time.Millisecond` 
means the max time an incoming request will wait to execute before being dropped its execution and return `429` response.

#### Example

```go
package main

import (
	"fmt"
	"github.com/ascarter/requestid"
	"github.com/gorilla/handlers"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	ddconfig "github.com/unionj-cloud/go-doudou/framework/internal/config"
	ddhttp "github.com/unionj-cloud/go-doudou/framework/http"
	"github.com/unionj-cloud/go-doudou/framework/logger"
	"github.com/unionj-cloud/go-doudou/framework/registry"
	"github.com/unionj-cloud/go-doudou/framework/tracing"
	"time"
	service "usersvc"
	"usersvc/config"
	"usersvc/transport/httpsrv"
)

func main() {
	ddconfig.InitEnv()
	conf := config.LoadFromEnv()

	logger.Init()

	if ddconfig.GddMode.Load() == "micro" {
		err := registry.NewNode()
		if err != nil {
			logrus.Panicln(fmt.Sprintf("%+v", err))
		}
		defer registry.Shutdown()
	}

	tracer, closer := tracing.Init()
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	svc := service.NewUsersvc(conf)

	handler := httpsrv.NewUsersvcHandler(svc)
	srv := ddhttp.NewDefaultHttpSrv()

	srv.AddMiddleware(ddhttp.Tracing, ddhttp.Metrics, ddhttp.BulkHead(1, 10*time.Millisecond), requestid.RequestIDHandler, handlers.CompressHandler, handlers.ProxyHeaders, ddhttp.Logger, ddhttp.Rest, ddhttp.Recover)
	srv.AddRoute(httpsrv.Routes(handler)...)
	srv.Run()
}
```  
&nbsp;
### Circuit Breaker / Timeout / Retry 
#### Usage
There is built-in [github.com/slok/goresilience](github.com/slok/goresilience) based Circuit Breaker / Timeout / Retry support in generated client code.
You don't need to do anything other than executing command: 
```shell
go-doudou svc http --handler -c go --doc
```  
The flag  `-c go` means generate go client code.
Then you will get two files in client folder: 
```shell
➜  usersvc git:(master) ✗ cd client    
➜  client git:(master) ✗ ll
total 32
-rw-r--r--  1 wubin1989  staff   7.9K  1 10 17:16 client.go
-rw-r--r--  1 wubin1989  staff   5.4K  1 10 17:16 clientproxy.go
```
For `client.go` file, all code will be overwritten each time you execute generation command.  
For `clientproxy.go` file, the existing code will not be changed, only new code will be appended. 

There is a default `goresilience.Runner` instance which has already been built-in circuit breaker, timeout and retry features for you, 
but if you need to customize it, you can pass `WithRunner(your_own_runner goresilience.Runner)` as `ProxyOption` parameter into 
`NewXXXClientProxy` function.

#### Example
```go 
package main

import (
	"fmt"
	"github.com/ascarter/requestid"
	"github.com/gorilla/handlers"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	ddconfig "github.com/unionj-cloud/go-doudou/framework/internal/config"
	ddhttp "github.com/unionj-cloud/go-doudou/framework/http"

	"github.com/unionj-cloud/go-doudou/framework/logger"
	"github.com/unionj-cloud/go-doudou/framework/registry"
	"github.com/unionj-cloud/go-doudou/framework/tracing"
	service "ordersvc"
	"ordersvc/config"
	"ordersvc/transport/httpsrv"
	"usersvc/client"
)

func main() {
	ddconfig.InitEnv()
	conf := config.LoadFromEnv()

	logger.Init()

	err := registry.NewNode()
	if err != nil {
		logrus.Panicln(fmt.Sprintf("%+v", err))
	}
	defer registry.Shutdown()

	tracer, closer := tracing.Init()
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	usersvcProvider := ddhttp.NewSmoothWeightedRoundRobinProvider("usersvc")
	usersvcClient := client.NewUsersvc(ddhttp.WithProvider(usersvcProvider))
	// if you don't need this resilience feature, you don't have to new the instance usersvcClientProxy.
	// you can just use usersvcClient.
	usersvcClientProxy := client.NewUsersvcClientProxy(usersvcClient)

	svc := service.NewOrdersvc(conf, nil, usersvcClientProxy)

	handler := httpsrv.NewOrdersvcHandler(svc)
	srv := ddhttp.NewDefaultHttpSrv()
	srv.AddMiddleware(ddhttp.Tracing, ddhttp.Metrics, requestid.RequestIDHandler, handlers.CompressHandler, handlers.ProxyHeaders, ddhttp.Logger, ddhttp.Rest, ddhttp.Recover)
	srv.AddRoute(httpsrv.Routes(handler)...)
	srv.Run()
}
```  
&nbsp;
### Log
#### Usage
There is a global `logrus.Entry` provided by `github.com/unionj-cloud/go-doudou/framework/logger` package. If `GDD_ENV` is set and is not set to `dev`,
it will be attached with some meta fields about service name, hostname, etc.

`logger` package implemented several exported package-level methods from `logrus`, so you can replace `logrus.Info()` with `logger.Info()` for example.
It also provided a `Init` function to help you configure `logrus.Logger` instance.

You can also configure log level by environment variable `GDD_LOG_LEVEL` and configure formatter type to `json` or `text` by environment variable `GDD_LOG_FORMAT`.

There are two built-in log related middlewares for you, `ddhttp.Metrics` and `ddhttp.Logger`. In short, `ddhttp.Metrics` is for printing brief log with limited 
information, while `ddhttp.Logger` is for printing detail log with request and response body, headers, opentracing span and some other information, and it only takes 
effect when environment variable `GDD_LOG_LEVEL` is set to `debug`.

#### Example
```go 
// you can use lumberjack to add log rotate feature to your service
logger.Init(logger.WithWritter(io.MultiWriter(os.Stdout, &lumberjack.Logger{
    Filename:   filepath.Join(os.Getenv("LOG_PATH"), fmt.Sprintf("%s.log", ddconfig.GddServiceName.Load())),
    MaxSize:    5,  // Max megabytes before log is rotated
    MaxBackups: 10, // Max number of old log files to keep
    MaxAge:     7,  // Max number of days to retain log files
    Compress:   true,
})))
```

#### ELK stack
`logger` package provided well support for ELK stack. To see example, please go to [go-doudou-guide](https://github.com/unionj-cloud/go-doudou-guide).

![elk](./elk.png)  
&nbsp;
### Jaeger
#### Usage
To add jaeger feature, you just need three steps:
1. Start jaeger
```shell
docker run -d --name jaeger \
  -p 6831:6831/udp \
  -p 16686:16686 \
  jaegertracing/all-in-one:1.29
```
2. Add two environment variables to your .env file
```shell
JAEGER_AGENT_HOST=localhost
JAEGER_AGENT_PORT=6831
```
3. Add three lines to your main function before new client and http server code
```go
tracer, closer := tracing.Init()
defer closer.Close()
opentracing.SetGlobalTracer(tracer)
```
Then your main function should like this
```go
package main

import (
	"fmt"
	"github.com/ascarter/requestid"
	"github.com/gorilla/handlers"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	ddconfig "github.com/unionj-cloud/go-doudou/framework/internal/config"
	ddhttp "github.com/unionj-cloud/go-doudou/framework/http"
	"github.com/unionj-cloud/go-doudou/framework/logger"
	"github.com/unionj-cloud/go-doudou/framework/registry"
	"github.com/unionj-cloud/go-doudou/framework/tracing"
	service "ordersvc"
	"ordersvc/config"
	"ordersvc/transport/httpsrv"
	"usersvc/client"
)

func main() {
	ddconfig.InitEnv()
	conf := config.LoadFromEnv()

	logger.Init()

	err := registry.NewNode()
	if err != nil {
		logrus.Panicln(fmt.Sprintf("%+v", err))
	}
	defer registry.Shutdown()

	tracer, closer := tracing.Init()
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	usersvcProvider := ddhttp.NewSmoothWeightedRoundRobinProvider("usersvc")
	usersvcClient := client.NewUsersvc(ddhttp.WithProvider(usersvcProvider))

	svc := service.NewOrdersvc(conf, nil, usersvcClient)

	handler := httpsrv.NewOrdersvcHandler(svc)
	srv := ddhttp.NewDefaultHttpSrv()
	srv.AddMiddleware(ddhttp.Tracing, ddhttp.Metrics, requestid.RequestIDHandler, handlers.CompressHandler, handlers.ProxyHeaders, ddhttp.Logger, ddhttp.Rest, ddhttp.Recover)
	srv.AddRoute(httpsrv.Routes(handler)...)
	srv.Run()
}
```
#### Screenshot
![jaeger1](./jaeger1.png)
![jaeger2](./jaeger2.png)  
&nbsp;
### Grafana / Prometheus
#### Usage
We implemented a service called `seed` for Prometheus service discovery based on [this blog](https://prometheus.io/blog/2018/07/05/implementing-custom-sd/).
Its source code is in [go-doudou-guide](https://github.com/unionj-cloud/go-doudou-guide) repo.

#### Screenshot
![grafana](./grafana.png)  
&nbsp;
### Configuration

Go-doudou use .env file to load environment variables to configure behaviors.

| Environment Variable    | Description                                                                                                                                                                                                                                                                        | Default   | Required |
| ----------------------- |------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------| -------- |
| GDD_BANNER              | Whether output banner to stdout or not, possible values are on and off.                                                                                                                                                                                                            | off       |          |
| GDD_BANNER_TEXT         |                                                                                                                                                                                                                                                                                    | Go-doudou |          |
| GDD_LOG_LEVEL           | Possible values are panic, fatal, error, warn, warning, info, debug, trace                                                                                                                                                                                                         | info      |          |
| GDD_LOG_FORMAT            | Set log format to text or json, accept values are text and json                                                                                                                                                                                                                    | text      |          |
| GDD_GRACE_TIMEOUT       | Graceful shutdown timeout for http server                                                                                                                                                                                                                                          | 15s       |          |
| GDD_WRITE_TIMEOUT       | Configure http.Server                                                                                                                                                                                                                                                              | 15s       |          |
| GDD_READ_TIMEOUT        | Configure http.Server                                                                                                                                                                                                                                                              | 15s       |          |
| GDD_IDLE_TIMEOUT        | Configure http.Server                                                                                                                                                                                                                                                              | 60s       |          |
| GDD_ROUTE_ROOT_PATH     | prefix GDD_ROUTE_ROOT_PATH to each of http api routes                                                                                                                                                                                                                              | ""        |          |
| GDD_SERVICE_NAME        | Service name that the node providing in the cluster.                                                                                                                                                                                                                               |           | Yes      |
| GDD_HOST                | Configure http.Server. Specifying host for the http server to listen on.                                                                                                                                                                                                           | ""        |          |
| GDD_PORT                | Configure http.Server. Specifying port for the http server to listen on.                                                                                                                                                                                                           | ""        |          |
| GDD_MODE                | Accept "mono" for monolith mode or "micro" for microservice mode                                                                                                                                                                                                                   |           |          |
| GDD_MANAGE_ENABLE       | Enable built-in api endpoints such as /go-doudou/doc, /go-doudou/openapi.json, /go-doudou/prometheus and /go-doudou/registry. Possible values are true and false.                                                                                                                  | false     |          |
| GDD_MANAGE_USER         | Http basic username for built-in api endpoints                                                                                                                                                                                                                                     | ""        |          |
| GDD_MANAGE_PASS         | Http basic password for built-in api endpoints                                                                                                                                                                                                                                     | ""        |          |
| GDD_MEM_SEED            | Seed address for join memberlist cluster. If empty or not set, this node will create a new cluster for other nodes to join.                                                                                                                                                        | ""        |          |
| GDD_MEM_NAME            | Only for dev and test use. Unique name of this node in cluster. if empty or not set, hostname will be used instead.                                                                                                                                                                | ""        |          |
| GDD_MEM_HOST            | Specify AdvertiseAddr attribute of memberlist config struct. if GDD_MEM_HOST starts with dot such as .seed-svc-headless.default.svc.cluster.local, it will be prefixed by hostname such as seed-2.seed-svc-headless.default.svc.cluster.local for supporting k8s stateful service. | ""        |          |
| GDD_MEM_PORT            | If empty or not set, an available port will be chosen randomly. Recommend specifying a port.                                                                                                                                                                                       | ""        |          |
| GDD_MEM_DEAD_TIMEOUT    | Dead node will be removed from node map if not received refute messages from it in GDD_MEM_DEAD_TIMEOUT duration                                                                                                                                                                   | 60s       |          |
| GDD_MEM_SYNC_INTERVAL   | Local node will synchronize states from other random node every GDD_MEM_SYNC_INTERVAL duration                                                                                                                                                                                     | 10s       |          |
| GDD_MEM_RECLAIM_TIMEOUT | Dead node will be replaced with new node with the same name but different full address in GDD_MEM_RECLAIM_TIMEOUT duration                                                                                                                                                         | 3s        |          |
| GDD_MEM_PROBE_INTERVAL | Do failure detecting every GDD_MEM_PROBE_INTERVAL duration                                                                                                                                                                                                                         | 5s        |          |
| GDD_MEM_PROBE_TIMEOUT | Probe fail if not receive ack message in GDD_MEM_PROBE_TIMEOUT duration                                                                                                                                                                                                            | 3s        |          |
| GDD_MEM_TCP_TIMEOUT | TCP request will timeout in GDD_MEM_TCP_TIMEOUT duration                                                                                                                                                                                                                           | 30s       |          |
| GDD_MEM_GOSSIP_NODES | Specify how many remote nodes you want to send gossip messages                                                                                                                                                                                                                     | 4         |          |
| GDD_MEM_GOSSIP_INTERVAL | Gossip messages in queue every GDD_MEM_GOSSIP_INTERVAL duration                                                                                                                                                                                                                    | 500ms     |          |
| GDD_MEM_SUSPICION_MULT | The multiplier for determining the time an inaccessible node is considered suspect before declaring it dead                                                                                                                                                                        | 6         |          |
| GDD_MEM_WEIGHT | Node weight for smooth weighted round-robin balancing                                                                                                                                                                                                                              | 0         |          |
| GDD_MEM_WEIGHT_INTERVAL | Node weight will be calculated every GDD_MEM_WEIGHT_INTERVAL                                                                                                                                                                                                                       | 5s        |          |
| GDD_RETRY_COUNT | Set resty client retry count                                                                                                                                                                                                                                                       | 0         |          |
&nbsp;
### Example

Please check [go-doudou-guide](https://github.com/unionj-cloud/go-doudou-guide)  
&nbsp;
### Notable tools

#### name

Command line tool for generating json tag of struct field. Please check [document ](./name/README.md).

#### ddl

DDL and dao layer generation command line tool based on [jmoiron/sqlx](https://github.com/jmoiron/sqlx). Please
check [document](./ddl/doc/README.md).  
&nbsp;
### TODO

Please reference [go-doudou kanban](https://github.com/unionj-cloud/go-doudou/projects/1)  
&nbsp;
## Credits
Here I give credit to [https://github.com/hashicorp/memberlist](https://github.com/hashicorp/memberlist) and all its contributors 
as go-doudou is relying on it to implement service register/discovery/fault tolerance feature.  

I also should give credit to [github.com/go-redis/redis_rate](github.com/go-redis/redis_rate) and all its contributors
as go-doudou is relying on it to implement redis based rate limit feature.
&nbsp;
## Community

Welcome to contribute to go-doudou by forking it and submitting pr or issues. If you like go-doudou, please give it a
star!

Slack invitation link: https://join.slack.com/t/go-doudou/shared_invite/zt-xzohc7ds-u7~aio6B8PELp5UtAdY~uw

Welcome to contact me from

- facebook: [https://www.facebook.com/bin.wu.94617999/](https://www.facebook.com/bin.wu.94617999/)
- twitter: [https://twitter.com/BINWU49205513](https://twitter.com/BINWU49205513)
- email: 328454505@qq.com
- wechat:  
  ![qrcode.png](qrcode.png)  
&nbsp;
## License

MIT

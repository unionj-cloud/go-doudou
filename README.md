## go-doudou
[![GoDoc](https://godoc.org/github.com/unionj-cloud/go-doudou?status.png)](https://godoc.org/github.com/unionj-cloud/go-doudou)
[![Build Status](https://travis-ci.com/unionj-cloud/go-doudou.svg?branch=main)](https://travis-ci.com/unionj-cloud/go-doudou)
[![codecov](https://codecov.io/gh/unionj-cloud/go-doudou/branch/main/graph/badge.svg?token=QRLPRAX885)](https://codecov.io/gh/unionj-cloud/go-doudou)
[![Go Report Card](https://goreportcard.com/badge/github.com/unionj-cloud/go-doudou)](https://goreportcard.com/report/github.com/unionj-cloud/go-doudou)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Funionj-cloud%2Fgo-doudou.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Funionj-cloud%2Fgo-doudou?ref=badge_shield)

[中文](./README_zh.md) [EN](./README.md)  

go-doudou（doudou pronounce /dəudəu/）is a gossip protocol and OpenAPI 3.0 spec based decentralized microservice framework. It supports monolith service application as well. It supports restful service only currently, but will support grpc in v2 version.



### Philosophy 

- Design First: We encourage designing your apis at the first place.
- Contract: We use OpenAPI 3.0 spec as a contract between server and client to reduce the communication cost between different dev teams and speed up development.
- Decentralization: We use gossip protocol to do service register and discovery to build a robust, scalable and decentralized service cluster. Thanks the awesome library memberlist by hashicorp.



### Features

- Low-code: design service interface to generate main function, routes, http handlers, service implementation stub, http client, OpenAPI 3.0 json spec and more.
- Support DNS address for service register and discovery
- Support monolith and microservices architecture
- Built-in client load balancing: currently only round-robin
- Built-in graceful shutdown
- Built-in live reloading by watching go files
- Built-in service apis documentation UI
- Built-in service registry UI
- Built-in prometheus middlewares: http_requests_total, response_status and http_response_time_seconds
- Built-in docker and k8s deployment support: dockerfile, deployment kind yaml file and statefulset kind yaml file
- Easy to learn, simple to use



### Support Golang Version

- go 1.13, 1.14, 1.15 with GO111MODULE=on
- go 1.16+
- < go 1.13: not test, maybe support

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

### TOC

  - [Install](#install)
  - [Hello World](#hello-world)
    - [Initialize project](#initialize-project)
    - [Define methods](#define-methods)
    - [Generate code](#generate-code)
    - [Run](#run)
    - [Deployment](#deployment)
      - [Build docker image and push to your repository](#build-docker-image-and-push-to-your-repository)
      - [Deploy](#deploy)
      - [Shutdown](#shutdown)
      - [Scale](#scale)
  - [Constraints](#constraints)
    - [Methods](#methods)
    - [Struct Parameters](#struct-parameters)
  - [Service register & discovery](#service-register--discovery)
  - [Client load balance](#client-load-balance)
  - [Example](#example)
  - [Notable tools](#notable-tools)
    - [name](#name)
    - [ddl](#ddl)
  - [TODO](#todo)
  - [Help](#help)
- [License](#license)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->



### Install

```shell
go get -v -u github.com/unionj-cloud/go-doudou/...@v0.6
```

### Hello World

#### Initialize project

```shell
➜  ~ go-doudou svc init helloworld
WARN[0000] Error loading .env file: open /Users/.env: no such file or directory 
1.16
helloworld
```
You can ignore the warning now.
```shell
➜  helloworld git:(master) ✗ ls -la -h
total 40
drwxr-xr-x   10 wubin1989  staff   320B  8 29 23:27 .
drwxr-xr-x+ 157 wubin1989  staff   4.9K  8 29 23:27 ..
-rw-r--r--    1 wubin1989  staff   2.0K  8 29 23:22 .env
drwxr-xr-x    5 wubin1989  staff   160B  8 29 23:26 .git
-rw-r--r--    1 wubin1989  staff   268B  8 29 23:22 .gitignore
drwxr-xr-x    6 wubin1989  staff   192B  8 29 23:27 .idea
-rw-r--r--    1 wubin1989  staff   707B  8 29 23:22 Dockerfile
-rw-r--r--    1 wubin1989  staff   442B  8 29 23:22 go.mod
-rw-r--r--    1 wubin1989  staff   253B  8 29 23:22 svc.go
drwxr-xr-x    3 wubin1989  staff    96B  8 29 23:22 vo
```
- Dockerfile：build docker image

- svc.go: design your rest apis by defining methods of Helloworld interface

- vo folder：define structs as data structure in http request body and response body, and also as OpenAPI 3.0 schemas

- .env: config file, go-doudou use it to load enviroment variables with GDD_ prefix

  

#### Define methods

There are some constraints, please read [Methods](#%E6%8E%A5%E5%8F%A3%E8%AE%BE%E8%AE%A1%E7%BA%A6%E6%9D%9F)和[Struct Parameters](#vo%E5%8C%85%E7%BB%93%E6%9E%84%E4%BD%93%E8%AE%BE%E8%AE%A1%E7%BA%A6%E6%9D%9F)

```go
package service

import (
	"context"
	"helloworld/vo"
)

type Helloworld interface {
	// You can define your service methods as your need. Below is an example.
	PageUsers(ctx context.Context, query vo.PageQuery) (code int, data vo.PageRet, err error)
}
```



#### Generate code

```shell
go-doudou svc http --handler -c go -o --doc
```
Let's see what are generated.
```shell
➜  helloworld git:(master) ✗ ls -la -h
total 328
drwxr-xr-x   20 wubin1989  staff   640B  8 31 12:34 .
drwxr-xr-x+ 157 wubin1989  staff   4.9K  8 31 12:36 ..
-rw-r--r--    1 wubin1989  staff   2.0K  8 29 23:45 .env
drwxr-xr-x    5 wubin1989  staff   160B  8 31 12:36 .git
-rw-r--r--    1 wubin1989  staff   268B  8 29 23:22 .gitignore
drwxr-xr-x    7 wubin1989  staff   224B  8 31 12:33 .idea
-rw-r--r--    1 wubin1989  staff   707B  8 29 23:22 Dockerfile
-rwxr-xr-x    1 wubin1989  staff    13K  8 31 12:35 app.log
drwxr-xr-x    3 wubin1989  staff    96B  8 29 23:44 client
drwxr-xr-x    3 wubin1989  staff    96B  8 29 23:44 cmd
drwxr-xr-x    3 wubin1989  staff    96B  8 29 23:44 config
drwxr-xr-x    3 wubin1989  staff    96B  8 29 23:44 db
-rw-r--r--    1 wubin1989  staff   536B  8 31 12:35 go.mod
-rw-r--r--    1 wubin1989  staff   115K  8 31 12:35 go.sum
-rwxr-xr-x    1 wubin1989  staff   1.9K  8 31 12:34 helloworld_openapi3.go
-rwxr-xr-x    1 wubin1989  staff   1.8K  8 31 12:34 helloworld_openapi3.json
-rw-r--r--    1 wubin1989  staff   253B  8 29 23:22 svc.go
-rw-r--r--    1 wubin1989  staff   413B  8 29 23:44 svcimpl.go
drwxr-xr-x    3 wubin1989  staff    96B  8 29 23:44 transport
drwxr-xr-x    3 wubin1989  staff    96B  8 29 23:22 vo
```
- helloworld_openapi3.json：OpenAPI 3.0 spec json documentation
- helloworld_openapi3.go: assgin OpenAPI 3.0 spec json string to a variable for serving online
- client：golang http client based on [resty](https://github.com/go-resty/resty)
- cmd：main.go file here
- config：config loading related
- db：function for connecting to database
- svcimpl.go：write your business logic here
- transport：http routes and handlers
- .env：put configs here



#### Run

Set GDD_MEM_SEED empty in .env file because there is no seed address before run our service now.

```shell
➜  helloworld git:(master) ✗ go run cmd/main.go
INFO[2021-08-31 21:35:47] Node 192.168.2.20 joined, supplying helloworld service 
WARN[2021-08-31 21:35:47] No seed found                                
INFO[2021-08-31 21:35:47] Memberlist created. Local node is Node 192.168.2.20, providing helloworld service at http://192.168.2.20:6060, memberlist port 50324 
 _____                     _                    _
|  __ \                   | |                  | |
| |  \/  ___   ______   __| |  ___   _   _   __| |  ___   _   _
| | __  / _ \ |______| / _` | / _ \ | | | | / _` | / _ \ | | | |
| |_\ \| (_) |        | (_| || (_) || |_| || (_| || (_) || |_| |
 \____/ \___/          \__,_| \___/  \__,_| \__,_| \___/  \__,_|
INFO[2021-08-31 21:35:47] ================ Registered Routes ================ 
INFO[2021-08-31 21:35:47] +-------------+--------+-------------------------+ 
INFO[2021-08-31 21:35:47] |    NAME     | METHOD |         PATTERN         | 
INFO[2021-08-31 21:35:47] +-------------+--------+-------------------------+ 
INFO[2021-08-31 21:35:47] | PageUsers   | POST   | /page/users             | 
INFO[2021-08-31 21:35:47] | GetDoc      | GET    | /go-doudou/doc          | 
INFO[2021-08-31 21:35:47] | GetOpenAPI  | GET    | /go-doudou/openapi.json | 
INFO[2021-08-31 21:35:47] | Prometheus  | GET    | /go-doudou/prometheus   | 
INFO[2021-08-31 21:35:47] | GetRegistry | GET    | /go-doudou/registry     | 
INFO[2021-08-31 21:35:47] +-------------+--------+-------------------------+ 
INFO[2021-08-31 21:35:47] =================================================== 
INFO[2021-08-31 21:35:47] Started in 431.269µs                         
INFO[2021-08-31 21:35:47] Http server is listening on :6060
```



#### Deployment

##### Build docker image and push to your repository

```shell
➜  helloworld git:(master) ✗ go-doudou svc push -r wubin1989
[+] Building 0.8s (13/13) FINISHED                                                                                                       
 => [internal] load build definition from Dockerfile                                                                                0.0s
 => => transferring dockerfile: 37B                                                                                                 0.0s
 => [internal] load .dockerignore                                                                                                   0.0s
 => => transferring context: 2B                                                                                                     0.0s
 => [internal] load metadata for docker.io/library/golang:1.13.4-alpine                                                             0.0s
 => [1/8] FROM docker.io/library/golang:1.13.4-alpine                                                                               0.0s
 => [internal] load build context                                                                                                   0.7s
 => => transferring context: 22.43MB                                                                                                0.6s
 => CACHED [2/8] WORKDIR /repo                                                                                                      0.0s
 => CACHED [3/8] ADD go.mod .                                                                                                       0.0s
 => CACHED [4/8] ADD go.sum .                                                                                                       0.0s
 => CACHED [5/8] ADD . ./                                                                                                           0.0s
 => CACHED [6/8] RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories                                   0.0s
 => CACHED [7/8] RUN apk add --no-cache bash tzdata                                                                                 0.0s
 => CACHED [8/8] RUN export GDD_VER=$(go list -mod=vendor -m -f '{{ .Version }}' github.com/unionj-cloud/go-doudou) && CGO_ENABLED  0.0s
 => exporting to image                                                                                                              0.0s
 => => exporting layers                                                                                                             0.0s
 => => writing image sha256:00365c58d0410d978aea462ec93323e20d879b15421e8eba29d8a17918660af8                                        0.0s
 => => naming to docker.io/library/helloworld                                                                                       0.0s

Use 'docker scan' to run Snyk tests against images to find vulnerabilities and learn how to fix them
The push refers to repository [docker.io/wubin1989/helloworld]
d0a9599b03e1: Pushed 
c3055fdf1a79: Layer already exists 
1c265a7f4c3e: Layer already exists 
f567cf5a5cf1: Layer already exists 
0b4acd902364: Layer already exists 
bbf9670b59e9: Layer already exists 
fdd6fb6fca5b: Layer already exists 
a17f85ec7605: Layer already exists 
2895b872dff5: Layer already exists 
eed8c158e67f: Layer already exists 
2033402d2275: Layer already exists 
77cae8ab23bf: Layer already exists 
v20210831125525: digest: sha256:5f75f7b43708d0619555f9bccbf0347e8db65319b83c65251015982ca6d23370 size: 2829
time="2021-08-31 12:55:53" level=info msg="image wubin1989/helloworld:v20210831125525 has been pushed successfully\n"
time="2021-08-31 12:55:53" level=info msg="k8s yaml has been created/updated successfully. execute command 'go-doudou svc deploy' to deploy service helloworld to k8s cluster\n"
```

then you should see there are two yaml files generated

```
➜  helloworld git:(master) ✗ ll
total 328
-rw-r--r--  1 wubin1989  staff   707B  8 29 23:22 Dockerfile
-rwxr-xr-x  1 wubin1989  staff    15K  8 31 12:55 app.log
drwxr-xr-x  3 wubin1989  staff    96B  8 29 23:44 client
drwxr-xr-x  3 wubin1989  staff    96B  8 29 23:44 cmd
drwxr-xr-x  3 wubin1989  staff    96B  8 29 23:44 config
drwxr-xr-x  3 wubin1989  staff    96B  8 29 23:44 db
-rw-r--r--  1 wubin1989  staff   536B  8 31 12:35 go.mod
-rw-r--r--  1 wubin1989  staff   115K  8 31 12:35 go.sum
-rw-r--r--  1 wubin1989  staff   817B  8 31 12:55 helloworld_deployment.yaml
-rwxr-xr-x  1 wubin1989  staff   1.9K  8 31 12:34 helloworld_openapi3.go
-rwxr-xr-x  1 wubin1989  staff   1.8K  8 31 12:34 helloworld_openapi3.json
-rw-r--r--  1 wubin1989  staff   867B  8 31 12:55 helloworld_statefulset.yaml
-rw-r--r--  1 wubin1989  staff   253B  8 29 23:22 svc.go
-rw-r--r--  1 wubin1989  staff   413B  8 29 23:44 svcimpl.go
drwxr-xr-x  3 wubin1989  staff    96B  8 29 23:44 transport
drwxr-xr-x  6 wubin1989  staff   192B  8 31 12:55 vendor
drwxr-xr-x  3 wubin1989  staff    96B  8 29 23:22 vo
```

- helloworld_deployment.yaml: k8s deploy file for stateless service, recommend for monolith architecture services
- helloworld_statefulset.yaml: k8s deploy file for stateful service, recommend for microservices  architecture services

##### Deploy

```shell
go-doudou svc deploy 
```

##### Shutdown

```shell
go-doudou svc shutdown
```



### Must Know

There are some constraints or notable things when you define your methods as exposed apis for client in svc.go file.

1. Only support GET, POST, PUT, DELETE http methods. If method name starts with one of Get/Post/Put/Delete, http method will be one of GET/POST/PUT/DELETE. If method name doesn't start with any of them, default http method is POST.
2. First input parameter MUST be context.Context.
3. Only support golang [built-in types](https://golang.org/pkg/builtin/), map with string key, custom structs in vo package, corresponding slice and pointer types for input and output parameters. When go-doudou generate code and OpenAPI 3.0 spec, it will scan structs in vo package. If there is a struct from other package, the struct fields cannot be known by go-doudou.
4. As special cases, it supports multipart.FileHeader for uploading file as input parameter, supports os.File for downloading file as output parameter.
5. NOT support alias types as field of a struct.
6. NOT support func, channel, interface and anonymous struct type as input and output parameter.
7. When execute  `go-doudou svc http --handler` , existing code in handlerimpl.go won't be overwritten. If you added or modified methods in svc.go, new code will be appended to handlerimpl.go.
8. When execute  `go-doudou svc http --handler` , existing code in handler.go will be overwritten, so don't modify handler.go file.
9. When execute  `go-doudou svc http`, before generating one file other than handler.go and handlerimpl.go files, go-doudou will check if it exists first, if already exists, do nothing.



### Service register & discovery

Go-doudou supports monolith and microservices architecture.
- `GDD_MODE=micro`：microservices architecture
- `GDD_MODE=mono`：monolith architecture 

There is service register code in main.go file.

```go
if ddconfig.GddMode.Load() == "micro" {
    node, err := registry.NewNode()
    if err != nil {
        logrus.Panicln(fmt.Sprintf("%+v", err))
    }
    logrus.Infof("Memberlist created. Local node is %s\n", node)
}
```
如果自己需要依赖其他服务，则除了需要把自己的服务注册到微服务集群之外，还需要加上实现服务发现的代码：

If dependent on other services, it should add service discovery code besides service register code.

```go
// service register
node, err := registry.NewNode()
if err != nil {
    logrus.Panicln(fmt.Sprintf("%+v", err))
}
logrus.Infof("%s joined cluster\n", node.String())

// call NewMemberlistServiceProvider to new a provider with name of the dependent service
usersvcProvider := ddhttp.NewMemberlistServiceProvider("usersvc", node)
// inject the provider into a client of the service
usersvcClient := client.NewUsersvc(client.WithProvider(usersvcProvider))

// inject the client into our service implementation instance
svc := service.NewOrdersvc(conf, conn, usersvcClient)
```



### Client load balance

Currently only one round robin strategy, welcome pr:)

```go
func (m *MemberlistServiceProvider) SelectServer() (string, error) {
	nodes, err := m.registry.Discover(m.name)
	if err != nil {
		return "", errors.Wrap(err, "SelectServer() fail")
	}
	next := int(atomic.AddUint64(&m.current, uint64(1)) % uint64(len(nodes)))
	m.current = uint64(next)
	selected := nodes[next]
	return selected.BaseUrl(), nil
}
```



### Example

Please check [go-doudou-guide](https://github.com/unionj-cloud/go-doudou-guide) 



### Notable tools

#### name

Command line tool for generating json tag of struct field. [Document](./name/README.md)



#### ddl

基于[jmoiron/sqlx](https://github.com/jmoiron/sqlx) 实现的同步数据库表结构和Go结构体的工具。还可以生成dao层代码。
[查看文档](./ddl/doc/README.md)



### TODO
Please reference [go-doudou kanban](https://github.com/unionj-cloud/go-doudou/projects/1)



### Help

希望大家跟我一起完善go-doudou，欢迎提pr和issue，欢迎扫码加作者微信提意见和需求。  
![qrcode.png](qrcode.png) 

社区钉钉群二维码，群号：31405977

![dingtalk.png](dingtalk.png)





## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Funionj-cloud%2Fgo-doudou.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Funionj-cloud%2Fgo-doudou?ref=badge_large)

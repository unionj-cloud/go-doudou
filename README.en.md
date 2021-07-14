## go-doudou
[![GoDoc](https://godoc.org/github.com/unionj-cloud/go-doudou?status.png)](https://godoc.org/github.com/unionj-cloud/go-doudou)
[![Build Status](https://travis-ci.com/unionj-cloud/go-doudou.svg?branch=main)](https://travis-ci.com/unionj-cloud/go-doudou)
[![codecov](https://codecov.io/gh/unionj-cloud/go-doudou/branch/main/graph/badge.svg?token=QRLPRAX885)](https://codecov.io/gh/unionj-cloud/go-doudou)
[![Go Report Card](https://goreportcard.com/badge/github.com/unionj-cloud/go-doudou)](https://goreportcard.com/report/github.com/unionj-cloud/go-doudou)

go-doudou(doudou pronounced doudou/dəudəu/)is a golang decentralized microservice agile development framework 
based on the gossip protocol for service registration and discovery,openapi3.0 specification for interface definition.  
go-doudou uses a set of command line tools to help developers quickly initialize one or a set of RESTful services.
Designing a set of apis by defining methods in the interface type,
then generating the main function to run your service,router and corresponding handler,client go code through command line tools.  
go-doudou advocates design first,by pre-designing and definition interfaces to generate code.
After modifying the definition,recovering or incrementally generating code to achieve rapid development.  
go-doudou canonizes the spirit of contract,statement the interface,regulates the cooperation between service providers and consumers,
and prompts the whole develop team to improve delivery efficiency through openapi3.0 protocol.  
go-doudou works for helping developers build a decentralized microservice system,connecting services in the cluster through the gossip protocol,
and using client load balancing to call other services.  


### Install

```shell
go get -v -u github.com/unionj-cloud/go-doudou/...@v0.4.6
```

### Usage

1. Initialize the project，taking the auth service as an example:
```shell
go-doudou svc init auth
```
Then the following project structure will be generated
```shell
➜  auth git:(master) ✗ ll
total 24
-rw-r--r--  1 wubin1989  staff   372B  7  2 17:20 Dockerfile
-rw-r--r--  1 wubin1989  staff   399B  7  2 17:20 go.mod
-rw-r--r--  1 wubin1989  staff   241B  7  2 17:20 svc.go
drwxr-xr-x  3 wubin1989  staff    96B  7  2 17:20 vo
```
- Dockerfile：Used to generate docker image
- svc.go：Interface design file,including interface type in it,you can defines methods in it
- vo folder：Your struct is defined inside, as the input and output parameters of the interface, and is also used to generate the schema in the openapi3.0 specification

2. Define method in interface in svc.go.If necessary, define the input and output struct in package vo.
   此处略，见下文的[接口设计约束](#%E6%8E%A5%E5%8F%A3%E8%AE%BE%E8%AE%A1%E7%BA%A6%E6%9D%9F)和[vo包结构体设计约束](#vo%E5%8C%85%E7%BB%93%E6%9E%84%E4%BD%93%E8%AE%BE%E8%AE%A1%E7%BA%A6%E6%9D%9F)
   

3. Generate HTTP interface code.
```shell
go-doudou svc http --handler -c go -o --doc
```
Some new folders have been added.
```shell
➜  auth git:(master) ✗ ls -la -h                  
total 280
drwxr-xr-x  17 wubin1989  staff   544B  7  2 17:43 .
drwxr-xr-x  11 wubin1989  staff   352B  7  2 17:40 ..
-rw-r--r--   1 wubin1989  staff   413B  7  2 17:43 .env
drwxr-xr-x   5 wubin1989  staff   160B  7  2 17:42 .git
-rw-r--r--   1 wubin1989  staff   268B  7  2 17:40 .gitignore
-rw-r--r--   1 wubin1989  staff   372B  7  2 17:40 Dockerfile
-rwxr-xr-x   1 wubin1989  staff   1.8K  7  2 17:40 auth_openapi3.json
drwxr-xr-x   3 wubin1989  staff    96B  7  2 17:40 client
drwxr-xr-x   3 wubin1989  staff    96B  7  2 17:40 cmd
drwxr-xr-x   4 wubin1989  staff   128B  7  2 17:40 config
drwxr-xr-x   3 wubin1989  staff    96B  7  2 17:40 db
-rw-r--r--   1 wubin1989  staff   614B  7  2 17:42 go.mod
-rw-r--r--   1 wubin1989  staff   111K  7  2 17:42 go.sum
-rw-r--r--   1 wubin1989  staff   241B  7  2 17:40 svc.go
-rw-r--r--   1 wubin1989  staff   369B  7  2 17:40 svcimpl.go
drwxr-xr-x   3 wubin1989  staff    96B  7  2 17:40 transport
drwxr-xr-x   3 wubin1989  staff    96B  7  2 17:40 vo
```
- auth_openapi3.json：Openapi3.0 standard json format interface document
- client：Client code contains golang interface , encapsulates the [resty module](https://github.com/go-resty/resty)
- cmd：Service start entry,You need to create dependent components or third-party service client instances in the main method and inject them into the service instances of this project.
- config：Related configuration file 
- db：Generate database connection
- svcimpl.go：Implement service logic here.
- transport：Contains generated http routes and handlers
- .env：Define environment variables.  


4. Delete configuration item GDD_SEED's value in the .env file,since there is no seed yet.
   

5. Run your service.
```shell
➜  auth git:(master) ✗ go run cmd/main.go
INFO[0000] Node wubindeMacBook-Pro.local joined, supplying auth service 
WARN[0000] No seed found                                
INFO[0000] Memberlist created. Local node is Node wubindeMacBook-Pro.local, providing auth service at 192.168.101.6, memberlist port 57157, service port 6060 
 _____                     _                    _
|  __ \                   | |                  | |
| |  \/  ___   ______   __| |  ___   _   _   __| |  ___   _   _
| | __  / _ \ |______| / _` | / _ \ | | | | / _` | / _ \ | | | |
| |_\ \| (_) |        | (_| || (_) || |_| || (_| || (_) || |_| |
 \____/ \___/          \__,_| \___/  \__,_| \__,_| \___/  \__,_|
INFO[2021-07-02 17:46:53] ================ Registered Routes ================ 
INFO[2021-07-02 17:46:53] +-----------+--------+-----------------+     
INFO[2021-07-02 17:46:53] |   NAME    | METHOD |     PATTERN     |     
INFO[2021-07-02 17:46:53] +-----------+--------+-----------------+     
INFO[2021-07-02 17:46:53] | PageUsers | POST   | /auth/pageusers |     
INFO[2021-07-02 17:46:53] +-----------+--------+-----------------+     
INFO[2021-07-02 17:46:53] =================================================== 
INFO[2021-07-02 17:46:53] Started in 468.696µs                         
INFO[2021-07-02 17:46:53] Http server is listening on :6060 
```  

Starting from step 6 is the steps related to deploying services, which requires a local docker environment to connect to the local or remote k8s service.  


6. make a image.
```shell
go-doudou svc push -r yourprivaterepositoryaddress
```  


7. Deploy on k8s
```shell
go-doudou svc deploy 
```  


8. Close service
```shell
go-doudou svc shutdown
```  


9. Scale service
```shell
go-doudou svc scale -n 3
```

### Notice

only http restful interface is supported yet, grpc and protobuffer are not supported


### interface design specification

1. Supports four http request methods: Post, Get, Delete, Put,corresponding to the interface method name, Post request by default.If the method name starts with Post/Get/Delete/Put, the http request method is one of the corresponding post/get/delete/put.
2. The type of the first input parameter is context.Context,which you don't need to change.You can use this parameter to achieve some effects. For example, when the client cancels the request, the processing logic can be stopped in time to save server resources.
3. The input and output parameters' type only support the built-in types of the Go language, the key is string type of map type, the custom struct in the vo package, and the corresponding slice type and pointer type of the above types.When go-doudou generates code and openapi documents, it scans the struct in the vo package. If the input and output parameters of the interface use the struct in a package other than the vo package, go-doudou cannot scan the fields of the structure.
4. In particular, the input parameter also supports the multipart.FileHeader type for file upload. The output also supports os.File type for file download.
5. func type, channel type, interface type and anonymous struct are not supported
6. Since the methods related to fetching Form parameters in the net/http package of go, such as FormValue, the parameter values obtained are all of string type. go-doudou uses the cobra and viper author spf13 [cast](https:/ /github.com/spf13/cast) module for type conversion,
   The code in the generated handlerimpl.go file may report a compilation error in the parsing of the form parameters. You can submit [issue](https://github.com/unionj-cloud/go-doudou/issues) to go-doudou, You can also modify it manually.
   When the interface method in svc.go is added, deleted, changed and the code generation command `go-doudou svc http --handler -c go -o --doc` is re-executed, the code in the handlerimpl.go file is generated incrementally. That is, the code generated before and the code manually modified by yourself will not be overwritten
7. The code in the handler.go file will be regenerated every time the go-doudou svc http command is executed, please do not manually modify the code inside.
8. Except for handler.go and handlerimpl.go, all files are first judged whether they exist, and then they are generated if they do not exist, otherwise, do nothing.

### package vo design specification

1. Struct's field type, only support Go language [built-in type](https://golang.org/pkg/builtin/), key is string type of map type, custom struct in vo package, **anonymous struct ** and the corresponding slice type and pointer type of the above types.
2. Func type, channel type, interface type are not supported.
3. Structure field type, type alias is not supported.

### Service registration and discovery
go-doudou supports both monolithic mode and microservice mode, which can be configured in the form of environment variables.
- `GDD_MODE=micro`：microservice mode
- `GDD_MODE=mono`：monolithic mode  
The generated cmd/main.go file has the following code：  
```go
if ddconfig.GddMode.Load() == "micro" {
    node, err := registry.NewNode()
    if err != nil {
        logrus.Panicln(fmt.Sprintf("%+v", err))
    }
    logrus.Infof("Memberlist created. Local node is %s\n", node)
}
```
You need to register your own service through the `registry.NewNode()` method，when other services depend on you, 
If you need to rely on other services, in addition to registering your services to the microservice cluster, you also need to add code to implement service discovery:
```go
// Register yourself and join the cluster
node, err := registry.NewNode()
if err != nil {
    logrus.Panicln(fmt.Sprintf("%+v", err))
}
logrus.Infof("%s joined cluster\n", node.String())

// Create a usersvc service provider when you need to rely on usersvc service
usersvcProvider := ddhttp.NewMemberlistServiceProvider("usersvc", node)
// Inject the provider of the usersvc service into the client instance of the usersvc service
usersvcClient := client.NewUsersvc(client.WithProvider(usersvcProvider))

// Inject the client instance of the usersvc service into your own service instance
svc := service.NewOrdersvc(conf, conn, usersvcClient)
```


### Client load balancing
Only a round robin load balancing strategy is implemented, welcome to submit pull request :)
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


### Demo

see[go-doudou-guide](https://github.com/unionj-cloud/go-doudou-guide) 


### kit

The kit package has some command line tools, which can be used after executing the installation command above 👆.

#### name

Generate the `json`tag behind the struct field according to the specified naming rules. [View document](./name/README.md)

#### ddl

A tool for synchronizing database table struct and Go struct based on [jmoiron/sqlx](https://github.com/jmoiron/sqlx). You can also generate DAO layer code.
[View document](./ddl/doc/README.md)

### Help
Welcome to mention pull request and issue, and welcome to scan the QR code and add the author's WeChat for comments and demands.Help me，hopefully work with me,to improve go-doudou.
![qrcode.png](qrcode.png)


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

2. 在svc.go文件里的interface里定义接口方法，在vo包里定义入参和出参结构体  
   此处略，见下文的[接口设计约束](#%E6%8E%A5%E5%8F%A3%E8%AE%BE%E8%AE%A1%E7%BA%A6%E6%9D%9F)和[vo包结构体设计约束](#vo%E5%8C%85%E7%BB%93%E6%9E%84%E4%BD%93%E8%AE%BE%E8%AE%A1%E7%BA%A6%E6%9D%9F)
   

3. 生成http接口代码
```shell
go-doudou svc http --handler -c go -o --doc
```
此时新增了一些文件夹
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
- auth_openapi3.json：openapi3.0规范的json格式接口文档
- client：包含golang的接口客户端代码，封装了[resty库](https://github.com/go-resty/resty)
- cmd：服务启动入口，需要在main方法里创建依赖的组件或者第三方服务客户端实例，注入本项目服务实例中
- config：配置文件相关
- db：生成数据库连接
- svcimpl.go：自定义服务的实现逻辑
- transport：包含生成的http routes和handlers
- .env：定义环境变量  


4. 将.env文件里的配置项GDD_SEED的值删掉，因为目前还没有种子  
   

5. 启动服务
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

从第6步开始是部署服务相关的步骤，需要本地有docker环境，连接到本地或者远程的k8s服务  


6. 打镜像
```shell
go-doudou svc push -r yourprivaterepositoryaddress
```  


7. 部署到k8s
```shell
go-doudou svc deploy 
```  


8. 关闭服务
```shell
go-doudou svc shutdown
```  


9. 伸缩服务
```shell
go-doudou svc scale -n 3
```


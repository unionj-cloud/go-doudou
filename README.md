## go-doudou
[![GoDoc](https://godoc.org/github.com/unionj-cloud/go-doudou?status.png)](https://godoc.org/github.com/unionj-cloud/go-doudou)
[![Build Status](https://travis-ci.com/unionj-cloud/go-doudou.svg?branch=main)](https://travis-ci.com/unionj-cloud/go-doudou)
[![codecov](https://codecov.io/gh/unionj-cloud/go-doudou/branch/main/graph/badge.svg?token=QRLPRAX885)](https://codecov.io/gh/unionj-cloud/go-doudou)
[![Go Report Card](https://goreportcard.com/badge/github.com/unionj-cloud/go-doudou)](https://goreportcard.com/report/github.com/unionj-cloud/go-doudou)

go-doudou（doudou发音/dəudəu/）是基于gossip协议做服务注册与发现，基于openapi 3.0规范做接口定义的go语言去中心化微服务敏捷开发框架。  
go-doudou通过一组命令行工具可以帮助开发者快速初始化一个或一组restful服务，通过在接口类中定义方法，即相当于设计了一组api，然后通过命令可以
生成启动服务的main方法，路由和相应的handler，以及go客户端代码。  
go-doudou主张设计优先，通过预先设计和定义接口，来生成代码，修改定义后，重新覆盖或者增量生成代码的方式来实现快速开发。  
go-doudou推崇契约精神，通过openapi 3.0协议来描述接口，规范服务提供方和消费方的合作，促使研发团队整体提高交付效率。
go-doudou致力于帮助开发者打造去中心化的微服务体系，通过gossip协议将集群内的服务连接起来，采用客户端负载均衡的方式调用其他服务，

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
### TOC

- [安装](#%E5%AE%89%E8%A3%85)
- [使用](#%E4%BD%BF%E7%94%A8)
- [工具箱](#%E5%B7%A5%E5%85%B7%E7%AE%B1)
  - [name](#name)
  - [ddl](#ddl)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->



### 安装

```shell
go get -v -u github.com/unionj-cloud/go-doudou/...@v0.3.3
```

### 使用
1. 以auth服务为例，初始化项目
```shell
go-doudou svc init auth
```
会生成如下项目结构
```shell
➜  auth git:(master) ✗ ll
total 24
-rw-r--r--  1 wubin1989  staff   372B  7  2 17:20 Dockerfile
-rw-r--r--  1 wubin1989  staff   399B  7  2 17:20 go.mod
-rw-r--r--  1 wubin1989  staff   241B  7  2 17:20 svc.go
drwxr-xr-x  3 wubin1989  staff    96B  7  2 17:20 vo
```
- Dockerfile：生成docker镜像
- svc.go：接口设计文件，里面是interface，在里面定义方法
- vo文件夹：里面定义struct，作为接口的入参和出参，也用于生成openapi3.0规范里的schema

2. 在svc.go文件里的interface里定义接口方法，在vo包里定义入参和出参结构体
   此处略，见下文详解  
   

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

4. 将.env文件里的配置项GDD_SEED的值生成空，因为目前还没有种子  
   

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

### 注意
暂时只支持http的restful接口，不支持grpc

### 工具箱

kit包有一些命令行工具，执行上面👆的安装命令后，就可以用了。

#### name

根据指定的命名规则生成结构体字段后面的`json`tag。[查看文档](./name/README.md)

#### ddl

基于[jmoiron/sqlx](https://github.com/jmoiron/sqlx) 实现的同步数据库表结构和Go结构体的工具。可以从结构体同步数据库表结构，也可以从数据库表结构生成结构体，还可以生成dao层代码。
[查看文档](./ddl/doc/README.md)









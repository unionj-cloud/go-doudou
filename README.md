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
- [注意](#%E6%B3%A8%E6%84%8F)
- [接口设计约束](#%E6%8E%A5%E5%8F%A3%E8%AE%BE%E8%AE%A1%E7%BA%A6%E6%9D%9F)
- [vo包结构体设计约束](#vo%E5%8C%85%E7%BB%93%E6%9E%84%E4%BD%93%E8%AE%BE%E8%AE%A1%E7%BA%A6%E6%9D%9F)
- [服务注册与发现](#%E6%9C%8D%E5%8A%A1%E6%B3%A8%E5%86%8C%E4%B8%8E%E5%8F%91%E7%8E%B0)
- [客户端负载均衡](#%E5%AE%A2%E6%88%B7%E7%AB%AF%E8%B4%9F%E8%BD%BD%E5%9D%87%E8%A1%A1)
- [Demo](#demo)
- [工具箱](#%E5%B7%A5%E5%85%B7%E7%AE%B1)
  - [name](#name)
  - [ddl](#ddl)
- [Help](#help)

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


### 注意

暂时只支持http的restful接口，不支持grpc和protobuffer


### 接口设计约束

1. 支持Post, Get, Delete, Put四种http请求方法，从接口方法名称来判断，默认是post请求，如果方法名以Post/Get/Delete/Put开头，
   则http请求方法分别为相对应的post/get/delete/put的其中一种  
2. 第一个入参的类型是context.Context，这个不要改，可以合理利用这个参数实现一些效果，比如当客户端取消请求，处理逻辑可以及时停止，节省服务器资源
3. 入参和出参的类型，仅支持go语言[内建类型](https://golang.org/pkg/builtin/) ，key为string类型的字典类型，vo包里自定义结构体以及上述类型相应的切片类型和指针类型。
   go-doudou生成代码和openapi文档的时候会扫描vo包里的结构体，如果接口的入参和出参里用了vo包以外的包里的结构体，go-doudou扫描不到结构体的字段。 
4. 特别的，入参还支持multipart.FileHeader类型，用于文件上传。出参还支持os.File类型，用于文件下载
5. 入参和出参的类型，不支持func类型，channel类型，接口类型和匿名结构体
6. 因为go的net/http包里的取Form参数相关的方法，比如FormValue，取到的参数值都是string类型的，go-doudou采用了cobra和viper的作者spf13大神的[cast](https://github.com/spf13/cast) 库做类型转换，
   生成的handlerimpl.go文件里的代码里解析表单参数的地方可能会报编译错误，可以给go-doudou提[issue](https://github.com/unionj-cloud/go-doudou/issues) ，也可以自己手动修改。
   当增删改了svc.go里的接口方法，重新执行代码生成命令`go-doudou svc http --handler -c go -o --doc`时，handlerimpl.go文件里的代码是增量生成的，
   即之前生成的代码和自己手动修改过的代码都不会被覆盖
7. handler.go文件里的代码在每次执行go-doudou svc http命令的时候都会重新生成，请不要手动修改里面的代码
8. 除handler.go和handlerimpl.go之外的其他文件，都是先判断是否存在，不存在才生成，存在就什么都不做


### vo包结构体设计约束

1. 结构体字段类型，仅支持go语言[内建类型](https://golang.org/pkg/builtin/) ，key为string类型的字典类型，vo包里自定义结构体，**匿名结构体**以及上述类型相应的切片类型和指针类型。
2. 结构体字段类型，不支持func类型，channel类型，接口类型
3. 结构体字段类型，不支持类型别名

### 服务注册与发现
go-doudou同时支持单体模式和微服务模式，以环境变量的方式配置。  
- `GDD_MODE=micro`：为微服务模式  
- `GDD_MODE=mono`：为单体模式  
在生成的cmd/main.go文件里有如下所示代码：  
```go
if ddconfig.GddMode.Load() == "micro" {
    node, err := registry.NewNode()
    if err != nil {
        logrus.Panicln(fmt.Sprintf("%+v", err))
    }
    logrus.Infof("Memberlist created. Local node is %s\n", node)
}
```
当只有其他服务依赖自己的时候，只需要把自己的服务通过`registry.NewNode()`方法注册上去即可。  
如果自己需要依赖其他服务，则除了需要把自己的服务注册到微服务集群之外，还需要加上实现服务发现的代码：
```go
// 注册自己并加入集群
node, err := registry.NewNode()
if err != nil {
    logrus.Panicln(fmt.Sprintf("%+v", err))
}
logrus.Infof("%s joined cluster\n", node.String())

// 需要依赖usersvc服务，那么就创建一个usersvc服务的provider
usersvcProvider := ddhttp.NewMemberlistServiceProvider("usersvc", node)
// 将usersvc服务的provider注入到usersvc服务的客户端实例里
usersvcClient := client.NewUsersvc(client.WithProvider(usersvcProvider))

// 将usersvc服务的客户端实例注入到自己的服务实例里
svc := service.NewOrdersvc(conf, conn, usersvcClient)
```


### 客户端负载均衡
暂时只实现了一种round robin的负载均衡策略，欢迎提pr:)
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

请参考[go-doudou-guide](https://github.com/unionj-cloud/go-doudou-guide) 


### 工具箱

kit包有一些命令行工具，执行上面👆的安装命令后，就可以用了。

#### name

根据指定的命名规则生成结构体字段后面的`json`tag。[查看文档](./name/README.md)

#### ddl

基于[jmoiron/sqlx](https://github.com/jmoiron/sqlx) 实现的同步数据库表结构和Go结构体的工具。还可以生成dao层代码。
[查看文档](./ddl/doc/README.md)


### Help
希望大家跟我一起完善go-doudou，欢迎提pr和issue，欢迎提意见和需求。欢迎扫描加作者微信交流。
![qrcode.png](qrcode.png)




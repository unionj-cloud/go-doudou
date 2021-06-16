## go-doudou
[![GoDoc](https://godoc.org/github.com/unionj-cloud/go-doudou?status.png)](https://godoc.org/github.com/unionj-cloud/go-doudou)
[![Build Status](https://travis-ci.com/unionj-cloud/go-doudou.svg?branch=main)](https://travis-ci.com/unionj-cloud/go-doudou)
[![codecov](https://codecov.io/gh/unionj-cloud/go-doudou/branch/main/graph/badge.svg?token=QRLPRAX885)](https://codecov.io/gh/unionj-cloud/go-doudou)
[![Go Report Card](https://goreportcard.com/badge/github.com/unionj-cloud/go-doudou)](https://goreportcard.com/report/github.com/unionj-cloud/go-doudou)

基于gossip协议的go语言微服务敏捷开发框架go-doudou（doudou发音/dəudəu/），包含必要的基础设施和提效的工具箱。还在开发中，敬请期待...

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
### TOC

- [安装](#%E5%AE%89%E8%A3%85)
- [工具箱](#%E5%B7%A5%E5%85%B7%E7%AE%B1)
  - [name](#name)
  - [ddl](#ddl)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->



### 安装

```shell
go get -v -u github.com/unionj-cloud/go-doudou/...@v0.1.7
```



### 工具箱

kit包有一些命令行工具，执行上面👆的安装命令后，就可以用了。

#### name

根据预设的命名规则生成结构体的Marshaler接口方法实现，省去了在结构体字段后面加`json`tag的工作。[查看文档](./name/README.md)

#### ddl

基于[jmoiron/sqlx](https://github.com/jmoiron/sqlx)实现的同步数据库表结构和Go结构体的工具。可以从结构体同步数据库表结构，也可以从数据库表结构生成结构体，还可以生成dao层代码。[查看文档](./ddl/doc/README.md)










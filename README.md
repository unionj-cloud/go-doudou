## go-doudou

基于gossip协议的go语言微服务敏捷开发框架go-doudou（doudou发音/dəudəu/），包含必要的基础设施和提效的工具箱。还在开发中，敬请期待...

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
### TOC

- [工具箱](#%E5%B7%A5%E5%85%B7%E7%AE%B1)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

### 工具箱

kit包有一些命令行工具：

- name: 根据预设的命名规则生成结构体的Marshaler接口方法实现，省去了在结构体字段后面加`json`tag的工作。[查看文档](./kit/name/README.md)
- ddl: 简单的orm工具。可以从结构体同步数据库表结构，也可以从数据库表结构生成结构体，还可以生成dao层代码。[查看文档](./kit/ddl/doc/README.md)
- ...








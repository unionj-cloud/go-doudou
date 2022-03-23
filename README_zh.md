<p align="center">
  <a href="https://go-doudou.github.io" target="_blank" rel="noopener noreferrer">
    <img width="180" src="https://go-doudou.github.io/hero.png" alt="Vite logo">
  </a>
</p>
<br/>
<p align="center">
  <a href="https://github.com/avelino/awesome-go"><img src="https://awesome.re/mentioned-badge.svg" alt="Mentioned in Awesome Go"></a>
  <a href="https://godoc.org/github.com/unionj-cloud/go-doudou"><img src="https://godoc.org/github.com/unionj-cloud/go-doudou?status.png" alt="GoDoc"></a>
  <a href="https://travis-ci.com/unionj-cloud/go-doudou"><img src="https://travis-ci.com/unionj-cloud/go-doudou.svg?branch=main" alt="Build Status"></a>
  <a href="https://github.com/unionj-cloud/go-doudou/actions/workflows/go.yml"><img src="https://github.com/unionj-cloud/go-doudou/actions/workflows/go.yml/badge.svg?branch=main" alt="Go"></a>
  <a href="https://codecov.io/gh/unionj-cloud/go-doudou"><img src="https://codecov.io/gh/unionj-cloud/go-doudou/branch/main/graph/badge.svg?token=QRLPRAX885" alt="codecov"></a>
  <a href="https://goreportcard.com/report/github.com/unionj-cloud/go-doudou"><img src="https://goreportcard.com/badge/github.com/unionj-cloud/go-doudou" alt="Go Report Card"></a>
  <a href="https://github.com/unionj-cloud/go-doudou"><img src="https://img.shields.io/github/v/release/unionj-cloud/go-doudou?style=flat-square" alt="Release"></a>
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License: MIT"></a>
  <a href="https://wakatime.com/badge/user/852bcf22-8a37-460a-a8e2-115833174eba/project/57c830f7-e507-4cb1-9fd1-feedd96685f6"><img src="https://wakatime.com/badge/user/852bcf22-8a37-460a-a8e2-115833174eba/project/57c830f7-e507-4cb1-9fd1-feedd96685f6.svg" alt="License: MIT"></a>
</p>
<br/>

# go-doudou

> 基于Gossip协议的去中心化Go语言微服务框架

- 💡 从Go语言接口类型开始，无须学习新的IDL语言（接口定义语言）
- 🛠️ 内建基于SWIM gossip协议的服务注册与发现机制，助你打造健壮、可弹性伸缩和去中心化的微服务集群
- 🔩 内建强大的代码生成器。在你定义完Go语言接口方法之后，唯一的工作只是实现你的独特创意
- ⚡ 生于云原生时代，内建命令行终端工具加速你的产品迭代
- 🔑 内建服务治理模块，支持远程配置管理、客户端负载均衡、熔断限流、隔仓、超时重试等等
- 📦️ 同时支持单体架构和微服务架构，可以自由设计你的系统

Go-doudou（doudou发音"兜兜"）是一个基于Gossip协议和OpenAPI v3接口描述规范的去中心化微服务框架。它同时支持开发单体架构的应用。目前，仅支持RESTful接口。

请阅读文档 [https://go-doudou.github.io/zh/](https://go-doudou.github.io/zh/) 了解更多。

## 感谢

Go-doudou是站在巨人的肩膀上开发而成的，在此感谢以下项目和它们的贡献者的无私付出：

- [hashicorp/memberlist](https://github.com/hashicorp/memberlist): go-doudou基于该库实现内建服务注册与发现和节点探活机制
- [gorilla/mux](https://github.com/gorilla/mux): go-doudou基于该库实现http路由
- [go-redis/redis_rate](github.com/go-redis/redis_rate): go-doudou基于该库实现基于Redis的跨节点限流机制
- [apolloconfig/agollo](https://github.com/apolloconfig/agollo): go-doudou基于该库实现了集成 [Apollo](https://github.com/apolloconfig/apollo) 的远程配置管理
- [nacos-group/nacos-sdk-go](https://github.com/nacos-group/nacos-sdk-go): go-doudou基于该库实现了集成 [Nacos](https://github.com/alibaba/nacos) 的服务注册与发现和远程配置管理

## 社区

欢迎加入go-doudou开发团队贡献代码。你可以fork本仓库并提交pr或者缺陷。如果你喜欢go-doudou，请给它一个Star！

你可以通过以下方式联系我

- 脸书: [https://www.facebook.com/bin.wu.94617999/](https://www.facebook.com/bin.wu.94617999/)
- 推特: [https://twitter.com/BINWU49205513](https://twitter.com/BINWU49205513)
- 邮箱: 328454505@qq.com
- 微信:  
  <img src="./qrcode.png" alt="wechat-group" width="240">
- 微信群:  
  <img src="./go-doudou-wechat-group.png" alt="wechat-group" width="240">
- QQ群:  
  <img src="./go-doudou-qq-group.png" alt="qq-group" width="240">

## 🔋 JetBrains开源授权

Go-doudou一直在JetBrains公司的免费开源授权下，通过GoLand IDE开发，在此表达我的感谢。

<a href="https://jb.gg/OpenSourceSupport" target="_blank"><img src="https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.png" alt="JetBrains Logo (Main) logo." width="300"></a>

## License

MIT

<p align="center">
  <a href="https://go-doudou.github.io" target="_blank" rel="noopener noreferrer">
    <img width="180" src="https://go-doudou.github.io/hero.png" alt="Vite logo">
  </a>
</p>
<br/>
<p align="center">
  <a href="https://github.com/avelino/awesome-go"><img src="https://awesome.re/mentioned-badge.svg" alt="Mentioned in Awesome Go"></a>
  <a href="https://godoc.org/github.com/unionj-cloud/go-doudou/v2"><img src="https://godoc.org/github.com/unionj-cloud/go-doudou/v2?status.png" alt="GoDoc"></a>
  <a href="https://github.com/unionj-cloud/go-doudou/v2/actions/workflows/go.yml"><img src="https://github.com/unionj-cloud/go-doudou/v2/actions/workflows/go.yml/badge.svg?branch=main" alt="Go"></a>
  <a href="https://codecov.io/gh/unionj-cloud/go-doudou"><img src="https://codecov.io/gh/unionj-cloud/go-doudou/branch/main/graph/badge.svg?token=QRLPRAX885" alt="codecov"></a>
  <a href="https://goreportcard.com/report/github.com/unionj-cloud/go-doudou/v2"><img src="https://goreportcard.com/badge/github.com/unionj-cloud/go-doudou/v2" alt="Go Report Card"></a>
  <a href="https://github.com/unionj-cloud/go-doudou/v2"><img src="https://img.shields.io/github/v/release/unionj-cloud/go-doudou?style=flat-square" alt="Release"></a>
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License: MIT"></a>
  <a href="https://wakatime.com/badge/user/852bcf22-8a37-460a-a8e2-115833174eba/project/57c830f7-e507-4cb1-9fd1-feedd96685f6"><img src="https://wakatime.com/badge/user/852bcf22-8a37-460a-a8e2-115833174eba/project/57c830f7-e507-4cb1-9fd1-feedd96685f6.svg" alt="License: MIT"></a>
</p>
<br/>

# go-doudou

> Lightweight Golang Microservice Framework

- 💡 Starts from golang interface, no need to learn new IDL(interface definition language).
- 🛠️ Built-in SWIM gossip protocol based service register and discovery mechanism to help you build a robust, scalable and decentralized service cluster.
- 🔩 Powerful code generator cli built-in. After defining your interface methods, your only job is implementing your awesome idea.
- ⚡ Born from the cloud-native era. Built-in CLI can speed up your product iteration.
- 🔑 Built-in service governance support including remote configuration management, client-side load balancer, rate limiter, circuit breaker, bulkhead, timeout, retry and more.
- 📦️ Supporting both monolith and microservice architectures gives you flexibility to design your system.

Go-doudou（doudou pronounce /dəudəu/）is OpenAPI 3.0 (for REST) spec and Protobuf v3 (for grpc) based lightweight microservice framework. It supports monolith service application as well.  

Read the Docs [https://go-doudou.github.io](https://go-doudou.github.io) to Learn More.

## Benchmark

![benchmark](./benchmark.png)

Machine: `MacBook Pro (16-inch, 2019)`  
CPU: `2.3 GHz 8 cores Intel Core i9`  
Memory: `16 GB 2667 MHz DDR4`
ProcessingTime: `0ms, 10ms, 100ms, 500ms`  
Concurrency: `1000`  
Duration: `30s`  
go-doudou Version: `v1.3.7`  

[Checkout the test code](https://github.com/wubin1989/go-web-framework-benchmark)

## Credits

Give credits to following repositories and all their contributors:
- [hashicorp/memberlist](https://github.com/hashicorp/memberlist): go-doudou is relying on it to implement service register/discovery/fault tolerance feature.
- [gorilla/mux](https://github.com/gorilla/mux): go-doudou is relying on it to implement http router.
- [go-redis/redis_rate](github.com/go-redis/redis_rate): go-doudou is relying on it to implement redis based rate limit feature
- [apolloconfig/agollo](https://github.com/apolloconfig/agollo): go-doudou is relying on it to implement remote configuration management support for [Apollo](https://github.com/apolloconfig/apollo)
- [nacos-group/nacos-sdk-go](https://github.com/nacos-group/nacos-sdk-go): go-doudou is relying on it to implement service discovery and remote configuration management support for [Nacos](https://github.com/alibaba/nacos)

## Community

Welcome to contribute to go-doudou by forking it and submitting pr or issues. If you like go-doudou, please give it a
star!

Welcome to contact me from

- Facebook: [https://www.facebook.com/bin.wu.94617999/](https://www.facebook.com/bin.wu.94617999/)
- Twitter: [https://twitter.com/BINWU49205513](https://twitter.com/BINWU49205513)
- Email: 328454505@qq.com
- WeChat:  
  <img src="./qrcode.png" alt="wechat-group" width="240">
- WeChat Group:  
  <img src="./go-doudou-wechat-group.png" alt="wechat-group" width="240">
- QQ group:  
  <img src="./go-doudou-qq-group.png" alt="qq-group" width="240">

## 🔋 JetBrains Open Source License

Go-doudou has been being developed with GoLand under the **free JetBrains Open Source license(s)** granted by JetBrains s.r.o., hence I would like to express my gratitude here.

<a href="https://jb.gg/OpenSourceSupport" target="_blank"><img src="https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.png" alt="JetBrains Logo (Main) logo." width="300"></a>

## License

MIT

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
  <a href="https://goproxy.cn"><img src="https://goproxy.cn/stats/github.com/unionj-cloud/go-doudou/badges/download-count.svg" alt="Goproxy.cn"></a>
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License: MIT"></a>
  <a href="https://wakatime.com/badge/user/852bcf22-8a37-460a-a8e2-115833174eba/project/57c830f7-e507-4cb1-9fd1-feedd96685f6"><img src="https://wakatime.com/badge/user/852bcf22-8a37-460a-a8e2-115833174eba/project/57c830f7-e507-4cb1-9fd1-feedd96685f6.svg" alt="License: MIT"></a>
</p>
<br/>

# go-doudou

> Gossip Protocol Distributed Golang Microservice Framework

- 💡 Starts from golang interface, no need to learn new IDL(interface definition language).
- 🛠️ Built-in SWIM gossip protocol based service register and discovery mechanism to help you build a robust, scalable and decentralized service cluster.
- 🔩 Powerful code generator cli built-in. After defining your interface methods, your only job is implementing your awesome idea.
- ⚡ Born from the cloud-native era. Built-in CLI can speed up your product iteration.
- 🔑 Built-in service governance support including client-side load balancer, rate limiter, circuit breaker, bulkhead, timeout, retry and more.
- 📦️ Supporting both monolith and microservice architectures gives you flexibility to design your system.

Go-doudou（doudou pronounce /dəudəu/）is a gossip protocol and OpenAPI 3.0 spec based decentralized microservice framework. It supports monolith service application as well. Currently, it supports RESTful service only.

[Read the Docs to Learn More](https://go-doudou.github.io).

## Credits

Give credits to following repositories and all their contributors:
- [https://github.com/hashicorp/memberlist](https://github.com/hashicorp/memberlist): go-doudou is relying on it to implement service register/discovery/fault tolerance feature.
- [github.com/go-redis/redis_rate](github.com/go-redis/redis_rate): go-doudou is relying on it to implement redis based rate limit feature

## Community

Welcome to contribute to go-doudou by forking it and submitting pr or issues. If you like go-doudou, please give it a
star!

Welcome to contact me from

- Facebook: [https://www.facebook.com/bin.wu.94617999/](https://www.facebook.com/bin.wu.94617999/)
- Twitter: [https://twitter.com/BINWU49205513](https://twitter.com/BINWU49205513)
- Email: 328454505@qq.com
- WeChat:  
  ![qrcode.png](qrcode.png)
- WeChat Group:  
  ![go-doudou-wechat-group.png](go-doudou-wechat-group.png)
- QQ group:  
  ![go-doudou-qq-group.png](go-doudou-qq-group.png)

## License

MIT

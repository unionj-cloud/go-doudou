# GRPC Resolver Nacos 

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/hang666/grpc_resolver_nacos)
![GitHub](https://img.shields.io/github/license/hang666/grpc_resolver_nacos)

本项目实现了Grpc中Nacos的服务发现

- 支持负载均衡中Nacos权重选择器
- 支持 grpc-gateway
- 使用 Nacos-sdk-go 同结构配置方便



## Installation

```
$ go get -u github.com/hang666/grpc_resolver_nacos
```

## Client Example

```go
// 创建NacosClientConfig
clientConfig := *constant.NewClientConfig(
	constant.WithNamespaceId(""),
	constant.WithTimeoutMs(5000),
	constant.WithNotLoadCacheAtStart(true),
	constant.WithLogDir("/tmp/nacos/log"),
	constant.WithCacheDir("/tmp/nacos/cache"),
	constant.WithLogLevel("debug"),
)
// 创建NacosServerConfigs (兼容多发现中心)
serverConfigs := []constant.ServerConfig{*constant.NewServerConfig("127.0.0.1", 8848, constant.WithContextPath("/nacos"))}
// 创建NacosNamingClient
client, err := clients.CreateNamingClient(map[string]interface{}{
	"serverConfigs": serverConfigs,
	"clientConfig":  &clientConfig,
})
// 添加Nacos配置 (支持多服务)
grpc_resolver_nacos.AddNacosConfig(grpc_resolver_nacos.NacosConfig{
	Label:              "user",			//Label与ServiceName一致即可
	ServiceName:        "user",			//Nacos内注册的服务名
	Clusters           	[]string{},
	GroupName          	"",
	NacosClientConfig:  clientConfig,
	NacosServerConfigs: serverConfigs,
	NacosClient:        client,
})
// target按照 nacos://ServiceName/ 填写即可，如上添加过的配置
// grpc-gateway中RegisterXXXHandlerFromEndpoint 如此target填写相同即可
conn, err := grpc.Dial("nacos://user/",
	grpc.WithTransportCredentials(insecure.NewCredentials()),
	grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "nacos_weight_balancer"}`))
```

## Load Balancing
```go
// 如下将 {"loadBalancingPolicy": "nacos_weight_balancer"} 添加DialOption即可
conn, err := grpc.Dial("nacos://user/",
	grpc.WithTransportCredentials(insecure.NewCredentials()),
	grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "nacos_weight_balancer"}`))
```

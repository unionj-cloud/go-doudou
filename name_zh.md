## name

根据指定的命名规则生成结构体字段后面的`json`tag。默认生成策略是**首字母小写的驼峰命名策略**，同时支持蛇形命名。  
未导出的字段会跳过，只修改导出字段的json标签。


<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
### TOC

- [命令行参数](#%E5%91%BD%E4%BB%A4%E8%A1%8C%E5%8F%82%E6%95%B0)
- [用法](#%E7%94%A8%E6%B3%95)
- [TODO](#todo)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->



### 命令行参数

```shell
➜  go-doudou git:(main) ✗ go-doudou name -h    
A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.

Usage:
  go-doudou name [flags]

Flags:
  -f, --file string       absolute path of vo file
  -h, --help              help for name
  -o, --omitempty         whether omit empty value or not
  -s, --strategy string   name of strategy, currently only support "lowerCamel" and "snake" (default "lowerCamel")

Global Flags:
      --config string   config file (default is $HOME/.go-doudou.yaml)
```



### 用法

- 在go文件里写上`//go:generate go-doudou name --file $GOFILE`，不限位置，最好是写在上方。目前的实现是对整个文件的所有struct都生效。

```go
//go:generate go-doudou name --file $GOFILE

type Event struct {
	Name      string
	EventType int
}

type TestName struct {
	Age    age
	School []struct {
		Name string
		Addr struct {
			Zip   string
			Block string
			Full  string
		}
	}
	EventChan chan Event
	SigChan   chan int
	Callback  func(string) bool
	CallbackN func(param string) bool
}
```

- 在项目根路径下执行命令`go generate ./...`

```go
type Event struct {
	Name      string `json:"name"`
	EventType int    `json:"eventType"`
}

type TestName struct {
	Age    age `json:"age"`
	School []struct {
		Name string `json:"name"`
		Addr struct {
			Zip   string `json:"zip"`
			Block string `json:"block"`
			Full  string `json:"full"`
		} `json:"addr"`
	} `json:"school"`
	EventChan chan Event              `json:"eventChan"`
	SigChan   chan int                `json:"sigChan"`
	Callback  func(string) bool       `json:"callback"`
	CallbackN func(param string) bool `json:"callbackN"`
}
```


### TODO

+ [x] 支持omitempty
+ [x] 蛇形命名策略






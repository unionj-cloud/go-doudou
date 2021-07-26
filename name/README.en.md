## name

Generate the `json`tag behind the struct field according to the specified naming rules. The default strategy is **camel case with lowercase first letter**, and also supports snake case.  


### Command line flags

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



### Usage

- Write `//go:generate go-doudou name --file $GOFILE` in the go file. There is no limit to the position, it is best to write it at the top. The current implementation is effective for all structs of the entire file.

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

- Execute the command `go generate ./...` under the project's root path

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

+ [x] Support omitempty
+ [x] Snake case
+ [ ] Only valid for the struct marked with `//go:generate name -file $GOFILE` above, not for the entire file's structs






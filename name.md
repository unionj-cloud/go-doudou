## name

Command line tool for generating json tag of struct field. Default strategy is camel case with first letter lowercased. Support snake case as well.  
Unexported fields will be skipped, only modify json tag of each exported field.


<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
### TOC

- [Flags](#flags)
- [Usage](#usage)
- [TODO](#todo)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->



### Flags

```shell
➜  go-doudou git:(main) go-doudou name -h   
WARN[0000] Error loading .env file: open /Users/wubin1989/workspace/cloud/.env: no such file or directory 
bulk add or update struct fields json tag

Usage:
  go-doudou name [flags]

Flags:
  -f, --file string       absolute path of vo file
  -h, --help              help for name
  -o, --omitempty         whether omit empty value or not
  -s, --strategy string   name of strategy, currently only support "lowerCamel" and "snake" (default "lowerCamel")
```



### Usage

- Put `//go:generate go-doudou name --file $GOFILE` into go file

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

- Execute  `go generate ./...` at the same folder

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
+ [x] Snake case strategy






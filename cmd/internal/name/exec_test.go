package name

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/pathutils"
)

const initCode = `package testdata

// 筛选条件
type PageFilter struct {
	// 真实姓名，前缀匹配
	Name string
	// 所属部门ID
	Dept int
}

//排序条件
type Order struct {
	Col  string	` + "`" + `json:"-"` + "`" + `
	sort string
}

type PageRet struct {
	Items    interface{}
	PageNo   int
	PageSize int
	Total    int
	HasNext  bool
}

type Field struct {
	Name   string
	Type   string
	Format string
}

type Base struct {
	Index string
	Type  string
}

type MappingPayload struct {
	Base
	Fields []Field
	Index  string
}
`

func TestName_Exec(t *testing.T) {
	type fields struct {
		File      string
		Strategy  string
		Omitempty bool
	}
	tests := []struct {
		name     string
		fields   fields
		want     string
		initCode string
	}{
		{
			name: "",
			fields: fields{
				File:     pathutils.Abs("testdata/vo.go"),
				Strategy: lowerCamelStrategy,
			},
			initCode: initCode,
			want: `package testdata

// 筛选条件
type PageFilter struct {
	// 真实姓名，前缀匹配
	Name string ` + "`" + `json:"name"` + "`" + `
	// 所属部门ID
	Dept int ` + "`" + `json:"dept"` + "`" + `
}

//排序条件
type Order struct {
	Col  string ` + "`" + `json:"-"` + "`" + `
	sort string
}

type PageRet struct {
	Items    interface{} ` + "`" + `json:"items"` + "`" + `
	PageNo   int         ` + "`" + `json:"pageNo"` + "`" + `
	PageSize int         ` + "`" + `json:"pageSize"` + "`" + `
	Total    int         ` + "`" + `json:"total"` + "`" + `
	HasNext  bool        ` + "`" + `json:"hasNext"` + "`" + `
}

type Field struct {
	Name   string ` + "`" + `json:"name"` + "`" + `
	Type   string ` + "`" + `json:"type"` + "`" + `
	Format string ` + "`" + `json:"format"` + "`" + `
}

type Base struct {
	Index string ` + "`" + `json:"index"` + "`" + `
	Type  string ` + "`" + `json:"type"` + "`" + `
}

type MappingPayload struct {
	Base
	Fields []Field ` + "`" + `json:"fields"` + "`" + `
	Index  string  ` + "`" + `json:"index"` + "`" + `
}
`,
		},
		{
			name: "",
			fields: fields{
				File:      pathutils.Abs("testdata/vo.go"),
				Strategy:  lowerCamelStrategy,
				Omitempty: true,
			},
			initCode: initCode,
			want: `package testdata

// 筛选条件
type PageFilter struct {
	// 真实姓名，前缀匹配
	Name string ` + "`" + `json:"name,omitempty"` + "`" + `
	// 所属部门ID
	Dept int ` + "`" + `json:"dept,omitempty"` + "`" + `
}

//排序条件
type Order struct {
	Col  string ` + "`" + `json:"-"` + "`" + `
	sort string
}

type PageRet struct {
	Items    interface{} ` + "`" + `json:"items,omitempty"` + "`" + `
	PageNo   int         ` + "`" + `json:"pageNo,omitempty"` + "`" + `
	PageSize int         ` + "`" + `json:"pageSize,omitempty"` + "`" + `
	Total    int         ` + "`" + `json:"total,omitempty"` + "`" + `
	HasNext  bool        ` + "`" + `json:"hasNext,omitempty"` + "`" + `
}

type Field struct {
	Name   string ` + "`" + `json:"name,omitempty"` + "`" + `
	Type   string ` + "`" + `json:"type,omitempty"` + "`" + `
	Format string ` + "`" + `json:"format,omitempty"` + "`" + `
}

type Base struct {
	Index string ` + "`" + `json:"index,omitempty"` + "`" + `
	Type  string ` + "`" + `json:"type,omitempty"` + "`" + `
}

type MappingPayload struct {
	Base
	Fields []Field ` + "`" + `json:"fields,omitempty"` + "`" + `
	Index  string  ` + "`" + `json:"index,omitempty"` + "`" + `
}
`,
		},
		{
			name: "",
			fields: fields{
				File:      pathutils.Abs("testdata/vo.go"),
				Strategy:  snakeStrategy,
				Omitempty: false,
			},
			initCode: initCode,
			want: `package testdata

// 筛选条件
type PageFilter struct {
	// 真实姓名，前缀匹配
	Name string ` + "`" + `json:"name"` + "`" + `
	// 所属部门ID
	Dept int ` + "`" + `json:"dept"` + "`" + `
}

//排序条件
type Order struct {
	Col  string ` + "`" + `json:"-"` + "`" + `
	sort string
}

type PageRet struct {
	Items    interface{} ` + "`" + `json:"items"` + "`" + `
	PageNo   int         ` + "`" + `json:"page_no"` + "`" + `
	PageSize int         ` + "`" + `json:"page_size"` + "`" + `
	Total    int         ` + "`" + `json:"total"` + "`" + `
	HasNext  bool        ` + "`" + `json:"has_next"` + "`" + `
}

type Field struct {
	Name   string ` + "`" + `json:"name"` + "`" + `
	Type   string ` + "`" + `json:"type"` + "`" + `
	Format string ` + "`" + `json:"format"` + "`" + `
}

type Base struct {
	Index string ` + "`" + `json:"index"` + "`" + `
	Type  string ` + "`" + `json:"type"` + "`" + `
}

type MappingPayload struct {
	Base
	Fields []Field ` + "`" + `json:"fields"` + "`" + `
	Index  string  ` + "`" + `json:"index"` + "`" + `
}
`,
		},
		{
			name: "",
			fields: fields{
				File:      pathutils.Abs("testdata/vo1.go"),
				Strategy:  snakeStrategy,
				Omitempty: false,
			},
			initCode: `package testdata

import "time"

// comment for alia age
type age int

type Event struct {
	Name      string
	EventType int
}

type TestAlias struct {
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

type ta TestAlias

type tt time.Time

type mm map[string]interface{}

type MyInter interface {
	Speak() error
}

type starM *time.Time
`,
			want: `package testdata

import "time"

// comment for alia age
type age int

type Event struct {
	Name      string ` + "`" + `json:"name"` + "`" + `
	EventType int    ` + "`" + `json:"event_type"` + "`" + `
}

type TestAlias struct {
	Age    age ` + "`" + `json:"age"` + "`" + `
	School []struct {
		Name string ` + "`" + `json:"name"` + "`" + `
		Addr struct {
			Zip   string ` + "`" + `json:"zip"` + "`" + `
			Block string ` + "`" + `json:"block"` + "`" + `
			Full  string ` + "`" + `json:"full"` + "`" + `
		} ` + "`" + `json:"addr"` + "`" + `
	} ` + "`" + `json:"school"` + "`" + `
	EventChan chan Event              ` + "`" + `json:"event_chan"` + "`" + `
	SigChan   chan int                ` + "`" + `json:"sig_chan"` + "`" + `
	Callback  func(string) bool       ` + "`" + `json:"callback"` + "`" + `
	CallbackN func(param string) bool ` + "`" + `json:"callback_n"` + "`" + `
}

type ta TestAlias

type tt time.Time

type mm map[string]interface{}

type MyInter interface {
	Speak() error
}

type starM *time.Time
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := Name{
				File:      tt.fields.File,
				Strategy:  tt.fields.Strategy,
				Omitempty: tt.fields.Omitempty,
			}
			receiver.Exec()
			f, err := os.Open(tt.fields.File)
			if err != nil {
				t.Fatal(err)
			}
			content, err := ioutil.ReadAll(f)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.want, string(content))
			ioutil.WriteFile(tt.fields.File, []byte(tt.initCode), os.ModePerm)
		})
	}
}

func TestPanic(t *testing.T) {
	receiver := Name{
		File:      pathutils.Abs("testdata/vo1.go"),
		Strategy:  "unknownstrategy",
		Omitempty: false,
	}
	assert.Panics(t, func() {
		receiver.Exec()
	})
}

func TestPanic2(t *testing.T) {
	receiver := Name{
		File:      pathutils.Abs("testdata/vo2"),
		Strategy:  lowerCamelStrategy,
		Omitempty: false,
	}
	assert.Panics(t, func() {
		receiver.Exec()
	})
}

func ExampleNameForm_Exec() {
	receiver := Name{
		File:      pathutils.Abs("testdata/vo.go"),
		Strategy:  lowerCamelStrategy,
		Omitempty: false,
		Form:      true,
	}
	receiver.Exec()
	// output:
}

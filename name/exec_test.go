package name

import (
	"github.com/unionj-cloud/go-doudou/pathutils"
	"io/ioutil"
	"os"
	"testing"
)

const initCode = `package testfiles

// 筛选条件
type PageFilter struct {
	// 真实姓名，前缀匹配
	Name string
	// 所属部门ID
	Dept int
}

//排序条件
type Order struct {
	Col  string
	Sort string
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
		name   string
		fields fields
		want   string
	}{
		{
			name: "",
			fields: fields{
				File:     pathutils.Abs("testfiles/vo.go"),
				Strategy: lowerCamelStrategy,
			},
			want: `package testfiles

// 筛选条件
type PageFilter struct {
	// 真实姓名，前缀匹配
	Name string ` + "`" + `json:"name"` + "`" + `
	// 所属部门ID
	Dept int ` + "`" + `json:"dept"` + "`" + `
}

//排序条件
type Order struct {
	Col  string ` + "`" + `json:"col"` + "`" + `
	Sort string ` + "`" + `json:"sort"` + "`" + `
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
				File:      pathutils.Abs("testfiles/vo.go"),
				Strategy:  lowerCamelStrategy,
				Omitempty: true,
			},
			want: `package testfiles

// 筛选条件
type PageFilter struct {
	// 真实姓名，前缀匹配
	Name string ` + "`" + `json:"name,omitempty"` + "`" + `
	// 所属部门ID
	Dept int ` + "`" + `json:"dept,omitempty"` + "`" + `
}

//排序条件
type Order struct {
	Col  string ` + "`" + `json:"col,omitempty"` + "`" + `
	Sort string ` + "`" + `json:"sort,omitempty"` + "`" + `
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
				File:      pathutils.Abs("testfiles/vo.go"),
				Strategy:  snakeStrategy,
				Omitempty: false,
			},
			want: `package testfiles

// 筛选条件
type PageFilter struct {
	// 真实姓名，前缀匹配
	Name string ` + "`" + `json:"name"` + "`" + `
	// 所属部门ID
	Dept int ` + "`" + `json:"dept"` + "`" + `
}

//排序条件
type Order struct {
	Col  string ` + "`" + `json:"col"` + "`" + `
	Sort string ` + "`" + `json:"sort"` + "`" + `
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
			if string(content) != tt.want {
				t.Errorf("want %s, got %s\n", tt.want, string(content))
			}
			ioutil.WriteFile(tt.fields.File, []byte(initCode), os.ModePerm)
		})
	}
}

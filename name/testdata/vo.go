package testdata

// 筛选条件
type PageFilter struct {
	// 真实姓名，前缀匹配
	Name string
	// 所属部门ID
	Dept int
}

//排序条件
type Order struct {
	Col  string	`json:"-"`
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

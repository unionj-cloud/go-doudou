package main

//go:generate go-doudou name --file $GOFILE -o

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

// 分页筛选条件
type PageQuery struct {
	Filter PageFilter
	Page   Page
}

type PageRet struct {
	Items    interface{}
	PageNo   int
	PageSize int
	Total    int
	HasNext  bool
}

type queryType int
type queryLogic int

const (
	SHOULD queryLogic = iota + 1
	MUST
	MUSTNOT
)

const (
	TERMS queryType = iota + 1
	MATCHPHRASE
	RANGE
	PREFIX
	WILDCARD
	EXISTS
)

type esFieldType string

const (
	TEXT    esFieldType = "text"
	KEYWORD esFieldType = "keyword"
	DATE    esFieldType = "date"
	LONG    esFieldType = "long"
	INTEGER esFieldType = "integer"
	SHORT   esFieldType = "short"
	DOUBLE  esFieldType = "double"
	FLOAT   esFieldType = "float"
	BOOL    esFieldType = "boolean"
)

type Base struct {
	Index string
	Type  string
}

type QueryCond struct {
	Pair       map[string][]interface{}
	QueryLogic queryLogic
	QueryType  queryType
	Children   []QueryCond
}

type Sort struct {
	Field     string
	Ascending bool
}

type Paging struct {
	StartDate  string
	EndDate    string
	DateField  string
	QueryConds []QueryCond
	Skip       int
	Limit      int
	Sortby     []Sort
}

type BulkSavePayload struct {
	Base
	Docs []map[string]interface{}
}

type SavePayload struct {
	Base
	Doc map[string]interface{}
}

type BulkDeletePayload struct {
	Base
	DocIds []string
}

type PagePayload struct {
	Base
	Paging
}

type PageResult struct {
	Page        int
	PageSize    int
	Total       int
	Docs        []map[string]interface{}
	HasNextPage bool
}

type StatPayload struct {
	Base
	Paging
	Aggr interface{}
}

type RandomPayload struct {
	Base
	Paging
}

type CountPayload struct {
	Base
	Paging
}

type Field struct {
	Name   string
	Type   esFieldType
	Format string
}

type MappingPayload struct {
	Base
	Fields []Field
}

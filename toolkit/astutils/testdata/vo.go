package main

//go:generate go-doudou name --file $GOFILE -o

// 筛选条件
type PageFilter struct {
	// 真实姓名，前缀匹配
	Name string `json:"name,omitempty"`
	// 所属部门ID
	Dept int `json:"dept,omitempty"`
}

// 排序条件
type Order struct {
	Col                string `json:"col,omitempty"`
	Sort, Name, Banana string
}

// 分页筛选条件
type PageQuery struct {
	Filter PageFilter `json:"filter,omitempty"`
	Page   Page       `json:"page,omitempty"`
}

type PageRet struct {
	Items    interface{} `json:"items,omitempty"`
	PageNo   int         `json:"pageNo,omitempty"`
	PageSize int         `json:"pageSize,omitempty"`
	Total    int         `json:"total,omitempty"`
	HasNext  bool        `json:"hasNext,omitempty"`
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
	Index string `json:"index,omitempty"`
	Type  string `json:"type,omitempty"`
}

type QueryCond struct {
	Pair       map[string][]interface{} `json:"pair,omitempty"`
	QueryLogic queryLogic               `json:"queryLogic,omitempty"`
	QueryType  queryType                `json:"queryType,omitempty"`
	Children   []QueryCond              `json:"children,omitempty"`
}

type Sort struct {
	Field     string `json:"field,omitempty"`
	Ascending bool   `json:"ascending,omitempty"`
}

type Paging struct {
	StartDate  string      `json:"startDate,omitempty"`
	EndDate    string      `json:"endDate,omitempty"`
	DateField  string      `json:"dateField,omitempty"`
	QueryConds []QueryCond `json:"queryConds,omitempty"`
	Skip       int         `json:"skip,omitempty"`
	Limit      int         `json:"limit,omitempty"`
	Sortby     []Sort      `json:"sortby,omitempty"`
}

type BulkSavePayload struct {
	Base
	Docs []map[string]interface{} `json:"docs,omitempty"`
}

type SavePayload struct {
	Base
	Doc map[string]interface{} `json:"doc,omitempty"`
}

type BulkDeletePayload struct {
	Base
	DocIds []string `json:"docIds,omitempty"`
}

type PagePayload struct {
	Base
	Paging
}

type PageResult struct {
	Page        int                      `json:"page,omitempty"`
	PageSize    int                      `json:"pageSize,omitempty"`
	Total       int                      `json:"total,omitempty"`
	Docs        []map[string]interface{} `json:"docs,omitempty"`
	HasNextPage bool                     `json:"hasNextPage,omitempty"`
}

type StatPayload struct {
	Base
	Paging
	Aggr interface{} `json:"aggr,omitempty"`
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
	Name   string      `json:"name,omitempty"`
	Type   esFieldType `json:"type,omitempty"`
	Format string      `json:"format,omitempty"`
}

type MappingPayload struct {
	Base
	Fields []Field `json:"fields,omitempty"`
}

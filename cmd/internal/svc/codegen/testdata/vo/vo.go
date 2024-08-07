package vo

import "github.com/shopspring/decimal"

//go:generate go-doudou name --file $GOFILE -o

// 筛选条件
type PageFilter struct {
	// 真实姓名，前缀匹配
	Name string
	// 所属部门ID
	Dept int
}

// 排序条件
type Order struct {
	Col  string
	Sort string
}

type Page struct {
	// 排序规则
	Orders []Order
	// 页码
	PageNo int
	// 每页行数
	Size int
	User UserVo
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
	Price    decimal.Decimal
}

type UserVo struct {
	Id    int
	Name  string
	Phone string
	Dept  string
}

package vo

import "github.com/unionj-cloud/go-doudou/ddl/query"

//go:generate go-doudou name --file $GOFILE

type Ret struct {
	Code int
	Data interface{}
	Msg  string
}

type PageFilter struct {
	// 真实姓名，前缀匹配
	Name string
	// 所属部门ID
	Dept int
}

// 分页筛选条件
type PageQuery struct {
	filter PageFilter
	page   query.Page
}

type UserVo struct {
	Id    int
	Name  string
	Phone string
	Dept  string
}

package vo

//go:generate go-doudou name --file $GOFILE --form

type PageFilter struct {
	// 真实姓名，前缀匹配
	Name string `json:"name" form:"name"`
	// 所属部门
	Dept string `json:"dept" form:"dept"`
}

type Order struct {
	Col  string `json:"col" form:"col"`
	Sort string `json:"sort" form:"sort"`
}

type Page struct {
	// 排序规则
	Orders []Order `json:"orders" form:"orders"`
	// 页码
	PageNo int `json:"pageNo" form:"pageNo"`
	// 每页行数
	Size int `json:"size" form:"size"`
}

// 分页筛选条件
type PageQuery struct {
	Filter PageFilter `json:"filter" form:"filter"`
	Page   Page       `json:"page" form:"page"`
}

type PageRet struct {
	Items    interface{} `json:"items" form:"items"`
	PageNo   int         `json:"pageNo" form:"pageNo"`
	PageSize int         `json:"pageSize" form:"pageSize"`
	Total    int         `json:"total" form:"total"`
	HasNext  bool        `json:"hasNext" form:"hasNext"`
}

type UserVo struct {
	Id       int    `json:"id" form:"id"`
	Username string `json:"username" form:"username"`
	Name     string `json:"name" form:"name"`
	Phone    string `json:"phone" form:"phone"`
	Dept     string `json:"dept" form:"dept"`
}

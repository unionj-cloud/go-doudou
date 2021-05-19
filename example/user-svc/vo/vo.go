package vo

//go:generate go-doudou name --file $GOFILE

type Ret struct {
	Code int
	Data interface{}
	Msg  string
	Err  error
	Mmm  map[string]UserVo
}

// 用户注册表单
type SignUpForm struct {
	// 用户名
	Username string
	// 密码
	Password string
	// 再次输入密码
	PassConfirm string
}

// 用户登录表单
type LogInForm struct {
	// 用户名
	Username string
	// 密码
	Password string
}

type PageFilter struct {
	// 真实姓名，前缀匹配
	Name string
	// 所属部门ID
	Dept int
}

const (
	Asc  = "asc"
	Desc = "desc"
)

type Order struct {
	Col  string
	Sort string
}

type Page struct {
	Orders [3]Order
	Offset int
	Size   int
}

type PageRet struct {
	Items    interface{}
	PageNo   int
	PageSize int
	Total    int
	HasNext  bool
}

// 分页筛选条件
type PageQuery struct {
	Filter PageFilter
	Page   Page
}

type UserVo struct {
	Id    int
	Name  string
	Phone string
	Dept  string
}

type Auth struct {
	Token string
	User  UserVo
}

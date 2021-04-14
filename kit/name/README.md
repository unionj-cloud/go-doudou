## name

根据预设的命名规则生成结构体的Marshaler接口方法实现，省去了在结构体字段后面加`json`tag的工作。默认生成策略是**首字母小写的驼峰命名策略**。暂时只支持这一种策略。欢迎提pr。



### 命令行参数

```shell
➜  example git:(main) ✗ name -h                                                         
Usage of name:
  -file string
    	absolute path of vo file
  -strategy string
    	name of strategy (default "lowerCaseNamingStrategy")
```



### 用法

- 在go文件里写上`//go:generate name -file $GOFILE`，不限位置，最好是写在上方。目前的实现是对整个文件的所有struct都生效。

```go
//go:generate name -file $GOFILE
type Student struct {
	School
	Company

	Name string
	Age  int

	TestScore int

	IsPaid bool
}

type School struct {
	Name     string
	Address  string
	CreateAt time.Time
}

type Company struct {
	Name     string
	Biz      string
	CreateAt time.Time
}
```

- 在项目根路径下执行命令`go generate ./...`

```go
func (object Student) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	objectMap[strategies.LowerCaseConvert("School")] = object.School
	objectMap[strategies.LowerCaseConvert("Company")] = object.Company
	objectMap[strategies.LowerCaseConvert("Name")] = object.Name
	objectMap[strategies.LowerCaseConvert("Age")] = object.Age
	objectMap[strategies.LowerCaseConvert("TestScore")] = object.TestScore
	objectMap[strategies.LowerCaseConvert("IsPaid")] = object.IsPaid
	return json.Marshal(objectMap)
}

func (object School) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	objectMap[strategies.LowerCaseConvert("Name")] = object.Name
	objectMap[strategies.LowerCaseConvert("Address")] = object.Address
	objectMap[strategies.LowerCaseConvert("CreateAt")] = object.CreateAt
	return json.Marshal(objectMap)
}

func (object Company) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	objectMap[strategies.LowerCaseConvert("Name")] = object.Name
	objectMap[strategies.LowerCaseConvert("Biz")] = object.Biz
	objectMap[strategies.LowerCaseConvert("CreateAt")] = object.CreateAt
	return json.Marshal(objectMap)
}
```

- 使用

  ```go
  func main() {
  	sch := vo.School{
  		Name:     "Beijing University",
  		Address:  "Beijing",
  		CreateAt: time.Now(),
  	}
  	company := vo.Company{
  		Name:     "Stq",
  		CreateAt: time.Now().AddDate(0, -1, 0),
  	}
  	stu := vo.Student{
  		Name:      "wubin",
  		Age:       30,
  		TestScore: 99,
  		School:    sch,
  		Company:   company,
  	}
  
  	data, err := json.Marshal(stu)
  	if err != nil {
  		panic(err)
  	}
  	fmt.Println(string(data))
  }
  ```

  输出

  ```json
  {"age":30,"company":{"biz":"","createAt":"2021-03-14T11:09:29.606039+08:00","name":"Stq"},"isPaid":false,"name":"wubin","school":{"address":"Beijing","createAt":"2021-04-14T11:09:29.606039+08:00","name":"Beijing University"},"testScore":99}
  ```



### TODO

+ [ ] 蛇形命名策略
+ [ ] 只针对上方标注了`//go:generate name -file $GOFILE`的结构体生效，而不是对整个文件的结构体都生效






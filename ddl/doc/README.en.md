## ddl

A tool for synchronizing database table struct and Go struct based on [jmoiron/sqlx](https://github.com/jmoiron/sqlx). Does not support index update, and does not support foreign keys.
**Add transaction support for dao layer code**, outsourcing a layer of abstraction in sqlx.Tx and sqlx.DB.


### Quick start

- ```
  git clone git@github.com:unionj-cloud/go-doudou.git
  ```

- ```
  cd go-doudou/example
  ```

- Current example directory structure

  ![example1](./example1.jpeg)

- ```
  go-doudou ddl --dao --pre=biz_ --domain=ddl/domain --env=ddl/.env
  ```

- Example directory structure after executing the above command

  ![example2](./example2.jpeg)

  Generated table structure

  ![table](./table.jpeg)

- ```
   go run ddl/main.go
  ```

  You can see the command line output

  ```
  ➜  example git:(main) ✗ go run ddl/main.go
  INFO[2021-04-13 17:52:42] {Items:[{ID:5 Name:Biden Phone:13893997878 Age:70 No:46 School:Harvard Univ. IsStudent:true Base:{CreateAt:2021-04-13 17:25:46 +0800 CST UpdateAt:2021-04-13 17:25:46 +0800 CST DeleteAt:<nil>}}] PageNo:1 PageSize:1 Total:1 HasNext:false} 
  ```



Insert two records into the database table

![data](./data.jpeg)



### Command line flags

```
➜  ~ go-doudou ddl -h
A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.

Usage:
  go-doudou ddl [flags]

Flags:
  -d, --dao             If true, generate dao code.
      --df string       Name of dao folder. (default "dao")
      --domain string   Path of domain folder. (default "domain")
      --env string      Path of database connection config .env file (default ".env")
  -h, --help            help for ddl
      --pre string      Table name prefix. e.g.: prefix biz_ for biz_product.
  -r, --reverse         If true, generate domain code from database. If false, update or create database tables from domain code.

Global Flags:
      --config string   config file (default is $HOME/.go-doudou.yaml)
```



### API

#### Example

```
//dd:table
type User struct {
	ID        int    `dd:"pk;auto"`
	Name      string `dd:"index:name_phone_idx,2;default:'jack'"`
	Phone     string `dd:"index:name_phone_idx,1;default:'13552053960';extra:comment '手机号'"`
	Age       int    `dd:"index"`
	No        int    `dd:"unique"`
	School    string `dd:"null;default:'harvard';extra:comment '学校'"`
	IsStudent bool

	Base
}
```

- A comment "//dd:table" should be added above the struct definition
- The struct field label named "dd"



#### Struct tag

##### pk

primary key

##### auto

auto increment

##### type

Represents the database field type. Not required. The default mapping rules are as follows:

| Supported Go language types (including pointers) | Database field types |
| :------------------------: | :------------: |
|            Int             |      Int       |
|           Int64            |     bigint     |
|          float32           |     Float      |
|          float64           |     Double     |
|           string           |  varchar(255)  |
|            bool            |    tinyint     |
|         time.Time          |    datetime    |

##### default

Represents the default value. You need to distinguish between built-in database functions or literals. If it is a database built-in function or an expression composed of a database built-in function, single quotes are not required. If it is literal, you need to add single quotes yourself.

##### extra

Represents other field information, such as "on update CURRENT_TIMESTAMP"，"comment '手机号'"
**Note: Because the ddl tool uses English ";" and ":" to parse tags, please do not use English ";" and ":" in the comment.**

##### index

index

- Format: "index: index name, sort, ascending order" or "index"
- Index name: string. If multiple fields use the same index name, it means that the index is a union index. Not required. If the index name is not declared, the default name is: column name+_idx
- Sort: integer. Required if index name is declared
- Ascending order: string. Only "asc" and "desc" are supported. Not required. Default is "asc"

##### unique

unique index. The format and usage are the same as index

##### null

Indicates that a null value can be stored. Note: If the type of the struct field is a pointer type, the default is nullable

##### unsigned

unsigned



#### dao layer interface

##### InsertXXX

Insert record

##### UpsertXXX

If there is a conflicting field value, such as a primary key value or a field value with a unique index, perform an update operation, otherwise an insert operation.

##### UpsertXXXNoneZero

Same as UpsertXXX. The difference is that only the non-[zero value](https://golang.org/ref/spec#The_zero_value) defined by the non-Go language specification is inserted or updated.

##### DeleteXXXs

Delete multiple records

##### UpdateXXX

Update record

##### UpdateXXXNoneZero

Same as UpdateXXX. The difference is that only the non-[zero value](https://golang.org/ref/spec#The_zero_value) defined by the non-Go language specification is updated. 

##### UpdateXXXs

Update multiple records

##### UpdateXXXsNoneZero

Same as UpdateXXXs. The difference is that only the non-[zero value](https://golang.org/ref/spec#The_zero_value) defined by the non-Go language specification is updated.

##### GetXXX

 Select record

##### SelectXXXs

Select multiple records

##### CountXXXs

Count multiple records

##### PageXXXs

Pagination

##### Transaction
Example：
```go
func (receiver *StockImpl) processExcel(ctx context.Context, f multipart.File, sheet string) (err error) {
	types := []string{"food", "cook"}
	var (
		xlsx *excelize.File
		rows [][]string
		tx   ddl.Tx
	)
	xlsx, err = excelize.OpenReader(f)
	if err != nil {
		return errors.Wrap(err, "")
	}
	rows, err = xlsx.GetRows(sheet)
	if err != nil {
		return errors.Wrap(err, "")
	}
	colNum := len(rows[0])
	rows = rows[1:]
	gdddb := ddl.GddDB{receiver.db}
	// begin task
	tx, err = gdddb.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "")
	}
	defer func() {
		if r := recover(); r != nil {
			if _err := tx.Rollback(); _err != nil {
				err = errors.Wrap(_err, "")
				return
			}
			if e, ok := r.(error); ok {
				err = errors.Wrap(e, "")
			} else {
				err = errors.New(fmt.Sprint(r))
			}
		}
	}()
	// Pass tx as an instance of the ddl.Querier interface to the dao layer implementation class
	mdao := dao.NewMaterialDao(tx)
	for _, item := range rows {
		if len(item) == 0 {
			goto END
		}
		row := make([]string, colNum)
		copy(row, item)
		name := row[0]
		price := cast.ToFloat32(row[1])
		spec := row[2]
		pieces := cast.ToInt(row[3])
		amount := cast.ToInt(row[4])
		note := row[5]
		totalMount := pieces * amount
		if _, err = mdao.Upsert(ctx, domain.Material{
			Name:        name,
			Amount:      amount,
			Price:       price,
			TotalAmount: totalMount,
			Spec:        spec,
			Pieces:      pieces,
			Type:        int8(sliceutils.IndexOf(sheet, types)),
			Note:        note,
		}); err != nil {
			// Error rollback
			if _err := tx.Rollback(); _err != nil {
				return errors.Wrap(_err, "")
			}
			return errors.Wrap(err, "")
		}
	}
END:
	// commit at last
	if err = tx.Commit(); err != nil {
		if _err := tx.Rollback(); _err != nil {
			return errors.Wrap(_err, "")
		}
		return errors.Wrap(err, "")
	}
	return err
}
```

#### Query Dsl

##### Example

```go
func ExampleCriteria() {
	
	query := C().Col("name").Eq(Literal("wubin")).
              Or(C().Col("school").Eq(Literal("havard"))).
              And(C().Col("age").Eq(Literal(18)))
	fmt.Println(query.Sql())

	query = C().Col("name").Eq(Literal("wubin")).
              Or(C().Col("school").Eq(Literal("havard"))).
              And(C().Col("delete_at").IsNotNull())
	fmt.Println(query.Sql())

	query = C().Col("name").Eq(Literal("wubin")).
              Or(C().Col("school").In(Literal("havard"))).
              And(C().Col("delete_at").IsNotNull())
	fmt.Println(query.Sql())

	query = C().Col("name").Eq(Literal("wubin")).
              Or(C().Col("school").In(Literal([]string{"havard", "beijing unv"}))).
              And(C().Col("delete_at").IsNotNull())
	fmt.Println(query.Sql())

	var d int
	var e int
	d = 10
	e = 5

	query = C().Col("name").Eq(Literal("wubin")).
              Or(C().Col("age").In(Literal([]*int{&d, &e}))).
              And(C().Col("delete_at").IsNotNull())
	fmt.Println(query.Sql())

	// Output:
	// ((`name` = 'wubin' or `school` = 'havard') and `age` = '18')
	// ((`name` = 'wubin' or `school` = 'havard') and `delete_at` is not null)
	// ((`name` = 'wubin' or `school` in ('havard')) and `delete_at` is not null)
	// ((`name` = 'wubin' or `school` in ('havard','beijing unv')) and `delete_at` is not null)
	// ((`name` = 'wubin' or `age` in ('10','5')) and `delete_at` is not null)
}
```



##### API

Mainly includes 1 interface and 6 structs



###### Q interface

```go
type Q interface {
	Sql() string
	And(q Q) Q
	Or(q Q) Q
}
```

- Call the Sql method to return the finally spliced where statement

- And method means "and" in sql

- Or method means "or" in sql



###### criteria struct

```go
type criteria struct {
	col  string
	val  Val
	asym arithsymbol.ArithSymbol
}
```

- col represents the name of the table field
- val represents the table field value
- asym represents an arithmetic operator, optional values:
    - Eq: `=`
    - Ne: `!=`
    - Gt: `>`
    - Lt: `<`
    - Gte: `>=`
    - Lte: `<=`
    - Is: `is`
    - Not: `is not`
    - In: `in`



Val struct

```go
type Val struct {
	Data interface{}
	Type valtypeenum.ValType
}
```

- Data represents the table field value
- Type represents the value type, optional:
    - Func: Represents database built-in functions or expressions composed of built-in functions
    - Null: Represents the null of the database
    - Literal: Represents the value different from Func and Null



###### where struct

```
type where struct {
   lsym     logicsymbol.LogicSymbol
   children []Q
}
```

- lsym represents a logical operator, optional values:
    - And: `and`
    - Or: `or`
- children represents sub-query conditions, and groups are formed by the logical relationship represented by lsym, and every two sub-conditions form a group. Each sub-condition can be either a `criteria` or a `where`. As long as the Q interface is implemented, a sub-condition can be made, and a group can be formed by lsym and another sub-condition.



###### Page struct

```go
type Page struct {
	Orders []Order
	Offset int
	Size   int
}
```

- Orders means sorting
- Offset means offset
- Size indicates how many rows there are on a page



Order struct

```go
type Order struct {
	Col  string
	Sort sortenum.Sort
}
```

- Col means table field
- Sort represents ascending and descending order, optional value: asc/desc



PageRet struct

```go
type PageRet struct {
	Items    interface{}
	PageNo   int
	PageSize int
	Total    int
	HasNext  bool
}
```

- Items represent row data
- PageNo represents the current page number
- PageSize indicates how many rows there are on a page
- Total means total
- HasNext indicates whether there is a next page


### TODO

+ [x] dao layer supports transaction
+ [ ] Support index update
+ [ ] Support foreign keys




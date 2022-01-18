## ddl

DDL and dao layer generation command line tool based on [jmoiron/sqlx](https://github.com/jmoiron/sqlx).



<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
### TOC

- [Features](#features)
- [Flags](#flags)
- [Quickstart](#quickstart)
- [API](#api)
  - [Example](#example)
  - [Tags](#tags)
    - [pk](#pk)
    - [auto](#auto)
    - [type](#type)
    - [default](#default)
    - [extra](#extra)
    - [index](#index)
    - [unique](#unique)
    - [null](#null)
    - [unsigned](#unsigned)
  - [Dao layer code](#dao-layer-code)
    - [CRUD](#crud)
    - [Transaction](#transaction)
  - [Query Dsl](#query-dsl)
    - [Example](#example-1)
    - [Q](#q)
    - [criteria](#criteria)
    - [Val](#val)
    - [where](#where)
- [TODO](#todo)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->



### Features

- Create/Update table from go struct
- Create/Update go struct from table
- Generate dao layer code with basic crud operations



### Flags

```shell
➜  ~ go-doudou ddl -h
migration tool between database table structure and golang struct

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
```



### Quickstart

- Install go-doudou

  ```shell
  go get -v github.com/unionj-cloud/go-doudou@v0.9.5
  ```

- Clone demo repository

  ```
  git clone git@github.com:unionj-cloud/ddldemo.git
  ```

- Update database table struct and generate dao layer code

  ```shell
  go-doudou ddl --dao --pre=ddl_
  ```

  ```shell
  ➜  ddldemo git:(main) ls -la dao
  total 56
  drwxr-xr-x   6 wubin1989  staff   192  9  1 00:28 .
  drwxr-xr-x  14 wubin1989  staff   448  9  1 00:28 ..
  -rw-r--r--   1 wubin1989  staff   953  9  1 00:28 base.go
  -rw-r--r--   1 wubin1989  staff    45  9  1 00:28 userdao.go
  -rw-r--r--   1 wubin1989  staff  9125  9  1 00:28 userdaoimpl.go
  -rw-r--r--   1 wubin1989  staff  5752  9  1 00:28 userdaosql.go
  ```

- Run main function

  ```
  ➜  ddldemo git:(main) go run main.go       
  INFO[0000] user jack's id is 14                         
  INFO[0000] returned user jack's id is 14                
  INFO[0000] returned user jack's average score is 97.534
  ```

- Delete domain and dao folder

  ```shell
  ➜  ddldemo git:(main) rm -rf dao && rm -rf domain
  ```

- Generate go struct and dao layer code from database

  ```shell
  go-doudou ddl --reverse --dao --pre=ddl_
  ```

  ```shell
  ➜  ddldemo git:(main) ✗ ll
  total 272
  -rw-r--r--  1 wubin1989  staff   1.0K  9  1 00:27 LICENSE
  -rw-r--r--  1 wubin1989  staff    85B  9  1 00:27 Makefile
  -rw-r--r--  1 wubin1989  staff     9B  9  1 00:27 README.md
  drwxr-xr-x  6 wubin1989  staff   192B  9  1 00:35 dao
  drwxr-xr-x  3 wubin1989  staff    96B  9  1 00:35 domain
  -rw-r--r--  1 wubin1989  staff   339B  9  1 00:27 go.mod
  -rw-r--r--  1 wubin1989  staff   116K  9  1 00:27 go.sum
  -rw-r--r--  1 wubin1989  staff   2.0K  9  1 00:27 main.go
  ```

- Run main function again

  ```
  ➜  ddldemo git:(main) ✗ go run main.go                          
  INFO[0000] user jack's id is 15                         
  INFO[0000] returned user jack's id is 15                
  INFO[0000] returned user jack's average score is 97.534 
  ```

  

### API

#### Example
```go
package domain

import "time"

type Base struct {
	CreateAt *time.Time `dd:"default:CURRENT_TIMESTAMP"`
	UpdateAt *time.Time `dd:"default:CURRENT_TIMESTAMP;extra:ON UPDATE CURRENT_TIMESTAMP"`
	DeleteAt *time.Time
}
```
```go
package domain

//dd:table
type Book struct {
  ID          int `dd:"pk;auto"`
  UserId      int `dd:"type:int"`
  PublisherId int	`dd:"fk:ddl_publisher,id,fk_publisher,ON DELETE CASCADE ON UPDATE NO ACTION"`

  Base
}
```
```go
package domain

//dd:table
type Publisher struct {
  ID   int `dd:"pk;auto"`
  Name string

  Base
}
```
```go
package domain

import "time"

//dd:table
type User struct {
  ID         int    `dd:"pk;auto"`
  Name       string `dd:"index:name_phone_idx,2;default:'jack'"`
  Phone      string `dd:"index:name_phone_idx,1;default:'13552053960';extra:comment '手机号'"`
  Age        int    `dd:"unsigned"`
  No         int    `dd:"type:int;unique"`
  UniqueCol  int    `dd:"type:int;unique:unique_col_idx,1"`
  UniqueCol2 int    `dd:"type:int;unique:unique_col_idx,2"`
  School     string `dd:"null;default:'harvard';extra:comment '学校'"`
  IsStudent  bool
  ArriveAt *time.Time `dd:"type:datetime;extra:comment '到货时间'"`
  Status   int8       `dd:"type:tinyint(4);extra:comment '0进行中
1完结
2取消'"`

  Base
}
```



#### Tags

##### pk

Primary key

##### auto

Autoincrement

##### type

Column type. Not required.

| Go Type（pointer） | Column Type  |
| :----------------: | :----------: |
| int, int16, int32  |     int      |
|       int64        |    bigint    |
|      float32       |    float     |
|      float64       |    double    |
|       string       | varchar(255) |
|     bool, int8     |   tinyint    |
|     time.Time      |   datetime   |
|  decimal.Decimal   | decimal(6,2) |

##### default

Default value. If value was database built-in function or expression made by built-in functions, not need single quote marks. If value was literal value, it should be quoted by single quote marks.

##### extra

Extra definition. Example: "on update CURRENT_TIMESTAMP"，"comment 'cellphone number'"  
**Note：don't use ; and : in comment**

##### index

- Format："index:Name,Order,Sort" or "index"
- Name: index name. string. If multiple fields use the same index name, the index will be created as composite index. Not required. Default index name is column name + _idx
- Order：int
- Sort：string. Only accept asc and desc. Not required. Default is asc

##### unique

Unique index. Usage is the same as index.

##### null

Nullable. **Note: if the field is a pointer, null is default.**

##### unsigned

Unsigned

##### fk

- Format："fk:ReferenceTableName,ReferenceTablePrimaryKey,Constraint,Action"  
- ReferenceTableName: reference table name
- ReferenceTablePrimaryKey: reference table primary key such as `id`
- Constraint: foreign key constraint such as `fk_publisher`
- Action: for example: `ON DELETE CASCADE ON UPDATE NO ACTION`



#### Dao layer code

##### CRUD

```go
package dao

import (
  "context"
  "github.com/unionj-cloud/go-doudou/ddl/query"
)

type Base interface {
  Insert(ctx context.Context, data interface{}) (int64, error)
  Upsert(ctx context.Context, data interface{}) (int64, error)
  UpsertNoneZero(ctx context.Context, data interface{}) (int64, error)
  DeleteMany(ctx context.Context, where query.Q) (int64, error)
  Update(ctx context.Context, data interface{}) (int64, error)
  UpdateNoneZero(ctx context.Context, data interface{}) (int64, error)
  UpdateMany(ctx context.Context, data interface{}, where query.Q) (int64, error)
  UpdateManyNoneZero(ctx context.Context, data interface{}, where query.Q) (int64, error)
  Get(ctx context.Context, id interface{}) (interface{}, error)
  SelectMany(ctx context.Context, where ...query.Q) (interface{}, error)
  CountMany(ctx context.Context, where ...query.Q) (int, error)
  PageMany(ctx context.Context, page query.Page, where ...query.Q) (query.PageRet, error)
}
```



##### Transaction
Example：
```go
func (receiver *StockImpl) processExcel(ctx context.Context, f multipart.File, sheet string) (err error) {
	types := []string{"food", "tool"}
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
    gdddb := wrapper.GddDB{receiver.db}
	// begin transaction
	tx, err = gdddb.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "")
	}
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			if e, ok := r.(error); ok {
				err = errors.Wrap(e, "")
			} else {
				err = errors.New(fmt.Sprint(r))
			}
		}
	}()
	// inject tx as ddl.Querier into dao layer implementation instance
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
			// rollback if err != nil
			_ = tx.Rollback()
			return errors.Wrap(err, "")
		}
	}
END:
	// commit
	if err = tx.Commit(); err != nil {
        _ = tx.Rollback()
		return errors.Wrap(err, "")
	}
	return err
}
```



#### Query Dsl

##### Example

```go
func ExampleCriteria() {

	query := C().Col("name").Eq("wubin").Or(C().Col("school").Eq("havard")).And(C().Col("age").Eq(18))
	fmt.Println(query.Sql())

	query = C().Col("name").Eq("wubin").Or(C().Col("school").Eq("havard")).And(C().Col("delete_at").IsNotNull())
	fmt.Println(query.Sql())

	query = C().Col("name").Eq("wubin").Or(C().Col("school").In("havard")).And(C().Col("delete_at").IsNotNull())
	fmt.Println(query.Sql())

	query = C().Col("name").Eq("wubin").Or(C().Col("school").In([]string{"havard", "beijing unv"})).And(C().Col("delete_at").IsNotNull())
	fmt.Println(query.Sql())

	query = C().Col("name").Eq("wubin").Or(C().Col("age").In([]int{5, 10})).And(C().Col("delete_at").IsNotNull())
	fmt.Println(query.Sql())

	query = C().Col("name").Ne("wubin").Or(C().Col("create_at").Lt("now()"))
	fmt.Println(query.Sql())

	page := Page{
		Orders: []Order{
			{
				Col:  "create_at",
				Sort: sortenum.Desc,
			},
		},
		Offset: 20,
		Size:   10,
	}
	page = page.Order(Order{
		Col:  "score",
		Sort: sortenum.Asc,
	})
	page = page.Limit(30, 5)
	fmt.Println(page.Sql())
	pageRet := NewPageRet(page)
	fmt.Println(pageRet.PageNo)

	fmt.Println(P().Order(Order{
		Col:  "score",
		Sort: sortenum.Asc,
	}).Limit(20, 10).Sql())

	query = C().Col("name").Eq("wubin").Or(C().Col("school").Eq("havard")).
		And(C().Col("age").Eq(18)).
		Or(C().Col("score").Gte(90))
	fmt.Println(query.Sql())

	page = P().Order(Order{
		Col:  "create_at",
		Sort: sortenum.Desc,
	}).Limit(0, 1)
	var where Q
	where = C().Col("project_id").Eq(1)
	where = where.And(C().Col("delete_at").IsNull())
	where = where.Append(page)
	fmt.Println(where.Sql())

	where = C().Col("project_id").Eq(1)
	where = where.And(C().Col("delete_at").IsNull())
	where = where.Append(String("for update"))
	fmt.Println(where.Sql())

	where = C().Col("cc.project_id").Eq(1)
	where = where.And(C().Col("cc.delete_at").IsNull())
	where = where.Append(String("for update"))
	fmt.Println(where.Sql())

	where = C().Col("cc.survey_id").Eq("abc").
		And(C().Col("cc.year").Eq(2021)).
		And(C().Col("cc.month").Eq(10)).
		And(C().Col("cc.stat_type").Eq(2)).Append(String("for update"))
	fmt.Println(where.Sql())

    where = C().Col("cc.name").Like("%ba%")
    fmt.Println(where.Sql())

	// Output:
	//((`name` = ? or `school` = ?) and `age` = ?) [wubin havard 18]
	//((`name` = ? or `school` = ?) and `delete_at` is not null) [wubin havard]
	//((`name` = ? or `school` in (?)) and `delete_at` is not null) [wubin havard]
	//((`name` = ? or `school` in (?,?)) and `delete_at` is not null) [wubin havard beijing unv]
	//((`name` = ? or `age` in (?,?)) and `delete_at` is not null) [wubin 5 10]
	//(`name` != ? or `create_at` < ?) [wubin now()]
	//order by `create_at` desc,`score` asc limit ?,? [30 5]
	//7
	//order by `score` asc limit ?,? [20 10]
	//(((`name` = ? or `school` = ?) and `age` = ?) or `score` >= ?) [wubin havard 18 90]
	//(`project_id` = ? and `delete_at` is null) order by `create_at` desc limit ?,? [1 0 1]
	//(`project_id` = ? and `delete_at` is null) for update [1]
	//(cc.`project_id` = ? and cc.`delete_at` is null) for update [1]
	//(((cc.`survey_id` = ? and cc.`year` = ?) and cc.`month` = ?) and cc.`stat_type` = ?) for update [abc 2021 10 2]
    //cc.`name` like ? [%ba%]
}
```



### TODO

+ [x] Support transaction in dao layer
+ [x] Support index update
+ [x] Support foreign key




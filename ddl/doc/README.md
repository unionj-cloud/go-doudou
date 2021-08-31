## ddl

DDL and dao layer generation command line tool based on [jmoiron/sqlx](https://github.com/jmoiron/sqlx).



### Features

- Create/Update table from go struct

- Create/Update go struct from table

- Generate dao layer code with basic crud operations

  

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
### TOC

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



### Flags

```shell
➜  ~ go-doudou ddl -h
WARN[0000] Error loading .env file: open /Users/.env: no such file or directory 
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
  go get -v -u github.com/unionj-cloud/go-doudou/...@v0.6.0
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
//dd:table
type User struct {
	ID        int    `dd:"pk;auto"`
	Name      string `dd:"index:name_phone_idx,2;default:'jack'"`
	Phone     string `dd:"index:name_phone_idx,1;default:'13552053960';extra:comment 'cellphone number'"`
	Age       int    `dd:"index"`
	No        int    `dd:"unique"`
	School    string `dd:"null;default:'harvard';extra:comment 'school'"`
	IsStudent bool

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



#### Dao layer code

##### CRUD

```go
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
	gdddb := ddl.GddDB{receiver.db}
	// begin transaction
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
			if _err := tx.Rollback(); _err != nil {
				return errors.Wrap(_err, "")
			}
			return errors.Wrap(err, "")
		}
	}
END:
	// commit
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



##### Q

```go
type Q interface {
	Sql() string
	And(q Q) Q
	Or(q Q) Q
}
```



##### criteria

```go
type criteria struct {
	col  string
	val  Val
	asym arithsymbol.ArithSymbol
}
```

- col: column name
- val: value
- asym：
  - Eq: `=`
  - Ne: `!=`
  - Gt: `>`
  - Lt: `<`
  - Gte: `>=`
  - Lte: `<=`
  - Is: `is`
  - Not: `is not`
  - In: `in`



##### Val

```go
type Val struct {
	Data interface{}
	Type valtypeenum.ValType
}
```

- Data: value
- Type
  - Func: database built-in function or expression made by built-in functions
  - Null: null
  - Literal: Literal value



##### where

```
type where struct {
   lsym     logicsymbol.LogicSymbol
   children []Q
}
```

- lsym
  - And: `and`
  - Or: `or`
  
- children: sub queries

  

### TODO

+ [x] Support transaction in dao layer
+ [ ] Support index update
+ [ ] Support foreign key




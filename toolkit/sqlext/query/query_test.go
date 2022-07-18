package query

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/unionj-cloud/go-doudou/toolkit/sqlext/sortenum"
	"testing"
)

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
}

func TestCriteriaAppend(t *testing.T) {
	sqlStatement := C().Col("name").Eq("wubin").Or(C().Col("school").Eq("havard")).
		And(C().Col("age").Eq(18)).
		Or(C().Col("score").Gte(90).And(C().Col("height").Gt(160).End(String("anything")).And(C().Col("height").Lte(170).
			Append(String("and favourite = 'Go'")))))
	str, _ := sqlStatement.Sql()
	require.Equal(t, "(((`name` = ? or `school` = ?) and `age` = ?) or (`score` >= ? and (`height` > ? and (`height` <= ? and favourite = 'Go'))))", str)
}

func TestWhereAppend(t *testing.T) {
	sqlStatement := C().Col("name").Eq("wubin").Or(C().Col("school").Eq("havard")).
		And(C().Col("age").Eq(18)).
		Or(C().Col("score").Gte(90).
			And(C().Col("height").Gt(160).And(C().Col("height").Lte(170))).
			Append(String("and favourite = 'Go'")))
	str, _ := sqlStatement.Sql()
	require.Equal(t, "(((`name` = ? or `school` = ?) and `age` = ?) or ((`score` >= ? and (`height` > ? and `height` <= ?)) and favourite = 'Go'))", str)
}

func ExampleEnd() {
	page := P().Order(Order{
		Col:  "create_at",
		Sort: sortenum.Desc,
	}).Limit(0, 1)
	var where Q
	where = C().Col("project_id").Eq(1)
	where = where.And(C().Col("delete_at").IsNull())
	where = where.End(page)
	fmt.Println(where.Sql())

	where = C().Col("project_id").Eq(1)
	where = where.And(C().Col("delete_at").IsNull())
	where = where.End(String("for update"))
	fmt.Println(where.Sql())

	where = C().Col("cc.project_id").Eq(1)
	where = where.And(C().Col("cc.delete_at").IsNull())
	where = where.End(String("for update"))
	fmt.Println(where.Sql())

	where = C().Col("cc.survey_id").Eq("abc").
		And(C().Col("cc.year").Eq(2021)).
		And(C().Col("cc.month").Eq(10)).
		And(C().Col("cc.stat_type").Eq(2)).End(String("for update"))
	fmt.Println(where.Sql())

	where = C().Col("cc.name").Like("%ba%")
	fmt.Println(where.Sql())

	page = P().Order(Order{
		Col:  "user.create_at",
		Sort: sortenum.Desc,
	}).Limit(0, 1)
	where = C().Col("project_id").Eq(1)
	where = where.And(C().Col("delete_at").IsNull())
	where = where.End(page)
	fmt.Println(where.Sql())

	where = C().Col("delete_at").IsNull().And(C().Col("op_code").NotIn([]int{1, 2, 3}))
	fmt.Println(where.Sql())

	// Output:
	//(`project_id` = ? and `delete_at` is null) order by `create_at` desc limit ?,? [1 0 1]
	//(`project_id` = ? and `delete_at` is null) for update [1]
	//(cc.`project_id` = ? and cc.`delete_at` is null) for update [1]
	//(((cc.`survey_id` = ? and cc.`year` = ?) and cc.`month` = ?) and cc.`stat_type` = ?) for update [abc 2021 10 2]
	//cc.`name` like ? [%ba%]
	//(`project_id` = ? and `delete_at` is null) order by user.`create_at` desc limit ?,? [1 0 1]
	//(`delete_at` is null and `op_code` not in (?,?,?)) [1 2 3]
}

func TestWhereAppend2(t *testing.T) {
	var where Q
	where = C().Col("left_number").Gt(0).Or(C().Col("left_number").Lt(0))
	where = where.And(C().Col("name").Ne("感谢参与"))
	where = where.And(C().Col("delete_at").IsNull())
	page := P().Order(Order{
		Col:  "order",
		Sort: sortenum.Desc,
	}).Limit(0, 10)
	where = where.Append(page)
	str, _ := where.Sql()
	require.Equal(t, "(((`left_number` > ? or `left_number` < ?) and `name` != ?) and `delete_at` is null) order by `order` desc limit ?,?", str)
}

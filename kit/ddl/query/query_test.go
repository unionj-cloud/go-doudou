package query

import (
	"fmt"
)

func ExampleCriteria() {

	query := C().Col("name").Eq(Literal("wubin")).Or(C().Col("school").Eq(Literal("havard"))).And(C().Col("age").Eq(Literal(18)))
	fmt.Println(query.Sql())
	query = C().Col("name").Eq(Literal("wubin")).Or(C().Col("school").Eq(Literal("havard"))).And(C().Col("delete_at").IsNotNull())
	fmt.Println(query.Sql())
	//query = C().Col("name").Eq(Literal("wubin")).Or(C().Col("school").Eq(Literal("havard"))).And(C().Col("delete_at").Gt(Literal(time.Now().Format(constants.FORMAT))))

	query = C().Col("name").Eq(Literal("wubin")).Or(C().Col("school").In(Literal("havard"))).And(C().Col("delete_at").IsNotNull())
	fmt.Println(query.Sql())

	query = C().Col("name").Eq(Literal("wubin")).Or(C().Col("school").In(Literal([]string{"havard", "beijing unv"}))).And(C().Col("delete_at").IsNotNull())
	fmt.Println(query.Sql())

	// Output:
	// ((`name` = 'wubin' or `school` = 'havard') and `age` = '18')
	// ((`name` = 'wubin' or `school` = 'havard') and `delete_at` is not null)
	// ((`name` = 'wubin' or `school` in ('havard')) and `delete_at` is not null)
	// ((`name` = 'wubin' or `school` in ('havard','beijing unv')) and `delete_at` is not null)
}

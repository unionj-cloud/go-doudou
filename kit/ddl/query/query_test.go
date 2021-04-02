package query

import "fmt"

func ExampleCriteria() {

	query := C().Col("name").Eq("wubin").Or(C().Col("school").Eq("havard")).And(C().Col("age").Eq("18"))

	fmt.Println(query.Sql())
	// Output:
	// ((name = 'wubin' or school = 'havard') and age = '18')
}

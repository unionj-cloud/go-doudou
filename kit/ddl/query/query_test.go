package query

import (
	"fmt"
	"github.com/unionj-cloud/go-doudou/kit/constants"
	"time"
)

func ExampleCriteria() {

	query := C().Col("name").Eq("wubin").Or(C().Col("school").Eq("havard")).And(C().Col("age").Eq(18))
	fmt.Println(query.Sql())
	query = C().Col("name").Eq("wubin").Or(C().Col("school").Eq("havard")).And(C().Col("delete_at").Gt(time.Now().Format(constants.FORMAT)))
	fmt.Println(query.Sql())

	// Output:
	// ((`name` = 'wubin' or `school` = 'havard') and `age` = '18')
	// ((`name` = 'wubin' or `school` = 'havard') and `delete_at` > '2021-04-03 20:48:23')
}

package main

import (
	"context"
	"fmt"

	"github.com/unionj-cloud/go-doudou/v2/toolkit/gorm_gen/examples/biz"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/gorm_gen/examples/conf"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/gorm_gen/examples/dal"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/gorm_gen/examples/dal/query"
)

func init() {
	dal.DB = dal.ConnectDB(conf.MySQLDSN).Debug()
}

func main() {
	// start your project here
	fmt.Println("hello world")
	defer fmt.Println("bye~")

	query.SetDefault(dal.DB)
	biz.Query(context.Background())
}

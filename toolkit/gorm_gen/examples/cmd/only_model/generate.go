package main

import (
	"github.com/unionj-cloud/go-doudou/v2/toolkit/gorm_gen"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/gorm_gen/examples/conf"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/gorm_gen/examples/dal"
)

func init() {
	dal.DB = dal.ConnectDB(conf.MySQLDSN).Debug()

	prepare(dal.DB) // prepare table for generate
}

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "/tmp/gentest/query",
	})

	g.UseDB(dal.DB)

	g.GenerateAllTable()

	g.Execute()
}

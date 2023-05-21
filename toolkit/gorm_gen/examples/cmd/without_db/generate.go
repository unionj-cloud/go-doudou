package main

import (
	"github.com/unionj-cloud/go-doudou/v2/toolkit/gorm_gen"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/gorm_gen/examples/dal/model"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "../../dal/query",
		Mode:    gen.WithDefaultQuery,
	})

	// generate from struct in project
	g.ApplyBasic(model.Customer{})

	g.Execute()
}

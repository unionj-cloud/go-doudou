package caches

import (
	"fmt"
	"github.com/goccy/go-reflect"

	"github.com/samber/lo"
	"github.com/wubin1989/gorm/callbacks"

	"github.com/wubin1989/gorm"
)

func buildIdentifier(db *gorm.DB) string {
	// Build query identifier,
	//	for that reason we need to compile all arguments into a string
	//	and concat them with the SQL query itself

	callbacks.BuildQuerySQL(db)
	var (
		identifier string
		query      string
		queryArgs  string
	)
	query = db.Statement.SQL.String()
	vars := lo.Map[interface{}, interface{}](db.Statement.Vars, func(item interface{}, index int) interface{} {
		if reflect.ValueOf(item).Kind() == reflect.Ptr {
			return reflect.ValueOf(item).Elem().Interface()
		}
		return item
	})
	queryArgs = fmt.Sprintf("%v", vars)
	identifier = fmt.Sprintf("%s-%s", query, queryArgs)

	return identifier
}

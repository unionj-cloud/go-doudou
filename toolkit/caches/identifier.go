package caches

import (
	"fmt"
	"gorm.io/gorm/callbacks"

	"gorm.io/gorm"
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
	queryArgs = fmt.Sprintf("%v", db.Statement.Vars)
	identifier = fmt.Sprintf("%s-%s", query, queryArgs)

	return identifier
}

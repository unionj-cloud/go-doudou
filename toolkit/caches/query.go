package caches

import "github.com/wubin1989/gorm"

type Query struct {
	Tags         []string
	Dest         interface{}
	RowsAffected int64
}

func (q *Query) replaceOn(db *gorm.DB) {
	SetPointedValue(db.Statement.Dest, q.Dest)
	SetPointedValue(&db.Statement.RowsAffected, &q.RowsAffected)
}

package caches

import (
	"testing"

	"github.com/wubin1989/gorm"
)

func Test_buildIdentifier(t *testing.T) {
	db := &gorm.DB{}
	db.Statement = &gorm.Statement{}
	db.Statement.SQL.WriteString("TEST-SQL")
	db.Statement.Vars = append(db.Statement.Vars, "test", 123, 12.3, true, false, []string{"test", "me"})

	actual := buildIdentifier(db)
	expected := "TEST-SQL-[test 123 12.3 true false [test me]]"
	if actual != expected {
		t.Errorf("buildIdentifier expected to return `%s` but got `%s`", expected, actual)
	}
}

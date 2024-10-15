package database

import (
	"os"
	"testing"
)

func TestInitialize(t *testing.T) {
	dir := "testdata/testpg"
	defer os.RemoveAll(dir)
	gen := GetOrmGenerator("gorm")
	gen.Initialize(OrmGeneratorConfig{
		Driver:      "postgres",
		Dsn:         "host=localhost user=corteza password=corteza dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		TablePrefix: "tutorial",
		Dir:         dir,
	})
	gen.GenGrpc()
}

module gorm.io/gen/tests

go 1.16

require (
	github.com/mattn/go-sqlite3 v1.14.16 // indirect
	github.com/unionj-cloud/go-doudou/v2 v2.1.9-0.20230825031202-41ef70f1be6f
	gorm.io/driver/mysql v1.5.1-0.20230509030346-3715c134c25b
	gorm.io/driver/sqlite v1.4.4
	gorm.io/gen v0.3.19
	gorm.io/gorm v1.25.1-0.20230505075827-e61b98d69677
	gorm.io/hints v1.1.1 // indirect
	gorm.io/plugin/dbresolver v1.4.0
)

replace github.com/unionj-cloud/go-doudou/v2 v2.1.9-0.20230825031202-41ef70f1be6f => ../../../

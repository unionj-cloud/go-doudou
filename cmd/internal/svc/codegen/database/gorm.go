package database

import (
	"github.com/unionj-cloud/go-doudou/v2/toolkit/errorx"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gen"
	"gorm.io/gorm"
)

const (
	GormKind = "gorm"
)

func init() {
	RegisterOrmGenerator(GormKind, &GormGenerator{})
}

var _ IOrmGenerator = (*GormGenerator)(nil)

type GormGenerator struct {
	Driver string
	Dsn    string
	Dir    string
	g      *gen.Generator
}

func (gg *GormGenerator) svcGo() {
	gg.g.GenerateSvcGo()
}

func (gg *GormGenerator) svcImplGo() {
	//gg.g.GenerateSvcImplGo()
}

func (gg *GormGenerator) dto() {
	gg.g.Execute()
}

func (gg *GormGenerator) GenService() {
	gg.dto()
	gg.svcGo()
	gg.svcImplGo()
}

const (
	driverMysql     = "mysql"
	driverPostgres  = "postgres"
	driverSqlite    = "sqlite"
	driverSqlserver = "sqlserver"
	driverTidb      = "tidb"
)

func (gg *GormGenerator) Initialize(conf OrmGeneratorConfig) {
	gg.Dir = conf.Dir
	gg.Driver = conf.Driver
	gg.Dsn = conf.Dsn
	var db *gorm.DB
	var err error
	switch gg.Driver {
	case driverMysql, driverTidb:
		conf := mysql.Config{
			DSN: gg.Dsn, // data source name
		}
		db, err = gorm.Open(mysql.New(conf))
	case driverPostgres:
		conf := postgres.Config{
			DSN: gg.Dsn,
		}
		db, err = gorm.Open(postgres.New(conf))
	case driverSqlite:
		db, err = gorm.Open(sqlite.Open(gg.Dsn))
	case driverSqlserver:
		db, err = gorm.Open(sqlserver.Open(gg.Dsn))
	default:
		errorx.Panic("Not support driver")
	}
	if err != nil {
		errorx.Panic(err.Error())
	}
	g := gen.NewGenerator(gen.Config{
		RootDir:       gg.Dir,
		OutPath:       gg.Dir + "/query",
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
		// if you want to assign field which has a default value in the `Create` API, set FieldCoverable true, reference: https://gorm.io/docs/create.html#Default-Values
		FieldCoverable: false,
		// if you want to generate field with unsigned integer type, set FieldSignable true
		FieldSignable: false,
		// if you want to generate index tags from database, set FieldWithIndexTag true
		FieldWithIndexTag: true,
		// if you want to generate type tags from database, set FieldWithTypeTag true
		FieldWithTypeTag: true,
		// if you need unit tests for query code, set WithUnitTest true
		WithUnitTest: false,
	})
	g.WithJSONTagNameStrategy(func(n string) string { return n + ",omitempty" })
	g.UseDB(db)
	g.ApplyBasic(g.GenerateAllTable()...)
	gg.g = g
}

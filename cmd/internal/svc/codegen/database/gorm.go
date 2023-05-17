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
	//TODO implement me
	panic("implement me")
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
		RootDir:        gg.Dir,
		OutPath:        gg.Dir + "/query",
		FieldNullable:  false,
		FieldCoverable: true,
		FieldSignable:  true,
		Mode:           gen.WithDefaultQuery | gen.WithQueryInterface,
	})
	g.UseDB(db)
	g.ApplyBasic(g.GenerateAllTable()...)
	gg.g = g
}

package database

import (
	"github.com/unionj-cloud/go-doudou/v2/framework/database"
	"gorm.io/gen"
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
}

func (gg *GormGenerator) svcGo() {
	//TODO implement me
	panic("implement me")
}

func (gg *GormGenerator) svcImplGo() {
	//TODO implement me
	panic("implement me")
}

func (gg *GormGenerator) dto() {
	g := gen.NewGenerator(gen.Config{
		OutPath:        "query",
		FieldNullable:  false,
		FieldCoverable: true,
		FieldSignable:  true,
		Mode:           gen.WithDefaultQuery | gen.WithQueryInterface,
	})
	g.UseDB(database.Db)
	g.ApplyBasic(g.GenerateAllTable()...)
	g.Execute()
}

func (gg *GormGenerator) GenService() {
	gg.dto()
	gg.svcGo()
	gg.svcImplGo()
}

func (gg *GormGenerator) SetConfig(conf OrmGeneratorConfig) {
	gg.Dir = conf.Dir
	gg.Driver = conf.Driver
	gg.Dsn = conf.Dsn
}

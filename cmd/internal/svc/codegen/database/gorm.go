package database

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

func (gen *GormGenerator) svcGo() {
	//TODO implement me
	panic("implement me")
}

func (gen *GormGenerator) svcImplGo() {
	//TODO implement me
	panic("implement me")
}

func (gen *GormGenerator) dto() {
	//TODO implement me
	panic("implement me")
}

func (gen *GormGenerator) GenService() {
	gen.dto()
	gen.svcGo()
	gen.svcImplGo()
}

func (gen *GormGenerator) SetConfig(conf OrmGeneratorConfig) {
	gen.Dir = conf.Dir
	gen.Driver = conf.Driver
	gen.Dsn = conf.Dsn
}

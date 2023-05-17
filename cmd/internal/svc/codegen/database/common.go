package database

var OrmGeneratorRegistry = make(map[OrmKind]IOrmGenerator)

type OrmKind string

func RegisterOrmGenerator(kind OrmKind, instance IOrmGenerator) {
	OrmGeneratorRegistry[kind] = instance
}

func GetOrmGenerator(kind OrmKind) IOrmGenerator {
	if gen, ok := OrmGeneratorRegistry[kind]; ok {
		return gen
	}
	return nil
}

type OrmGeneratorConfig struct {
	Driver string
	Dsn    string
	Dir    string
}

type IOrmGenerator interface {
	svcGo()
	svcImplGo()
	dto()
	Initialize(conf OrmGeneratorConfig)
	GenService()
}

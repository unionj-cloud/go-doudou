package database

import (
	"github.com/samber/lo"
	"github.com/unionj-cloud/toolkit/astutils"
	"path/filepath"
	"strings"

	"github.com/gobwas/glob"
	"github.com/unionj-cloud/toolkit/errorx"
	"github.com/unionj-cloud/toolkit/executils"
	"github.com/unionj-cloud/toolkit/stringutils"
	"github.com/wubin1989/gen"
	"github.com/wubin1989/gen/field"
	"github.com/wubin1989/gorm"
	"github.com/wubin1989/mysql"
	"github.com/wubin1989/postgres"
	"github.com/wubin1989/sqlite"
	"github.com/wubin1989/sqlserver"
)

const (
	GormKind = "gorm"
)

func init() {
	gg := &GormGenerator{}
	ag := &AbstractBaseGenerator{
		impl:   gg,
		runner: executils.CmdRunner{},
	}
	gg.AbstractBaseGenerator = ag
	RegisterOrmGenerator(GormKind, gg)
}

var _ IOrmGenerator = (*GormGenerator)(nil)

type GormGenerator struct {
	*AbstractBaseGenerator
}

func (gg *GormGenerator) svcGo() {
	gg.g.GenerateSvcGo()
}

func (gg *GormGenerator) svcImplGrpc() {
	gg.g.GenerateSvcImplGrpc()
}

func (gg *GormGenerator) svcImplRest() {
	gg.g.GenerateSvcImplRest()
}

func (gg *GormGenerator) orm() {
	gg.g.Execute()
}

func (gg *GormGenerator) fix() {
	//dir, _ := filepath.Abs(gg.Dir)
	//var files []string
	//err := filepath.Walk(dir, astutils.Visit(&files))
	//if err != nil {
	//	panic(err)
	//}
	//for _, file := range files {
	//	if filepath.Ext(file) != ".go" {
	//		continue
	//	}
	//	source, err := ioutil.ReadFile(file)
	//	if err != nil {
	//		panic(err)
	//	}
	//	fileContent := string(source)
	//	fileContent = strings.ReplaceAll(fileContent, "github.com/wubin1989/gen", "github.com/wubin1989/gen")
	//	ioutil.WriteFile(file, []byte(fileContent), os.ModePerm)
	//}
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
	gg.Client = false
	gg.TablePrefix = strings.TrimSuffix(conf.TablePrefix, ".")
	gg.TableGlob = conf.TableGlob
	gg.TableExcludeGlob = conf.TableExcludeGlob
	gg.GenGenGo = conf.GenGenGo
	gg.CaseConverter = conf.CaseConverter
	gg.ProtoGenerator = conf.ProtoGenerator
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
		if stringutils.IsNotEmpty(gg.TablePrefix) {
			db.Exec(`set search_path='` + gg.TablePrefix + `'`)
		}
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
	g := gormgen.NewGenerator(gormgen.Config{
		RootDir:       gg.Dir,
		OutPath:       gg.Dir + "/query",
		Mode:          gormgen.WithDefaultQuery | gormgen.WithQueryInterface,
		FieldNullable: true,
		// if you want to assign field which has a default value in the `Create` API, set FieldCoverable true
		FieldCoverable: true,
		// if you want to generate field with unsigned integer type, set FieldSignable true
		FieldSignable: false,
		// if you want to generate index tags from database, set FieldWithIndexTag true
		FieldWithIndexTag: true,
		// if you want to generate type tags from database, set FieldWithTypeTag true
		FieldWithTypeTag: true,
		// if you need unit tests for query code, set WithUnitTest true
		WithUnitTest: false,
	})
	g.WithOpts(gormgen.FieldGORMTag("", func(tag field.GormTag) field.GormTag {
		if defaultTag, ok := tag["default"]; ok {
			if len(defaultTag) > 0 {
				t := defaultTag[0]
				idx := strings.Index(t, "::")
				if idx > -1 {
					tag.Set("default", t[:idx])
				}
			}
		}
		return tag
	}))
	g.WithJSONTagNameStrategy(func(n string, t string) string {
		tagContent := gg.CaseConverter(n)
		if conf.Omitempty {
			tagContent = tagContent + ",omitempty"
		}
		return tagContent
	})
	if stringutils.IsNotEmpty(conf.TypeMapping) {
		typeMapping := make(map[string]func(gorm.ColumnType) (dataType string))
		lo.ForEach(strings.Split(conf.TypeMapping, ","), func(item string, index int) {
			kv := strings.Split(item, ":")
			columnType := kv[0]
			goType := kv[1]
			typeMapping[columnType] = func(columnType gorm.ColumnType) (dataType string) {
				if n, ok := columnType.Nullable(); ok && n {
					return "*" + goType
				}
				return goType
			}
		})
		g.WithDataTypeMap(typeMapping)
	}
	g.UseDB(db)
	g.GenGenGo = gg.GenGenGo
	var models []interface{}
	if stringutils.IsNotEmpty(gg.TableGlob) {
		g.FilterTableGlob = glob.MustCompile(gg.TableGlob)
	}
	if stringutils.IsNotEmpty(gg.TableExcludeGlob) {
		g.ExcludeTableGlob = glob.MustCompile(gg.TableExcludeGlob)
	}

	if stringutils.IsEmpty(gg.TableGlob) && stringutils.IsEmpty(gg.TableExcludeGlob) {
		models = g.GenerateAllTable(
			gormgen.FieldType(conf.Soft, "gorm.DeletedAt"),
			gormgen.FieldGenType(conf.Soft, "Time"),
		)
	} else {
		models = g.GenerateFilteredTables(
			gormgen.FieldType(conf.Soft, "gorm.DeletedAt"),
			gormgen.FieldGenType(conf.Soft, "Time"))
	}
	g.ApplyBasic(models...)
	g.ConfigPackage = astutils.GetPkgPath(filepath.Join(gg.Dir, "config"))
	gg.g = g
}

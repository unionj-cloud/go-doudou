package database

import (
	"github.com/unionj-cloud/go-doudou/v2/framework/internal/config"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/cast"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/errorx"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/stringutils"
	"github.com/unionj-cloud/go-doudou/v2/toolkit/zlogger"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"time"
)

const (
	driverMysql      = "mysql"
	driverPostgres   = "postgres"
	driverSqlite     = "sqlite"
	driverSqlserver  = "sqlserver"
	driverTidb       = "tidb"
	driverClickhouse = "clickhouse"
)

var Db *gorm.DB

func init() {
	if cast.ToBoolOrDefault(config.GddDBDisableAutoConfigure.Load(), config.DefaultGddDBDisableAutoConfigure) {
		return
	}
	dsn := config.GddDBDsn.Load()
	if stringutils.IsEmpty(dsn) {
		return
	}
	driver := config.GddDBDriver.Load()
	if stringutils.IsEmpty(driver) {
		errorx.Panic("Database driver is missing")
	}
	var err error
	switch driver {
	case driverMysql, driverTidb:
		conf := mysql.Config{
			DSN:                           dsn, // data source name
			SkipInitializeWithVersion:     cast.ToBoolOrDefault(config.GddDBMysqlSkipInitializeWithVersion.Load(), config.DefaultGddDBMysqlSkipInitializeWithVersion),
			DefaultStringSize:             uint(cast.ToIntOrDefault(config.GddDBMysqlDefaultStringSize.Load(), config.DefaultGddDBMysqlDefaultStringSize)),
			DisableWithReturning:          cast.ToBoolOrDefault(config.GddDBMysqlDisableWithReturning.Load(), config.DefaultGddDBMysqlDisableWithReturning),
			DisableDatetimePrecision:      cast.ToBoolOrDefault(config.GddDBMysqlDisableDatetimePrecision.Load(), config.DefaultGddDBMysqlDisableDatetimePrecision),
			DontSupportRenameIndex:        cast.ToBoolOrDefault(config.GddDBMysqlDontSupportRenameIndex.Load(), config.DefaultGddDBMysqlDontSupportRenameIndex),
			DontSupportRenameColumn:       cast.ToBoolOrDefault(config.GddDBMysqlDontSupportRenameColumn.Load(), config.DefaultGddDBMysqlDontSupportRenameColumn),
			DontSupportForShareClause:     cast.ToBoolOrDefault(config.GddDBMysqlDontSupportForShareClause.Load(), config.DefaultGddDBMysqlDontSupportForShareClause),
			DontSupportNullAsDefaultValue: cast.ToBoolOrDefault(config.GddDBMysqlDontSupportNullAsDefaultValue.Load(), config.DefaultGddDBMysqlDontSupportNullAsDefaultValue),
			DontSupportRenameColumnUnique: cast.ToBoolOrDefault(config.GddDBMysqlDontSupportRenameColumnUnique.Load(), config.DefaultGddDBMysqlDontSupportRenameColumnUnique),
		}
		Db, err = gorm.Open(mysql.New(conf), &gorm.Config{})
	case driverPostgres:
		conf := postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: cast.ToBoolOrDefault(config.GddDBPostgresPreferSimpleProtocol.Load(), config.DefaultGddDBPostgresPreferSimpleProtocol),
			WithoutReturning:     cast.ToBoolOrDefault(config.GddDBPostgresWithoutReturning.Load(), config.DefaultGddDBPostgresWithoutReturning),
		}
		Db, err = gorm.Open(postgres.New(conf), &gorm.Config{})
	case driverSqlite:
		Db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	case driverSqlserver:
		Db, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	case driverClickhouse:
		Db, err = gorm.Open(clickhouse.Open(dsn), &gorm.Config{})
	default:
		errorx.Panic("Not support driver")
	}
	if err != nil {
		errorx.Panic(err.Error())
	}
	sqlDB, err := Db.DB()
	if err != nil {
		errorx.Panic(err.Error())
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(cast.ToIntOrDefault(config.GddDBMaxIdleConns.Load(), config.DefaultGddDBMaxIdleConns))

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(cast.ToIntOrDefault(config.GddDBMaxOpenConns.Load(), config.DefaultGddDBMaxOpenConns))

	maxLifetime, err := time.ParseDuration(config.GddDBConnMaxLifetime.Load())
	if err != nil {
		zlogger.Debug().Msgf("Parse %s %s as time.Duration failed: %s, use default %s instead.\n", string(config.GddDBConnMaxLifetime), config.GddDBConnMaxLifetime.Load(), err.Error(), config.DefaultGddDBConnMaxLifetime)
		maxLifetime = config.DefaultGddDBConnMaxLifetime
	}
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(maxLifetime)

	maxIdleTime, err := time.ParseDuration(config.GddDBConnMaxIdleTime.Load())
	if err != nil {
		zlogger.Debug().Msgf("Parse %s %s as time.Duration failed: %s, use default %s instead.\n", string(config.GddDBConnMaxIdleTime), config.GddDBConnMaxIdleTime.Load(), err.Error(), config.DefaultGddDBConnMaxIdleTime)
		maxIdleTime = config.DefaultGddDBConnMaxIdleTime
	}
	sqlDB.SetConnMaxIdleTime(maxIdleTime)
}

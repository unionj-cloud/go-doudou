package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstantsExists(t *testing.T) {
	assert.Equal(t, "mysql", DriverMysql)
	assert.Equal(t, "postgres", DriverPostgres)
	assert.Equal(t, "sqlite", DriverSqlite)
	assert.Equal(t, "sqlserver", DriverSqlserver)
	assert.Equal(t, "tidb", DriverTidb)
	assert.Equal(t, "clickhouse", DriverClickhouse)
	assert.Equal(t, "dm", DriverDm)
}

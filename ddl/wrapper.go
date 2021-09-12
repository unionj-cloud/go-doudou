package ddl

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

// DB wraps sqlx.Tx and sqlx.DB https://github.com/jmoiron/sqlx/issues/344#issuecomment-318372779
type DB interface {
	Querier
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (Tx, error)
	Close() error
}

// Tx transaction
type Tx interface {
	Querier
	Commit() error
	Rollback() error
}

// Querier common operations for sqlx.Tx and sqlx.DB
type Querier interface {
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Rebind(query string) string
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

// GddDB wraps sqlx.DB
type GddDB struct {
	*sqlx.DB
}

// GddTx wraps sqlx.Tx
type GddTx struct {
	*sqlx.Tx
}

// BeginTxx begins a transaction
func (g *GddDB) BeginTxx(ctx context.Context, opts *sql.TxOptions) (Tx, error) {
	tx, err := g.DB.BeginTxx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return GddTx{tx}, nil
}

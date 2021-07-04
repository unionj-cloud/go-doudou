package ddl

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type DB interface {
	Querier
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (Tx, error)
	Close() error
}

type Tx interface {
	Querier
	Commit() error
	Rollback() error
}

type Querier interface {
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Rebind(query string) string
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type GddDB struct {
	*sqlx.DB
}

type GddTx struct {
	*sqlx.Tx
}

func (g *GddDB) BeginTxx(ctx context.Context, opts *sql.TxOptions) (Tx, error) {
	tx, err := g.DB.BeginTxx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return GddTx{tx}, nil
}

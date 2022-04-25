package wrapper

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/unionj-cloud/go-doudou/toolkit/caller"
	"github.com/unionj-cloud/go-doudou/toolkit/sqlext/logger"
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
	BindNamed(query string, arg interface{}) (string, []interface{}, error)
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

// GddDB wraps sqlx.DB
type GddDB struct {
	*sqlx.DB
	logger logger.ISqlLogger
}

type GddDBOption func(*GddDB)

func WithLogger(logger logger.ISqlLogger) GddDBOption {
	return func(g *GddDB) {
		g.logger = logger
	}
}

func NewGddDB(db *sqlx.DB, options ...GddDBOption) DB {
	g := &GddDB{
		DB:     db,
		logger: logger.NewSqlLogger(logrus.StandardLogger()),
	}
	for _, opt := range options {
		opt(g)
	}
	return g
}

func (g *GddDB) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	if q, args, err := g.DB.BindNamed(query, arg); err != nil {
		return nil, errors.Wrap(err, caller.NewCaller().String())
	} else {
		g.logger.Log(q, args...)
	}
	return g.DB.NamedExecContext(ctx, query, arg)
}

func (g *GddDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	g.logger.Log(query, args...)
	return g.DB.ExecContext(ctx, query, args...)
}

func (g *GddDB) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	g.logger.Log(query, args...)
	return g.DB.GetContext(ctx, dest, query, args...)
}

func (g *GddDB) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	g.logger.Log(query, args...)
	return g.DB.SelectContext(ctx, dest, query, args...)
}

// BeginTxx begins a transaction
func (g *GddDB) BeginTxx(ctx context.Context, opts *sql.TxOptions) (Tx, error) {
	tx, err := g.DB.BeginTxx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &GddTx{tx, g.logger}, nil
}

// GddTx wraps sqlx.Tx
type GddTx struct {
	*sqlx.Tx
	logger logger.ISqlLogger
}

func (g *GddTx) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	if q, args, err := g.Tx.BindNamed(query, arg); err != nil {
		return nil, errors.Wrap(err, caller.NewCaller().String())
	} else {
		g.logger.Log(q, args...)
	}
	return g.Tx.NamedExecContext(ctx, query, arg)
}

func (g *GddTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	g.logger.Log(query, args...)
	return g.Tx.ExecContext(ctx, query, args...)
}

func (g *GddTx) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	g.logger.Log(query, args...)
	return g.Tx.GetContext(ctx, dest, query, args...)
}

func (g *GddTx) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	g.logger.Log(query, args...)
	return g.Tx.SelectContext(ctx, dest, query, args...)
}

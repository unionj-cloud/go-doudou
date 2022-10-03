package wrapper

import (
	"context"
	"database/sql"
	"github.com/go-redis/cache/v8"
	"github.com/jmoiron/sqlx"
	"github.com/lithammer/shortuuid/v4"
	"github.com/pkg/errors"
	"github.com/unionj-cloud/go-doudou/toolkit/caller"
	"github.com/unionj-cloud/go-doudou/toolkit/sqlext/logger"
	"time"
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
	logger      logger.SqlLogger
	cacheStore  *cache.Cache
	redisKeyTTL time.Duration
}

type GddDBOption func(*GddDB)

func WithLogger(logger logger.SqlLogger) GddDBOption {
	return func(g *GddDB) {
		g.logger = logger
	}
}

func WithCache(store *cache.Cache) GddDBOption {
	return func(g *GddDB) {
		g.cacheStore = store
	}
}

func WithRedisKeyTTL(ttl time.Duration) GddDBOption {
	return func(g *GddDB) {
		g.redisKeyTTL = ttl
	}
}

func NewGddDB(db *sqlx.DB, options ...GddDBOption) DB {
	g := &GddDB{
		DB:          db,
		logger:      logger.NewSqlLogger(),
		redisKeyTTL: time.Hour,
	}
	for _, opt := range options {
		opt(g)
	}
	return g
}

func (g *GddDB) NamedExecContext(ctx context.Context, query string, arg interface{}) (ret sql.Result, err error) {
	var (
		q    string
		args []interface{}
	)
	defer func() {
		g.logger.LogWithErr(ctx, err, nil, q, args...)
	}()
	q, args, err = g.DB.BindNamed(query, arg)
	err = errors.Wrap(err, caller.NewCaller().String())
	if err != nil {
		return nil, err
	}
	ret, err = g.DB.NamedExecContext(ctx, query, arg)
	err = errors.Wrap(err, caller.NewCaller().String())
	return
}

func (g *GddDB) ExecContext(ctx context.Context, query string, args ...interface{}) (ret sql.Result, err error) {
	defer func() {
		g.logger.LogWithErr(ctx, err, nil, query, args...)
	}()
	ret, err = g.DB.ExecContext(ctx, query, args...)
	err = errors.Wrap(err, caller.NewCaller().String())
	return
}

func (g *GddDB) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error) {
	hit := true
	defer func() {
		g.logger.LogWithErr(ctx, err, &hit, query, args...)
	}()
	if g.cacheStore != nil {
		err = g.cacheStore.Once(&cache.Item{
			Key:   shortuuid.NewWithNamespace(logger.PopulatedSql(query, args...)),
			Value: dest,
			TTL:   g.redisKeyTTL,
			Do: func(*cache.Item) (interface{}, error) {
				hit = false
				err = g.DB.GetContext(ctx, dest, query, args...)
				err = errors.Wrap(err, caller.NewCaller().String())
				return dest, err
			},
		})
		return
	}
	hit = false
	err = g.DB.GetContext(ctx, dest, query, args...)
	err = errors.Wrap(err, caller.NewCaller().String())
	return
}

func (g *GddDB) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error) {
	hit := true
	defer func() {
		g.logger.LogWithErr(ctx, err, &hit, query, args...)
	}()
	if g.cacheStore != nil {
		err = g.cacheStore.Once(&cache.Item{
			Key:   shortuuid.NewWithNamespace(logger.PopulatedSql(query, args...)),
			Value: dest,
			TTL:   g.redisKeyTTL,
			Do: func(*cache.Item) (interface{}, error) {
				hit = false
				err = g.DB.SelectContext(ctx, dest, query, args...)
				err = errors.Wrap(err, caller.NewCaller().String())
				return dest, err
			},
		})
		return
	}
	hit = false
	err = g.DB.SelectContext(ctx, dest, query, args...)
	err = errors.Wrap(err, caller.NewCaller().String())
	return
}

// BeginTxx begins a transaction
func (g *GddDB) BeginTxx(ctx context.Context, opts *sql.TxOptions) (Tx, error) {
	tx, err := g.DB.BeginTxx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &GddTx{tx, g.logger, g.cacheStore, g.redisKeyTTL}, nil
}

// GddTx wraps sqlx.Tx
type GddTx struct {
	*sqlx.Tx
	logger      logger.SqlLogger
	cacheStore  *cache.Cache
	redisKeyTTL time.Duration
}

func (g *GddTx) NamedExecContext(ctx context.Context, query string, arg interface{}) (ret sql.Result, err error) {
	var (
		q    string
		args []interface{}
	)
	defer func() {
		g.logger.LogWithErr(ctx, err, nil, q, args...)
	}()
	q, args, err = g.Tx.BindNamed(query, arg)
	err = errors.Wrap(err, caller.NewCaller().String())
	if err != nil {
		return nil, err
	}
	ret, err = g.Tx.NamedExecContext(ctx, query, arg)
	err = errors.Wrap(err, caller.NewCaller().String())
	return
}

func (g *GddTx) ExecContext(ctx context.Context, query string, args ...interface{}) (ret sql.Result, err error) {
	defer func() {
		g.logger.LogWithErr(ctx, err, nil, query, args...)
	}()
	ret, err = g.Tx.ExecContext(ctx, query, args...)
	err = errors.Wrap(err, caller.NewCaller().String())
	return
}

func (g *GddTx) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error) {
	hit := true
	defer func() {
		g.logger.LogWithErr(ctx, err, &hit, query, args...)
	}()
	if g.cacheStore != nil {
		err = g.cacheStore.Once(&cache.Item{
			Key:   shortuuid.NewWithNamespace(logger.PopulatedSql(query, args...)),
			Value: dest,
			TTL:   g.redisKeyTTL,
			Do: func(*cache.Item) (interface{}, error) {
				hit = false
				err = g.Tx.GetContext(ctx, dest, query, args...)
				err = errors.Wrap(err, caller.NewCaller().String())
				return dest, err
			},
		})
		return
	}
	hit = false
	err = g.Tx.GetContext(ctx, dest, query, args...)
	err = errors.Wrap(err, caller.NewCaller().String())
	return
}

func (g *GddTx) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) (err error) {
	hit := true
	defer func() {
		g.logger.LogWithErr(ctx, err, &hit, query, args...)
	}()
	if g.cacheStore != nil {
		err = g.cacheStore.Once(&cache.Item{
			Key:   shortuuid.NewWithNamespace(logger.PopulatedSql(query, args...)),
			Value: dest,
			TTL:   g.redisKeyTTL,
			Do: func(*cache.Item) (interface{}, error) {
				hit = false
				err = g.Tx.SelectContext(ctx, dest, query, args...)
				err = errors.Wrap(err, caller.NewCaller().String())
				return dest, err
			},
		})
		return
	}
	hit = false
	err = g.Tx.SelectContext(ctx, dest, query, args...)
	err = errors.Wrap(err, caller.NewCaller().String())
	return
}

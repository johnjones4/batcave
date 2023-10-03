package pgstore

import (
	"context"
	"main/core"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type pgWrapper struct {
	pool *pgxpool.Pool
}

type PGStore struct {
	pool core.Database
	log  logrus.FieldLogger
}

func New(ctx context.Context, log logrus.FieldLogger, url string) (*PGStore, error) {
	pool, err := pgxpool.Connect(ctx, url)
	if err != nil {
		return nil, err
	}

	impl := &PGStore{}
	impl.pool = &pgWrapper{pool}
	impl.log = log
	return impl, nil
}

func (p *pgWrapper) Exec(ctx context.Context, sql string, arguments ...interface{}) error {
	_, err := p.pool.Exec(ctx, sql, arguments...)
	return err
}

func (p *pgWrapper) Query(ctx context.Context, sql string, args ...interface{}) (core.Rows, error) {
	return p.pool.Query(ctx, sql, args...)
}

func (p *pgWrapper) QueryRow(ctx context.Context, sql string, args ...interface{}) core.Row {
	return p.pool.QueryRow(ctx, sql, args...)
}

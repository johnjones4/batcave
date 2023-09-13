package pgstore

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type PGStore struct {
	pool *pgxpool.Pool
	log  logrus.FieldLogger
}

func New(ctx context.Context, log logrus.FieldLogger, url string) (*PGStore, error) {
	pool, err := pgxpool.Connect(ctx, url)
	if err != nil {
		return nil, err
	}

	impl := &PGStore{}
	impl.pool = pool
	impl.log = log
	return impl, nil
}

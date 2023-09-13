package pgstore

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PGStore struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, url string) (*PGStore, error) {
	pool, err := pgxpool.Connect(ctx, url)
	if err != nil {
		return nil, err
	}

	impl := &PGStore{}
	impl.pool = pool
	return impl, nil
}

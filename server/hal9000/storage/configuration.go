package storage

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Configuration struct {
	DatabaseURL string
	UsersPath   string
	ClientsPath string
}

func Connect(configuration Configuration) (*pgxpool.Pool, error) {
	return pgxpool.Connect(context.Background(), configuration.DatabaseURL)
}

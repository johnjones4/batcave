package storage

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/johnjones4/hal-9000/server/hal9000/core"
)

type StateStore struct {
	pool *pgxpool.Pool
}

func NewStateStore(pool *pgxpool.Pool) *StateStore {
	return &StateStore{
		pool: pool,
	}
}

func (ss *StateStore) GetState(client core.Client) (string, error) {
	row := ss.pool.QueryRow(context.Background(), "SELECT state FROM states WHERE client = $1", client.ID)
	var state string
	err := row.Scan(&state)
	if err != nil {
		if err == pgx.ErrNoRows {
			return core.StateDefault, nil
		}
		return "", err
	}
	return state, nil
}

func (ss *StateStore) SetState(client core.Client, state string) error {
	rows, err := ss.pool.Query(context.Background(), "SELECT state FROM states WHERE client = $1", client.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		_, err = ss.pool.Exec(context.Background(), "UPDATE states SET state = $1 WHERE client = $2", state, client.ID)
		return err
	}
	_, err = ss.pool.Exec(context.Background(), "INSERT INTO states (state,client) VALUES ($1,$2)", state, client.ID)
	return err
}

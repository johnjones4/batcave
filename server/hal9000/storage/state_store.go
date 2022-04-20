package storage

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/johnjones4/hal-9000/server/hal9000/core"
)

type StateStore interface {
	GetState(client core.Client) (string, error)
	SetState(client core.Client, state string) error
}

type DatabaseStateStore struct {
	pool *pgxpool.Pool
}

func NewDatabaseStateStore(pool *pgxpool.Pool) StateStore {
	return &DatabaseStateStore{
		pool: pool,
	}
}

func (ss *DatabaseStateStore) GetState(client core.Client) (string, error) {
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

func (ss *DatabaseStateStore) SetState(client core.Client, state string) error {
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

type MemoryStateStore struct {
	states map[string]string
}

func NewMemoryStateStore() StateStore {
	return &MemoryStateStore{
		states: make(map[string]string),
	}
}

func (ss *MemoryStateStore) GetState(client core.Client) (string, error) {
	s, ok := ss.states[client.ID]
	if !ok {
		return core.StateDefault, nil
	}
	return s, nil
}

func (ss *MemoryStateStore) SetState(client core.Client, state string) error {
	ss.states[client.ID] = state
	return nil
}

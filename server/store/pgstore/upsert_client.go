package pgstore

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v4"
)

func (s *PGStore) UpsertClient(ctx context.Context, source string, clientId string, info any) error {
	infoStr, err := json.Marshal(info)
	if err != nil {
		return nil
	}

	now := time.Now().UTC()

	var dummy string
	err = s.pool.QueryRow(ctx, "SELECT client_id FROM clients_registry WHERE source = $1 AND client_id = $2", source, clientId).Scan(&dummy)
	if err != nil && err != pgx.ErrNoRows {
		return err
	}

	if err == pgx.ErrNoRows {
		_, err = s.pool.Exec(
			ctx,
			"INSERT INTO clients_registry (source, client_id, info, created, updated) VALUES ($1, $2, $3, $4, $5)",
			source,
			clientId,
			string(infoStr),
			now,
			now,
		)
		if err != nil {
			return err
		}
	} else {
		_, err = s.pool.Exec(
			ctx,
			"UPDATE clients_registry SET info = $1, updated = $2 WHERE source = $3 AND client_id = $4",
			string(infoStr),
			now,
			source,
			clientId,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

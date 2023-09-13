package pgstore

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v4"
)

func (s *PGStore) GetClient(ctx context.Context, source, clientId string, infoReceiver any) error {
	var info string
	err := s.pool.QueryRow(ctx, "SELECT info FROM clients_registry WHERE source = $1 AND client_id = $2", source, clientId).Scan(&info)
	if err == pgx.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(info), infoReceiver)
	if err != nil {
		return err
	}
	return nil
}

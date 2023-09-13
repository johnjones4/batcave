package pgstore

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v4"
)

func (s *PGStore) UpdateIntentEmbedding(ctx context.Context, intent string, embedding []float32) error {
	embeddingjson, err := json.Marshal(embedding)
	if err != nil {
		return err
	}

	var empty string
	err = s.pool.QueryRow(ctx, "SELECT intent_label FROM intents WHERE intent_label = $1", intent).Scan(&empty)
	if err != nil && err != pgx.ErrNoRows {
		return err
	} else if err != nil {
		_, err = s.pool.Exec(ctx, "INSERT INTO intents (intent_label, embedding) VALUES ($1, $2)", intent, string(embeddingjson))
		if err != nil {
			return err
		}
	} else {
		_, err = s.pool.Exec(ctx, "UPDATE intents SET embedding = $1 WHERE intent_label = $2", string(embeddingjson), intent)
		if err != nil {
			return err
		}
	}

	return nil
}

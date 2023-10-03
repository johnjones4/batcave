package pgstore

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v4"
)

func (s *PGStore) ClosestMatchingIntent(ctx context.Context, embedding []float32) (string, error) {
	embeddingjson, err := json.Marshal(embedding)
	if err != nil {
		return "", err
	}

	var intent string
	var distance float32
	err = s.pool.QueryRow(ctx, "SELECT intent_label, embedding <=> $1 as distance FROM intents ORDER BY distance LIMIT 1", string(embeddingjson)).Scan(&intent, &distance)
	if err == pgx.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	s.log.Debugf("Found intent %s with distance %f", intent, distance)
	if distance > 0.25 {
		return "", nil
	}

	return intent, nil
}

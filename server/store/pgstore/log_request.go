package pgstore

import (
	"context"
	"main/core"
	"time"
)

func (s *PGStore) LogRequest(ctx context.Context, res *core.Request) error {
	return s.pool.Exec(ctx, "INSERT INTO requests (event_id, timestamp, source, client_id, latitude, longitude, message_text) VALUES ($1,$2,$3,$4,$5,$6,$7)", res.EventId, time.Now().UTC(), res.Source, res.ClientID, res.Coordinate.Latitude, res.Coordinate.Longitude, res.Message.Text)
}

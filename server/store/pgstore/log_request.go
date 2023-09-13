package pgstore

import (
	"context"
	"main/core"
	"time"
)

// TODO log source
func (s *PGStore) LogRequest(ctx context.Context, res *core.Request) error {
	_, err := s.pool.Exec(ctx, "INSERT INTO requests (event_id, timestamp, client_id, latitude, longitude, message_text) VALUES ($1,$2,$3,$4,$5,$6)", res.EventId, time.Now().UTC(), res.ClientID, res.Coordinate.Latitude, res.Coordinate.Longitude, res.Message.Text)
	return err
}

package pgstore

import (
	"context"
	"main/core"
	"time"
)

func (s *PGStore) LogResponse(ctx context.Context, _ *core.Request, res *core.Response) error {
	_, err := s.pool.Exec(ctx, "INSERT INTO responses (event_id, timestamp, message_text, media_url, media_type, action) VALUES ($1,$2,$3,$4,$5,$6)", res.EventId, time.Now().UTC(), res.Message.Text, res.Media.URL, res.Media.Type, res.Action)
	return err
}

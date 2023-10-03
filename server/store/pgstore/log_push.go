package pgstore

import (
	"context"
	"main/core"
	"time"
)

func (s *PGStore) LogPush(ctx context.Context, clientId string, push *core.PushMessage) error {
	return s.pool.Exec(
		ctx,
		"INSERT INTO pushes (event_id, timestamp, client_id, message_text, media_url, media_type) VALUES ($1,$2,$3,$4,$5,$6)",
		push.EventId,
		time.Now().UTC(),
		clientId,
		push.Message.Text,
		push.Media.URL,
		push.Media.Type,
	)
}

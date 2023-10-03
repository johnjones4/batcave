package pgstore

import "context"

func (s *PGStore) ClearScheduledEvent(ctx context.Context, id string) error {
	return s.pool.Exec(ctx, "DELETE FROM scheduled_events WHERE event_id = $1", id)
}

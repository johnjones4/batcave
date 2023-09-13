package pgstore

import "context"

func (s *PGStore) ClearScheduledEvent(ctx context.Context, id string) error {
	_, err := s.pool.Exec(ctx, "DELETE FROM scheduled_events WHERE event_id = $1", id)
	return err
}

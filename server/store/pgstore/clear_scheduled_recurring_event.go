package pgstore

import "context"

func (s *PGStore) ClearScheduledRecurringEvent(ctx context.Context, id string) error {
	_, err := s.pool.Exec(ctx, "DELETE FROM recurring_events WHERE event_id = $1", id)
	return err
}

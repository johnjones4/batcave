package pgstore

import "context"

func (s *PGStore) ClearScheduledRecurringEvent(ctx context.Context, id string) error {
	return s.pool.Exec(ctx, "DELETE FROM recurring_events WHERE event_id = $1", id)
}

package pgstore

import (
	"context"
	"time"
)

func (s *PGStore) UpdateRecurringEventTimestamp(ctx context.Context, id string, stamp time.Time) error {
	_, err := s.pool.Exec(ctx, "UPDATE recurring_events SET last_run = $1 WHERE event_id = $2", stamp, id)
	return err
}

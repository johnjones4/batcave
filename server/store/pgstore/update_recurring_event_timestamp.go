package pgstore

import (
	"context"
	"time"
)

func (s *PGStore) UpdateRecurringEventTimestamp(ctx context.Context, id string, stamp time.Time) error {
	return s.pool.Exec(ctx, "UPDATE recurring_events SET last_run = $1 WHERE event_id = $2", stamp, id)
}

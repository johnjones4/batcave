package pgstore

import (
	"context"
	"encoding/json"
	"main/core"
	"time"

	"github.com/google/uuid"
)

func (s *PGStore) ScheduleRecurringEvent(ctx context.Context, event *core.ScheduledRecurringEvent) error {
	event.ID = uuid.NewString()
	info, err := json.Marshal(event.Info)
	if err != nil {
		return err
	}
	_, err = s.pool.Exec(
		ctx,
		"INSERT INTO recurring_events (event_id, source, client_id, intent, scheduled, last_run, created, info) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)",
		event.ID,
		event.Source,
		event.ClientId,
		event.Intent,
		event.Scheduled,
		event.LastRun,
		time.Now().UTC(),
		string(info),
	)
	if err != nil {
		return err
	}
	return nil
}

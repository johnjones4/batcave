package pgstore

import (
	"context"
	"encoding/json"
	"main/core"
	"time"

	"github.com/google/uuid"
)

func (s *PGStore) ScheduleEvent(ctx context.Context, event *core.ScheduledEvent) error {
	event.ID = uuid.NewString()
	info, err := json.Marshal(event.Info)
	if err != nil {
		return err
	}
	_, err = s.pool.Exec(
		ctx,
		"INSERT INTO scheduled_events (event_id, event_type, scheduled, created, info) VALUES ($1,$2,$3,$4,$5)",
		event.ID,
		event.EventType,
		event.Scheduled.UTC(),
		time.Now().UTC(),
		string(info),
	)
	if err != nil {
		return err
	}
	return nil
}

package pgstore

import (
	"context"
	"main/core"
	"time"
)

func (s *PGStore) ReadyEvents(ctx context.Context, frontier time.Time, eventType string, infoParser func(event *core.ScheduledEvent, info string) error) ([]core.ScheduledEvent, error) {
	rows, err := s.pool.Query(ctx, "SELECT event_id, source, client_id, event_type, scheduled, info FROM scheduled_events WHERE event_type = $1 AND scheduled <= $2", eventType, frontier.UTC())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]core.ScheduledEvent, 0)

	for rows.Next() {
		var row core.ScheduledEvent
		var info string
		err = rows.Scan(&row.ID, &row.Source, &row.ClientId, &row.EventType, &row.Scheduled, &info)
		if err != nil {
			return nil, err
		}

		err = infoParser(&row, info)
		if err != nil {
			return nil, err
		}

		events = append(events, row)
	}

	return events, nil
}

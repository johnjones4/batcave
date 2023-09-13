package pgstore

import (
	"context"
	"main/core"
)

func (s *PGStore) GetReadyEvents(ctx context.Context, eventType string, infoParser func(event *core.ScheduledEvent, info string) error) ([]core.ScheduledEvent, error) {
	rows, err := s.pool.Query(ctx, "SELECT event_id, event_type, scheduled, info FROM scheduled_events WHERE scheduled <= CURRENT_TIMESTAMP")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]core.ScheduledEvent, 0)

	for rows.Next() {
		var row core.ScheduledEvent
		var info string
		err = rows.Scan(&row.ID, &row.EventType, &row.Scheduled, &info)
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

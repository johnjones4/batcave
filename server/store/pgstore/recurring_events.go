package pgstore

import (
	"context"
	"encoding/json"
	"main/core"
)

func (s *PGStore) RecurringEvents(ctx context.Context) ([]core.ScheduledRecurringEvent, error) {
	rows, err := s.pool.Query(ctx, "SELECT event_id, source, client_id, intent, scheduled, last_run, info FROM recurring_events")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]core.ScheduledRecurringEvent, 0)

	for rows.Next() {
		var row core.ScheduledRecurringEvent
		var info string
		err = rows.Scan(&row.ID, &row.Source, &row.ClientId, &row.Intent, &row.Scheduled, &row.LastRun, &info)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(info), &row.Info)
		if err != nil {
			return nil, err
		}

		events = append(events, row)
	}

	return events, nil
}

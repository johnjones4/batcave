package types

import "time"

type LogEventRow struct {
	Timestamp time.Time   `json:"timestamp"`
	Event     string      `json:"event"`
	Info      interface{} `json:"info"`
}

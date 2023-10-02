package core

import (
	"context"
	"time"
)

type LLM interface {
	Completion(ctx context.Context, prompt string) (string, error)
	Embedding(ctx context.Context, text string) ([]float32, error)
}

type STT interface {
	SpeechToText(ctx context.Context, wavBytes []byte) (string, error)
}

type ScheduledEventCore struct {
	ID       string
	Source   string
	ClientId string
}

type ScheduledEvent struct {
	ScheduledEventCore
	EventType string
	Scheduled time.Time
	Info      any
}

type ScheduledRecurringEvent struct {
	ScheduledEventCore
	Intent    string
	Scheduled string
	LastRun   time.Time
	Info      map[string]any
}

package core

import (
	"context"
	"time"
)

type PushLogger interface {
	LogPush(ctx context.Context, clientId string, push *PushMessage) error
}

type IntentEmbeddingStore interface {
	UpdateIntentEmbedding(ctx context.Context, intent string, embedding []float32) error
	ClosestMatchingIntent(ctx context.Context, embedding []float32) (string, error)
}

type Scheduler interface {
	ScheduleEvent(ctx context.Context, event *ScheduledEvent) error
	ReadyEvents(ctx context.Context, frontier time.Time, eventType string, infoParser func(event *ScheduledEvent, info string) error) ([]ScheduledEvent, error)
	ClearScheduledEvent(ctx context.Context, id string) error
	ScheduleRecurringEvent(ctx context.Context, event *ScheduledRecurringEvent) error
	ClearScheduledRecurringEvent(ctx context.Context, id string) error
	RecurringEvents(ctx context.Context) ([]ScheduledRecurringEvent, error)
	UpdateRecurringEventTimestamp(ctx context.Context, id string, stamp time.Time) error
}

type ClientRegistry interface {
	UpsertClient(ctx context.Context, source string, clientId string, info any) error
	Client(ctx context.Context, source, clientId string, infoParser func(client *Client, info string) error) (Client, error)
	UserForClient(ctx context.Context, source, clientId string) (User, error)
	ClientsForUser(ctx context.Context, userId string, infoParser func(client *Client, info string) error) ([]Client, error)
}

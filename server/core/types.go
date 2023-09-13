package core

import (
	"context"
	"time"
)

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Message struct {
	Text  string `json:"text"`
	Audio string `json:"audio"`
}

type Request struct {
	EventId    string     `json:"eventId"`
	Message    Message    `json:"message"`
	Source     string     `json:"source"`
	ClientID   string     `json:"clientId"`
	Coordinate Coordinate `json:"coordinate"`
}

type IntentMetadata struct {
	IntentParseCompletion string
	IntentParseReceiver   map[string]any
}

type Media struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}

type Response struct {
	EventId string  `json:"eventId"`
	Message Message `json:"message"`
	Action  string  `json:"action"`
	Media   Media   `json:"media"`
}

type PushMessage struct {
	EventId string  `json:"eventId"`
	Message Message `json:"message"`
	Media   Media   `json:"media"`
}

type PushLogger interface {
	LogPush(ctx context.Context, clientId string, push *PushMessage) error
}

type LLM interface {
	Completion(ctx context.Context, prompt string) (string, error)
	Embedding(ctx context.Context, text string) ([]float32, error)
}

type IntentEmbeddingStore interface {
	UpdateIntentEmbedding(ctx context.Context, intent string, embedding []float32) error
	GetClosestMatchingIntent(ctx context.Context, embedding []float32) (string, error)
}

type RequestProcessor func(ctx context.Context, req *Request) error

type ResponseProcessor func(ctx context.Context, req *Request, res *Response) error

type IntentMatcher interface {
	Match(ctx context.Context, req *Request) (IntentActor, IntentMetadata, error)
}

type IntentActor interface {
	IntentLabel() string
	IntentParsePrompt(req *Request) string
	ActOnIntent(ctx context.Context, req *Request, md *IntentMetadata) (Response, error)
}

type Env struct {
	HttpHost           string `env:"HTTP_HOST"`
	ServiceConfig      string `env:"SERVICE_CONFIG"`
	OpenAIKey          string `env:"OPENAI_API_KEY"`
	IntentDescriptions string `env:"INTENT_DESCRIPTIONS"`
	DatabaseURL        string `env:"DATABASE_URL"`
	TelegramToken      string `env:"TELEGRAM_TOKEN"`
}

type ClientProvider interface {
	SendToClient(ctx context.Context, clientId string, message PushMessage) (bool, error)
}

type ScheduledEvent struct {
	ID        string
	EventType string
	Scheduled time.Time
	Info      any
}

type Scheduler interface {
	ScheduleEvent(ctx context.Context, event *ScheduledEvent) error
	GetReadyEvents(ctx context.Context, eventType string, infoParser func(event *ScheduledEvent, info string) error) ([]ScheduledEvent, error)
	ClearScheduledEvent(ctx context.Context, id string) error
}

type ClientRegistry interface {
	UpsertClient(ctx context.Context, source string, clientId string, info any) error
	GetClient(ctx context.Context, source, clientId string, infoReceiver any) error
}

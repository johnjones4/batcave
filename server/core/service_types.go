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

type TuneIn interface {
	GetStreamURL(query string) (string, error)
}

type Push interface {
	SendLater(ctx context.Context, when time.Time, source string, clientId string, message PushMessage) error
}

type RecurringPush interface {
	SendRecurring(ctx context.Context, source string, clientId string, schedule string, intent string, info map[string]any) error
}

type HomeAssistantGroup struct {
	Names     []string `json:"names"`
	DeviceIds []string `json:"deviceIds"`
	ClientIds []string `json:"clientIds"`
}
type HomeAssistant interface {
	ToggleDeviceState(deviceId string, on bool) error
	Groups() []HomeAssistantGroup
}

type WeatherAlert struct {
	ID            string   `json:"id"`
	AffectedZones []string `json:"affectedZones"`
	Headline      string   `json:"headline"`
}

type WeatherForecastPeriod struct {
	StartTime        time.Time `json:"startTime"`
	EndTime          time.Time `json:"endTime"`
	DetailedForecast string    `json:"detailedForecast"`
	Name             string    `json:"name"`
	Temperature      float64   `json:"temperature"`
	TemperatureUnit  string    `json:"temperatureUnit"`
	WindSpeed        string    `json:"windSpeed"`
	WindDirection    string    `json:"windDirection"`
	Icon             string    `json:"icon"`
	IsDaytime        bool      `json:"isDaytime"`
}

type WeatherForecast struct {
	RadarURL string                  `json:"radarURL"`
	Forecast []WeatherForecastPeriod `json:"forecast"`
	Alerts   []WeatherAlert          `json:"alerts"`
}

type Weather interface {
	PredictWeather(coord Coordinate) (WeatherForecast, error)
}

type Geocoder interface {
	Geocode(q string) (Coordinate, error)
}

package storage

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/johnjones4/hal-9000/server/hal9000/core"
)

type InteractionLogger struct {
	pool *pgxpool.Pool
}

func NewInteractionLogger(pool *pgxpool.Pool) *InteractionLogger {
	return &InteractionLogger{pool}
}

type InteractionEvent struct {
	Request  core.Inbound  `json:"request"`
	Response core.Outbound `json:"response"`
}

func (il *InteractionLogger) Log(e InteractionEvent) error {
	bytes, err := json.Marshal(e)
	if err != nil {
		return err
	}
	_, err = il.pool.Exec(context.Background(), "INSERT INTO log (uuid,timestamp,event) VALUES ($1,$2,$3)", uuid.New().String(), time.Now(), string(bytes))
	return err
}

package storage

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/johnjones4/hal-9000/server/hal9000/core"
)

type InteractionEvent struct {
	Request  core.Inbound  `json:"request"`
	Response core.Outbound `json:"response"`
}

type InteractionLogger interface {
	Log(e InteractionEvent) error
}

type DatabaseInteractionLogger struct {
	pool *pgxpool.Pool
}

func NewDatabaseInteractionLogger(pool *pgxpool.Pool) InteractionLogger {
	return &DatabaseInteractionLogger{pool}
}

func (il *DatabaseInteractionLogger) Log(e InteractionEvent) error {
	bytes, err := json.Marshal(e)
	if err != nil {
		return err
	}
	_, err = il.pool.Exec(context.Background(), "INSERT INTO log (uuid,timestamp,event) VALUES ($1,$2,$3)", uuid.New().String(), time.Now(), string(bytes))
	return err
}

type TerminalInteractionLogger struct{}

func (il *TerminalInteractionLogger) Log(e InteractionEvent) error {
	log.Println(e)
	return nil
}

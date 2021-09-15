package service

import (
	"encoding/json"
	"fmt"
	"hal9000/types"
	"log"
	"time"
)

type loggerConcrete struct {
	log *log.Logger
}

func InitLogger() (types.Logger, error) {
	l := log.Default()
	return loggerConcrete{l}, nil
}

func (l loggerConcrete) LogEvent(event string, info interface{}) {
	eventBytes, err := json.Marshal(types.LogEventRow{
		Timestamp: time.Now(),
		Event:     event,
		Info:      info,
	})
	if err != nil {
		l.LogError(err)
		return
	}
	l.log.Printf("EVENT|%s\n", string(eventBytes))
}

func (l loggerConcrete) LogError(err error) {
	l.log.Printf("ERROR|%s\n", fmt.Sprint(err))
}

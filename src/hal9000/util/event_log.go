package util

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"
)

var file *os.File

type LogEventRow struct {
	Timestamp time.Time   `json:"timestamp"`
	Event     string      `json:"event"`
	Info      interface{} `json:"info"`
}

func InitLogger() error {
	logFile := path.Join(os.Getenv("LOG_DIR"), fmt.Sprintf("%d.txt", int(time.Now().Unix())))
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	file = f
	return nil
}

func LogEvent(event string, info interface{}) {
	eventBytes, err := json.Marshal(LogEventRow{time.Now(), event, info})
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = file.WriteString(string(eventBytes))
	if err != nil {
		fmt.Println(err)
		return
	}
}

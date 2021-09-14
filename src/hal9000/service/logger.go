package service

import (
	"encoding/json"
	"fmt"
	"hal9000/types"
	"os"
	"path"
	"time"
)

type loggerConcrete struct {
	file *os.File
}

func InitLogger() (types.Logger, error) {
	_ = os.Mkdir(os.Getenv("LOG_DIR"), os.ModePerm|os.ModeDir)
	logFile := path.Join(os.Getenv("LOG_DIR"), fmt.Sprintf("%d.txt", int(time.Now().Unix())))
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}
	return loggerConcrete{f}, nil
}

func (l loggerConcrete) LogEvent(event string, info interface{}) {
	eventBytes, err := json.Marshal(types.LogEventRow{
		Timestamp: time.Now(),
		Event:     event,
		Info:      info,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = l.file.WriteString(string(eventBytes) + "\n")
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (l loggerConcrete) LogError(err error) {
	fmt.Println(err)
}

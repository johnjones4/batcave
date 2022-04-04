package learning

import (
	"encoding/json"
	"os"

	"github.com/johnjones4/hal-9000/hal9000/core"
)

type InteractionLogger struct {
	file *os.File
}

func NewInteractionLogger(path string) (*InteractionLogger, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &InteractionLogger{f}, nil
}

type InteractionEvent struct {
	Request  core.Request  `json:"request"`
	Response core.Response `json:"response"`
}

func (il *InteractionLogger) Log(e InteractionEvent) error {
	bytes, err := json.Marshal(e)
	if err != nil {
		return err
	}
	row := append(bytes, '\n')
	_, err = il.file.Write(row)
	return err
}

func (il *InteractionLogger) Close() error {
	return il.file.Close()
}

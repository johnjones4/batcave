package worker

import (
	"context"
	"main/util"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type LogWorker struct {
	workerConcrete
	cfg util.ServerConfig
	msg chan string
}

func NewLogWorker(cfg util.ServerConfig, log logrus.FieldLogger) *LogWorker {
	return &LogWorker{
		cfg:            cfg,
		workerConcrete: newWorkerConcrete(log),
		msg:            make(chan string, 255),
	}
}

func (w *LogWorker) Start(ctx context.Context) {
	conn, _, err := websocket.DefaultDialer.Dial(w.cfg.URL("/api/client/log"), w.cfg.Headers())
	if err != nil {
		w.log.Errorf("error connecting to websocket: %s", err)
		return
	}
	defer func() {
		w.stopped <- true
		conn.Close()
	}()

	msgs := make(chan string)
	go func() {
		msgT, msg, err := conn.ReadMessage()
		if err != nil {
			w.log.Errorf("error reading websocket: %s", err)
			return
		}
		if msgT == websocket.TextMessage {
			msgs <- string(msg)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-w.stop:
			return
		case msg := <-msgs:
			w.msg <- msg
		}
	}
}

func (w *LogWorker) Stop() error {
	w.stop <- true
	<-w.stopped
	return nil
}

func (w *LogWorker) Chan() chan string {
	return w.msg
}

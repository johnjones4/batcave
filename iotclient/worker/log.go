package worker

import (
	"context"
	"main/util"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type LogWorker struct {
	workerConcrete
	socketworker
	msg chan string
}

func NewLogWorker(cfg util.ServerConfig, log logrus.FieldLogger) *LogWorker {
	return &LogWorker{
		workerConcrete: newWorkerConcrete(log),
		msg:            make(chan string, 255),
		socketworker: socketworker{
			url: "/api/client/log",
			cfg: cfg,
			log: log,
		},
	}
}

func (w *LogWorker) Start(ctx context.Context) {
	var conn *websocket.Conn
	conn = w.reconnect()
	if conn == nil {
		return
	}

	defer func() {
		conn.Close()
		w.stopped <- true
	}()

	msgs := make(chan string)
	go func() {
		msgT, msg, err := conn.ReadMessage()
		if err != nil {
			w.workerConcrete.log.Errorf("error reading websocket: %s", err)

			conn.Close()
			conn = w.reconnect()
			if conn == nil {
				return
			}
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

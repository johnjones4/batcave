package worker

import (
	"context"
	"main/core"
	"main/util"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type CommandWorker struct {
	workerConcrete
	socketworker
	sendQueue    chan core.Request
	receiveQueue chan core.Response
}

func NewCommandWorker(cfg util.ServerConfig, log logrus.FieldLogger) *CommandWorker {
	return &CommandWorker{
		workerConcrete: newWorkerConcrete(log),
		sendQueue:      make(chan core.Request, 32),
		receiveQueue:   make(chan core.Response, 32),
		socketworker: socketworker{
			url: "/api/client/converse",
			cfg: cfg,
			log: log,
		},
	}
}

func (w *CommandWorker) Setup() error {
	return nil
}

func (w *CommandWorker) Teardown() error {
	return nil
}

func (w *CommandWorker) Start(ctx context.Context) {
	var conn *websocket.Conn
	conn = w.socketworker.reconnect()
	if conn == nil {
		return
	}

	defer func() {
		conn.Close()
		w.stopped <- true
	}()

	needsReconnect := make(chan bool)
	reconnecting := make(chan bool)

	go func() {
		for {
			var response core.Response
			err := conn.ReadJSON(&response)
			if err != nil {
				w.workerConcrete.log.Errorf("error reading websocket: %s", err)
				needsReconnect <- true
				<-reconnecting
			}
			w.receiveQueue <- response
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-w.stop:
			return
		case <-needsReconnect:
			conn.Close()
			conn = w.socketworker.reconnect()
			if conn == nil {
				return
			}
			reconnecting <- true
		case request := <-w.sendQueue:
			err := conn.WriteJSON(request)
			if err != nil {
				w.workerConcrete.log.Errorf("error writing websocket: %s", err)
				conn.Close()
				conn = w.socketworker.reconnect()
				if conn == nil {
					return
				}
			}
		}
	}
}

func (w *CommandWorker) Stop() error {
	w.stop <- true
	<-w.stopped
	return nil
}

func (w *CommandWorker) Chan() chan core.Response {
	return w.receiveQueue
}

func (w *CommandWorker) SendChan() chan core.Request {
	return w.sendQueue
}

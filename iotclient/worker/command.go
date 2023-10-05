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
	cfg          util.ServerConfig
	sendQueue    chan core.Request
	receiveQueue chan core.Response
}

func NewCommandWorker(cfg util.ServerConfig, log logrus.FieldLogger) *CommandWorker {
	return &CommandWorker{
		workerConcrete: newWorkerConcrete(log),
		cfg:            cfg,
		sendQueue:      make(chan core.Request, 32),
		receiveQueue:   make(chan core.Response, 32),
	}
}

func (w *CommandWorker) Setup() error {
	return nil
}

func (w *CommandWorker) Teardown() error {
	return nil
}

func (w *CommandWorker) Start(ctx context.Context) {
	conn, _, err := websocket.DefaultDialer.Dial(w.cfg.URL("/api/client/converse"), w.cfg.Headers())
	if err != nil {
		w.log.Errorf("error connecting to websocket: %s", err)
		return
	}
	defer conn.Close()
	defer func() {
		w.stopped <- true
	}()

	go func() {
		for {
			var response core.Response
			err := conn.ReadJSON(&response)
			if err != nil {
				w.log.Errorf("error reading websocket: %s", err)
				return
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
		case request := <-w.sendQueue:
			err = conn.WriteJSON(request)
			if err != nil {
				w.log.Errorf("error writing websocket: %s", err)
				return
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

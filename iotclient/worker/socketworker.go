package worker

import (
	"main/util"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type socketworker struct {
	url string
	cfg util.ServerConfig
	log logrus.FieldLogger
}

func (w *socketworker) reconnect() *websocket.Conn {
	for i := 0; i < 10; i++ {
		w.log.Debugf("connecting to %s attempt %d", w.url, i)
		conn, _, err := websocket.DefaultDialer.Dial(w.cfg.URL(w.url), w.cfg.Headers())
		if err != nil {
			w.log.Errorf("error connecting to websocket: %s", err)
			time.Sleep(time.Second * time.Duration(i))
		} else {
			w.log.Debug("connection successful")
			return conn
		}
	}
	return nil
}

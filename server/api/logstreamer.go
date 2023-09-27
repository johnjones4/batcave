package api

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	maxWait = time.Second * 10
	minWait = time.Millisecond * 100
)

func (a *API) streamer(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		a.handleError(w, err, http.StatusBadRequest)
		return
	}
	defer c.Close()

	err = c.WriteMessage(websocket.TextMessage, []byte("Hello"))
	if err != nil {
		a.Log.Error(err)
		return
	}

	var lastStamp time.Time
	var wait = minWait
	for {
		var msg string
		a.logMsgLock.RLock()
		if a.logMsgStamp.After(lastStamp) {
			msg = a.logMsg
			lastStamp = a.logMsgStamp
		}
		a.logMsgLock.RUnlock()

		if msg == "" {
			time.Sleep(wait)
			wait = wait * 2
			if wait > maxWait {
				wait = maxWait
			}
			continue
		}

		err = c.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			a.Log.Error(err)
			return
		}
	}
}

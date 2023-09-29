package api

import (
	"net/http"

	"github.com/gorilla/websocket"
)

func (a *API) addLogListener() *logListener {
	a.logListenersLock.Lock()
	defer a.logListenersLock.Unlock()
	l := &logListener{
		channel: make(chan string, 255),
	}
	if a.logListeners == nil {
		a.logListeners = l
	} else {
		tail := a.logListeners
		for tail != l {
			if tail.next == nil {
				tail.next = l
			}
			tail = tail.next
		}
	}
	return l
}

func (a *API) _removeLogListener(tail *logListener, l *logListener) {
	if tail.next == l {
		tail.next = tail.next.next
	} else if tail.next != nil {
		a._removeLogListener(tail.next, l)
	}
}

func (a *API) removeLogListener(l *logListener) {
	a.logListenersLock.Lock()
	defer a.logListenersLock.Unlock()
	if a.logListeners == l {
		a.logListeners = l.next
	} else {
		a._removeLogListener(a.logListeners, l)
	}
}

func (a *API) streamer(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		a.handleError(w, err, http.StatusBadRequest)
		return
	}
	defer c.Close()

	listener := a.addLogListener()
	defer a.removeLogListener(listener)

	err = c.WriteMessage(websocket.TextMessage, []byte("Hello"))
	if err != nil {
		a.Log.Error(err)
		return
	}

	for msg := range listener.channel {
		err = c.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			a.Log.Error(err)
			return
		}
	}
}

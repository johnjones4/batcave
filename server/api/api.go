package api

import (
	"context"
	"main/core"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{}

type APIParams struct {
	IntentMatcher      core.IntentMatcher
	RequestProcessors  []core.RequestProcessor
	ResponseProcessors []core.ResponseProcessor
	Log                core.HookableLogger
	Telegram           core.TelegramSender
	ClientRegistry     core.ClientRegistry
	SocketSender       core.SocketSender
}

type logListener struct {
	next    *logListener
	channel chan string
}

type API struct {
	logMessages      chan string
	mux              *chi.Mux
	logListenersLock sync.Mutex
	logListeners     *logListener
	APIParams
}

func (a *API) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	a.mux.ServeHTTP(res, req)
}

func (a *API) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-a.logMessages:
			a.logListenersLock.Lock()
			l := a.logListeners
			for l != nil {
				l.channel <- msg
				l = l.next
			}
			a.logListenersLock.Unlock()
		}
	}
}

func (a *API) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (a *API) Fire(entry *logrus.Entry) error {
	str, err := entry.String()
	if err != nil {
		return err
	}
	go func() {
		a.logMessages <- str
	}()
	return nil
}

func New(params APIParams) *API {
	a := &API{
		mux:       chi.NewRouter(),
		APIParams: params,
	}

	a.mux.Use(middleware.RequestID)
	a.mux.Use(middleware.RealIP)
	a.mux.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: params.Log}))
	a.mux.Use(middleware.Recoverer)

	a.mux.Route("/api", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
		})

		r.Route("/client", func(r chi.Router) {
			r.Use(a.authMiddleware)
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			})
			r.Post("/message", a.message)
			r.Handle("/log", http.HandlerFunc(a.streamer))
			r.Handle("/converse", http.HandlerFunc(a.converse))
		})

		r.Post("/telegram", a.telegramHandler)
	})

	a.logMessages = make(chan string)
	a.Log.AddHook(a)

	return a
}

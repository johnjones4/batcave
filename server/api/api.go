package api

import (
	"main/core"
	"main/services/telegram"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

type APIParams struct {
	IntentMatcher      core.IntentMatcher
	RequestProcessors  []core.RequestProcessor
	ResponseProcessors []core.ResponseProcessor
	Log                logrus.FieldLogger
	Telegram           telegram.Telegram
}

type apiConcrete struct {
	mux *chi.Mux
	APIParams
}

func (a *apiConcrete) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	a.mux.ServeHTTP(res, req)
}

func New(params APIParams) http.Handler {
	a := apiConcrete{
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

		r.Post("/message", a.directHandler)
		r.Post("/telegram", a.telegramHandler)
	})

	return &a
}

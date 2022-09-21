package api

import (
	"net/http"

	"github.com/johnjones4/hal-9000/server/hal9000/runtime"
	"github.com/swaggest/rest/web"
)

type API struct {
	Host    string
	Runtime *runtime.Runtime
}

func (a *API) Run() error {
	s := web.DefaultService()

	// Init API documentation schema.
	s.OpenAPI.Info.Title = "Basic Example"
	s.OpenAPI.Info.WithDescription("This app showcases a trivial REST API.")
	s.OpenAPI.Info.Version = "v1.2.3"

	// Setup middlewares.
	s.Use(
		logRequest,
		// makeRequestVerifier(a.Runtime.ClientStore),
	)

	s.Post("/api/request", makeRequestHandler(a.Runtime))
	s.Get("/api/commands", makeCommandsHandler(a.Runtime))
	s.Get("/api/ping", makePingHandler(a.Runtime))
	s.Get("/api/info", makeInfoHandler(a.Runtime))

	return http.ListenAndServe(a.Host, s)
}

package api

import (
	"context"
	"net/http"

	"github.com/johnjones4/hal-9000/server/hal9000/core"
	"github.com/johnjones4/hal-9000/server/hal9000/runtime"
	"github.com/johnjones4/hal-9000/server/hal9000/security"

	"github.com/swaggest/usecase"
)

func makeCommandsHandler(r *runtime.Runtime) usecase.Interactor {
	type request struct {
		Authorization string `header:"Authorization"`
	}
	type response struct {
		Commands map[string]string `json:"commands"`
	}
	return usecase.NewIOI(new(request), new(response), func(ctx context.Context, input, output interface{}) error {
		in := input.(*request)
		out := output.(*response)

		username, err := r.TokenManager.UsernameForToken(in.Authorization)
		if err != nil {
			if err == security.ErrorTokenExpired {
				return apiError{
					status: http.StatusForbidden,
					code:   core.ErrorCodeExpiredToken,
					parent: err,
				}
			}
			return wrappedError(err, core.ErrorCodeStore)
		}

		user, err := r.UserStore.GetUser(username)
		if err != nil {
			return wrappedError(err, core.ErrorCodeStore)
		}

		state, err := r.StateStore.GetStateForUser(user.User)
		if err != nil {
			return wrappedError(err, core.ErrorCodeStore)
		}

		for _, intent := range r.Intents.Intents {
			for command, description := range intent.SupportedComandsForState(state) {
				out.Commands[command] = description
			}
		}

		return nil
	})
}

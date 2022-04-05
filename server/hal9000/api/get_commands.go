package api

import (
	"context"
	"net/http"

	"github.com/johnjones4/hal-9000/server/hal9000/core"
	"github.com/johnjones4/hal-9000/server/hal9000/intent"
	"github.com/johnjones4/hal-9000/server/hal9000/security"
	"github.com/johnjones4/hal-9000/server/hal9000/storage"

	"github.com/swaggest/usecase"
)

func makeCommandsHandler(intentSet *intent.IntentSet, userStore *storage.UserStore, stateStore *storage.StateStore, tm *security.TokenManager) usecase.Interactor {
	type request struct {
		Authorization string `header:"Authorization"`
	}
	type response struct {
		Commands map[string]string `json:"commands"`
	}
	return usecase.NewIOI(new(request), new(response), func(ctx context.Context, input, output interface{}) error {
		in := input.(*request)
		out := output.(*response)

		username, err := tm.UsernameForToken(in.Authorization)
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

		user, err := userStore.GetUser(username)
		if err != nil {
			return wrappedError(err, core.ErrorCodeStore)
		}

		state, err := stateStore.GetStateForUser(user.User)
		if err != nil {
			return wrappedError(err, core.ErrorCodeStore)
		}

		for _, intent := range intentSet.Intents {
			for command, description := range intent.SupportedComandsForState(state) {
				out.Commands[command] = description
			}
		}

		return nil
	})
}

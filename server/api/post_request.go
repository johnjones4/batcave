package api

import (
	"context"
	"main/core"
	"main/intent"
	"main/learning"
	"main/storage"
	"net/http"

	"github.com/swaggest/usecase"
)

func makeRequestHandler(intentSet *intent.IntentSet, userStore *storage.UserStore, stateStore *storage.StateStore, logger *learning.InteractionLogger, tm *TokenManager) usecase.Interactor {
	type request struct {
		core.RequestBody
		Authorization string `header:"Authorization"`
	}
	return usecase.NewIOI(new(request), new(core.ResponseBody), func(ctx context.Context, input, output interface{}) error {
		in := input.(*request)
		username, err := tm.usernameForToken(in.Authorization)
		if err != nil {
			if err == errorTokenExpired {
				return apiError{
					status: http.StatusForbidden,
					code:   errorCodeExpiredToken,
					parent: err,
				}
			}
			return wrappedError(err, errorCodeInternal)
		}

		user, err := userStore.GetUser(username)
		if err != nil {
			return wrappedError(err, errorCodeStore)
		}

		state, err := stateStore.GetStateForUser(user.User)
		if err != nil {
			return wrappedError(err, errorCodeStore)
		}

		request := core.Request{
			RequestBody: in.RequestBody,
			State:       state,
		}

		response, err := intentSet.ProcessRequest(request)
		if err != nil {
			return wrappedError(err, errorCodeReqestProcess)
		}

		err = logger.Log(learning.InteractionEvent{
			Request:  request,
			Response: response,
		})
		if err != nil {
			return wrappedError(err, errorCodeLog)
		}

		err = stateStore.SetStateForUSer(request.State)
		if err != nil {
			return wrappedError(err, errorCodeStore)
		}

		out := output.(*core.ResponseBody)
		*out = response.ResponseBody

		return nil
	})
}

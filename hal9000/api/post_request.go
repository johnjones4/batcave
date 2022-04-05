package api

import (
	"context"
	"net/http"

	"github.com/johnjones4/hal-9000/hal9000/core"
	"github.com/johnjones4/hal-9000/hal9000/intent"
	"github.com/johnjones4/hal-9000/hal9000/learning"
	"github.com/johnjones4/hal-9000/hal9000/security"
	"github.com/johnjones4/hal-9000/hal9000/storage"

	"github.com/swaggest/usecase"
)

func makeRequestHandler(intentSet *intent.IntentSet, userStore *storage.UserStore, stateStore *storage.StateStore, logger *learning.InteractionLogger, tm *security.TokenManager) usecase.Interactor {
	type request struct {
		core.InboundBody
		Authorization string `header:"Authorization"`
	}
	return usecase.NewIOI(new(request), new(core.OutboundBody), func(ctx context.Context, input, output interface{}) error {
		in := input.(*request)
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

		request, err := intent.Parse(in.InboundBody, state)
		if err != nil {
			return err //TODO
		}

		response, err := intentSet.ProcessRequest(request)
		if err != nil {
			return wrappedError(err, core.ErrorCodeReqestProcess)
		}

		err = logger.Log(learning.InteractionEvent{
			Request:  request,
			Response: response,
		})
		if err != nil {
			return wrappedError(err, core.ErrorCodeLog)
		}

		err = stateStore.SetStateForUSer(request.State)
		if err != nil {
			return wrappedError(err, core.ErrorCodeStore)
		}

		out := output.(*core.OutboundBody)
		*out = response.OutboundBody

		return nil
	})
}

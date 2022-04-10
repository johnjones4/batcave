package api

import (
	"context"
	"net/http"

	"github.com/johnjones4/hal-9000/server/hal9000/core"
	"github.com/johnjones4/hal-9000/server/hal9000/learning"
	"github.com/johnjones4/hal-9000/server/hal9000/runtime"
	"github.com/johnjones4/hal-9000/server/hal9000/security"

	"github.com/swaggest/usecase"
)

func makeRequestHandler(r *runtime.Runtime) usecase.Interactor {
	type request struct {
		core.InboundBody
		Authorization string `header:"Authorization"`
	}
	return usecase.NewIOI(new(request), new(core.OutboundBody), func(ctx context.Context, input, output interface{}) error {
		in := input.(*request)
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

		request, err := r.Parse(in.InboundBody, state)
		if err != nil {
			return wrappedError(err, core.ErrorCodeParsing)
		}

		response, err := r.Intents.ProcessRequest(request)
		if err != nil {
			return wrappedError(err, core.ErrorCodeReqestProcess)
		}

		err = r.Logger.Log(learning.InteractionEvent{
			Request:  request,
			Response: response,
		})
		if err != nil {
			return wrappedError(err, core.ErrorCodeLog)
		}

		err = r.StateStore.SetStateForUser(request.State)
		if err != nil {
			return wrappedError(err, core.ErrorCodeStore)
		}

		out := output.(*core.OutboundBody)
		*out = response.OutboundBody

		return nil
	})
}

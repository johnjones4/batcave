package api

import (
	"context"

	"github.com/johnjones4/hal-9000/server/hal9000/core"
	"github.com/johnjones4/hal-9000/server/hal9000/runtime"
	"github.com/johnjones4/hal-9000/server/hal9000/storage"

	"github.com/swaggest/usecase"
)

func makeRequestHandler(r *runtime.Runtime) usecase.Interactor {
	type request struct {
		core.InboundBody
		Client string `header:"User-Agent"`
	}
	return usecase.NewIOI(new(request), new(core.OutboundBody), func(ctx context.Context, input, output interface{}) error {
		in := input.(*request)

		client, err := r.ClientStore.GetClient(in.Client)
		if err != nil {
			return wrappedError(err, core.ErrorCodeClient)
		}

		state, err := r.StateStore.GetState(client.Client)
		if err != nil {
			return wrappedError(err, core.ErrorCodeStore)
		}

		request, err := r.Parse(in.InboundBody, client.Client, state)
		if err != nil {
			return wrappedError(err, core.ErrorCodeParsing)
		}

		response, err := r.Intents.ProcessRequest(request)
		if err != nil {
			return wrappedError(err, core.ErrorCodeReqestProcess)
		}

		err = r.Logger.Log(storage.InteractionEvent{
			Request:  request,
			Response: response,
		})
		if err != nil {
			return wrappedError(err, core.ErrorCodeLog)
		}

		err = r.StateStore.SetState(client.Client, response.State)
		if err != nil {
			return wrappedError(err, core.ErrorCodeStore)
		}

		out := output.(*core.OutboundBody)
		*out = response.OutboundBody

		return nil
	})
}

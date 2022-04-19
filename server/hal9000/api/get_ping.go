package api

import (
	"context"

	"github.com/johnjones4/hal-9000/server/hal9000/core"
	"github.com/johnjones4/hal-9000/server/hal9000/runtime"

	"github.com/swaggest/usecase"
)

func makePingHandler(r *runtime.Runtime) usecase.Interactor {
	type request struct {
		Client string `header:"User-Agent"`
	}
	type response struct {
		Pong bool `json:"pong"`
	}
	return usecase.NewIOI(new(request), new(response), func(ctx context.Context, input, output interface{}) error {
		in := input.(*request)
		out := output.(*response)

		client, err := r.ClientStore.GetClient(in.Client)
		if err != nil {
			return wrappedError(err, core.ErrorCodeClient)
		}

		_, err = r.StateStore.GetState(client.Client)
		if err != nil {
			return wrappedError(err, core.ErrorCodeStore)
		}

		out.Pong = true

		return nil
	})
}

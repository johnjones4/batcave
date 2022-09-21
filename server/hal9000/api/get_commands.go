package api

import (
	"context"

	"github.com/johnjones4/hal-9000/server/hal9000/core"
	"github.com/johnjones4/hal-9000/server/hal9000/runtime"

	"github.com/swaggest/usecase"
)

func makeCommandsHandler(r *runtime.Runtime) usecase.Interactor {
	type request struct {
		Client string `header:"User-Agent"`
	}
	type response struct {
		Commands map[string]core.CommandInfo `json:"commands"`
	}
	return usecase.NewIOI(new(request), new(response), func(ctx context.Context, input, output interface{}) error {
		in := input.(*request)
		out := output.(*response)

		client, err := r.ClientStore.GetClient(in.Client)
		if err != nil {
			return wrappedError(err, core.ErrorCodeClient)
		}

		state, err := r.StateStore.GetState(client.Client)
		if err != nil {
			return wrappedError(err, core.ErrorCodeStore)
		}

		out.Commands = make(map[string]core.CommandInfo)
		for _, intent := range r.Intents.Intents {
			for command, description := range intent.SupportedCommandsForState(state) {
				out.Commands[command] = description
			}
		}

		return nil
	})
}

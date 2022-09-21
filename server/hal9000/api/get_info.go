package api

import (
	"context"

	"github.com/johnjones4/hal-9000/server/hal9000/core"
	"github.com/johnjones4/hal-9000/server/hal9000/runtime"
	"github.com/johnjones4/hal-9000/server/hal9000/service"

	"github.com/swaggest/usecase"
)

func makeInfoHandler(r *runtime.Runtime) usecase.Interactor {
	type request struct {
		Latitude  float64 `query:"latitude"`
		Longitude float64 `query:"longitude"`
	}
	type response struct {
		Info map[string]interface{} `json:"info"`
	}
	return usecase.NewIOI(new(request), new(response), func(ctx context.Context, input, output interface{}) error {
		in := input.(*request)
		out := output.(*response)

		ctx1 := core.ContextWithCoordinates(ctx, core.Coordinate{
			Latitude:  in.Latitude,
			Longitude: in.Longitude,
		})

		out.Info = make(map[string]interface{}, 0)
		for _, s := range r.Intents.Services() {
			if infoService, ok := s.(service.InfoService); ok {
				name := infoService.Name()
				info, err := infoService.Info(ctx1)
				if err != nil {
					return err
				}
				out.Info[name] = info
			}
		}

		return nil
	})
}

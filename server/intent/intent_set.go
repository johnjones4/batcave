package intent

import (
	"fmt"
	"main/core"
)

type IntentSet struct {
	Intents []core.Intent
}

func (h *IntentSet) ProcessRequest(req core.Request) (core.Response, error) {
	for _, commandHandler := range h.Intents {
		for _, command := range commandHandler.SupportedComandsForState(req.State) {
			if req.Command == command {
				res, err := commandHandler.Execute(req)
				if err != nil {
					if ferr, ok := err.(core.FeedbackError); ok {
						return core.Response{
							ResponseBody: core.ResponseBody{
								Message: ferr.Error(),
							},
							State: req.State,
						}, nil
					}
					return core.Response{}, err
				}
				return res, nil
			}
		}
	}
	return core.Response{}, fmt.Errorf("no handler for %s", fmt.Sprint(req))
}

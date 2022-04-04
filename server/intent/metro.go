package intent

import (
	"fmt"
	"main/core"
	"main/service"
	"strings"
)

type Metro struct {
	Service *service.Metro
}

func (c *Metro) SupportedComandsForState(s core.State) []string {
	if s.State != core.StateDefault {
		return []string{}
	}
	return []string{
		"metro",
	}
}

func (c *Metro) Execute(req core.Request) (core.Response, error) {
	info, err := c.Service.GetArrivals(req.Location)
	if err != nil {
		return core.Response{}, err
	}

	arrivals := make([]string, len(info))
	for i, arrival := range info {
		arrivals[i] = fmt.Sprintf("%s -> %s (%s)", arrival.Line, arrival.Destination, arrival.Min)
	}

	return core.Response{
		ResponseBody: core.ResponseBody{
			Message: "Here are the upcoming arrivals near you:\n" + strings.Join(arrivals, "\n"),
		},
		State: req.State,
	}, nil
}

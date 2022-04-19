package intent

import (
	"fmt"
	"strings"

	"github.com/johnjones4/hal-9000/server/hal9000/core"
	"github.com/johnjones4/hal-9000/server/hal9000/service"
)

type Metro struct {
	Service *service.Metro
}

func (c *Metro) SupportedComandsForState(s string) map[string]core.CommandInfo {
	if s != core.StateDefault {
		return map[string]core.CommandInfo{}
	}
	return map[string]core.CommandInfo{
		"metro": {
			Description:  "Get the metro arrivals for the closest station.",
			RequiresBody: false,
		},
	}
}

func (c *Metro) Execute(req core.Inbound) (core.Outbound, error) {
	station, info, err := c.Service.GetArrivals(req.Location)
	if err != nil {
		return core.Outbound{}, err
	}

	arrivals := make([]string, len(info))
	for i, arrival := range info {
		arrivals[i] = fmt.Sprintf("%s %s (%s)", arrival.Line, arrival.Destination, arrival.Min)
	}

	return core.Outbound{
		OutboundBody: core.OutboundBody{
			Body: fmt.Sprintf("Upcoming arrivals for %s:\n%s", station.Name, strings.Join(arrivals, "\n")),
		},
		State: req.State,
	}, nil
}

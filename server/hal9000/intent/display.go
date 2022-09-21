package intent

import (
	"context"

	"github.com/johnjones4/hal-9000/server/hal9000/core"
	"github.com/johnjones4/hal-9000/server/hal9000/service"
	"github.com/johnjones4/hal-9000/server/hal9000/util"
)

type Display struct {
	DisplayServices []service.DisplayService
}

func (c *Display) Services() []core.Service {
	services := make([]core.Service, 0)
	for _, s := range c.DisplayServices {
		services = append(services, s)
	}
	return services
}

func (c *Display) SupportedCommandsForState(s string) map[string]core.CommandInfo {
	return map[string]core.CommandInfo{
		"display": {
			Description:  "Show a given display device's output",
			RequiresBody: true,
		},
	}
}

func (c *Display) Execute(req core.Inbound) (core.Outbound, error) {
	displayMap := make(map[string]service.Displayable)
	displays := make([]string, 0)
	for _, displayProvider := range c.DisplayServices {
		displaysList := displayProvider.Displays()
		for _, display := range displaysList {
			for _, name := range display.Names() {
				displays = append(displays, name)
				displayMap[name] = display
			}
		}
	}

	displayName := util.FindClosestMatchString(displays, req.Body)
	if displayName == "" {
		return core.Outbound{}, core.NewFeedbackError("Could not find anything to display")
	}

	display := displayMap[displayName]

	ctx := core.ContextWithCoordinates(context.Background(), req.Location)

	url, err := display.URL(ctx)
	if err != nil {
		return core.Outbound{}, err
	}

	return core.Outbound{
		OutboundBody: core.OutboundBody{
			URL: url,
		},
		State: req.State,
	}, nil
}

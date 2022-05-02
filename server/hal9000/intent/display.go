package intent

import (
	"github.com/johnjones4/hal-9000/server/hal9000/core"
	"github.com/johnjones4/hal-9000/server/hal9000/service"
	"github.com/johnjones4/hal-9000/server/hal9000/util"
)

type Display struct {
	Services []service.DisplayService
}

func (c *Display) SupportedComandsForState(s string) map[string]core.CommandInfo {
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
	for _, displayProvider := range c.Services {
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

	url, err := display.URL(req)
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

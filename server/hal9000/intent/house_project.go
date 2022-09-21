package intent

import (
	"fmt"

	"github.com/johnjones4/hal-9000/server/hal9000/core"
	"github.com/johnjones4/hal-9000/server/hal9000/service"
)

type HouseProjectConfiguration struct {
	ListId string
}
type HouseProject struct {
	Service       *service.Trello
	Configuration HouseProjectConfiguration
}

func (c *HouseProject) Services() []core.Service {
	services := make([]core.Service, 1)
	services[0] = c.Service
	return services
}

func (c *HouseProject) SupportedCommandsForState(s string) map[string]core.CommandInfo {
	if s != core.StateDefault {
		return map[string]core.CommandInfo{}
	}
	return map[string]core.CommandInfo{
		"house-project-add": {
			Description:  "Add a new house project",
			RequiresBody: true,
		},
	}
}

func (c *HouseProject) Execute(req core.Inbound) (core.Outbound, error) {
	name := req.Body

	if name == "" {
		return core.Outbound{}, core.NewFeedbackError("Please provide a project name")
	}

	url, err := c.Service.NewCard(c.Configuration.ListId, name)
	if err != nil {
		return core.Outbound{}, err
	}
	return core.Outbound{
		OutboundBody: core.OutboundBody{
			Body: fmt.Sprintf("Added \"%s\" to your projects list", name),
			URL:  url,
		},
		State: req.State,
	}, nil
}

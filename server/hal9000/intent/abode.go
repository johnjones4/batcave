package intent

import (
	"fmt"
	"strings"

	"github.com/johnjones4/hal-9000/server/hal9000/core"
	"github.com/johnjones4/hal-9000/server/hal9000/service"
	"github.com/johnjones4/hal-9000/server/hal9000/util"
)

type Abode struct {
	Service *service.Abode
}

const (
	AbodeCommandMode = "abode-mode"
	AbodeCommandInfo = "abode-info"
)

func (c *Abode) SupportedComandsForState(s string) map[string]core.CommandInfo {
	if s != core.StateDefault {
		return map[string]core.CommandInfo{}
	}
	return map[string]core.CommandInfo{
		AbodeCommandMode: {
			Description:  "Set Abode system mode.",
			RequiresBody: true,
		},
		AbodeCommandInfo: {
			Description:  "Get Abode devices status.",
			RequiresBody: false,
		},
	}
}

func (c *Abode) Execute(req core.Inbound) (core.Outbound, error) {
	switch req.Command {
	case AbodeCommandMode:
		mode := req.Body
		if !util.ArrayContains([]string{service.AbodeModeAway, service.AbodeModeHome, service.AbodeModeStandby}, mode) {
			return core.Outbound{}, core.NewFeedbackError(fmt.Sprintf("Unsupported mode: %s", mode))
		}

		err := c.Service.SetMode(mode)
		if err != nil {
			return core.Outbound{}, err
		}

		return core.Outbound{
			OutboundBody: core.OutboundBody{
				Body: fmt.Sprintf("Abode has been set to \"%s\"", mode),
			},
			State: req.State,
		}, nil
	case AbodeCommandInfo:
		panel, err := c.Service.GetPanel()
		if err != nil {
			return core.Outbound{}, err
		}

		statuses, err := c.Service.GetDeviceStatuses()
		if err != nil {
			return core.Outbound{}, err
		}

		statusStr := strings.Builder{}

		statusStr.WriteString(fmt.Sprintf("System: %s\n", panel.Mode.Label))

		for i, status := range statuses {
			statusStr.WriteString(fmt.Sprintf("%s: %s", status.Name, status.Status))
			if i < len(statuses)-1 {
				statusStr.WriteByte('\n')
			}
		}

		return core.Outbound{
			OutboundBody: core.OutboundBody{
				Body: "Abode status:\n" + statusStr.String(),
			},
			State: req.State,
		}, nil
	default:
		return core.Outbound{}, fmt.Errorf("no handler for %s", fmt.Sprint(req))
	}
}

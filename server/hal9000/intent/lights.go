package intent

import (
	"fmt"
	"strings"

	"github.com/johnjones4/hal-9000/server/hal9000/core"
	"github.com/johnjones4/hal-9000/server/hal9000/service"
	"github.com/johnjones4/hal-9000/server/hal9000/util"
)

const (
	lightsOn  = "on"
	lightsOff = "offs"
)

type Lights struct {
	Service *service.Kasa
}

func (c *Lights) SupportedComandsForState(s string) map[string]core.CommandInfo {
	if s != core.StateDefault {
		return map[string]core.CommandInfo{}
	}
	return map[string]core.CommandInfo{
		"lights": {
			Description:  "Turn the given lights on and off.",
			RequiresBody: true,
		},
	}
}

func (c *Lights) Execute(req core.Inbound) (core.Outbound, error) {
	names, mapped := c.Service.DeviceNamesAndMap()
	name := util.FindClosestMatchString(names, req.Body)
	if name == "" {
		return core.Outbound{}, core.NewFeedbackError("Could not find a matching device")
	}
	deviceGroup := mapped[name]

	newState := util.FindClosestMatchString([]string{lightsOn, lightsOff}, req.Body)
	if newState == "" {
		return core.Outbound{}, core.NewFeedbackError("Could not determine requested state")
	}

	for _, device := range deviceGroup.Devices {
		err := c.Service.SetStatus(device, newState == lightsOn)
		if err != nil {
			return core.Outbound{}, err
		}
	}

	return core.Outbound{
		OutboundBody: core.OutboundBody{
			Body: fmt.Sprintf("\"%s\" is now %s.", deviceGroup.PreferredName, newState),
		},
		State: req.State,
	}, nil
}

func determineStateRequest(body string) bool {
	onWords := []string{"on"}
	lc := strings.ToLower(body)
	for _, word := range onWords {
		if strings.Contains(lc, word) {
			return true
		}
	}
	return false
}

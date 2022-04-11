package intent

import (
	"fmt"
	"math"
	"strings"

	"github.com/johnjones4/hal-9000/server/hal9000/core"
	"github.com/johnjones4/hal-9000/server/hal9000/service"
	"github.com/johnjones4/hal-9000/server/hal9000/util"
)

type Lights struct {
	Service *service.Kasa
}

func (c *Lights) SupportedComandsForState(s core.State) map[string]core.CommandInfo {
	if s.State != core.StateDefault {
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
	deviceGroup, err := findDeviceInCommand(c.Service.DeviceGroups(), req.Body)
	if err != nil {
		return core.Outbound{}, err
	}

	newState := determineStateRequest(req.Body)

	for _, device := range deviceGroup.Devices {
		err = c.Service.SetStatus(device, newState)
		if err != nil {
			return core.Outbound{}, err
		}
	}

	newStateStr := "on"
	if !newState {
		newStateStr = "off"
	}

	return core.Outbound{
		OutboundBody: core.OutboundBody{
			Body: fmt.Sprintf("\"%s\" is now %s.", deviceGroup.PreferredName, newStateStr),
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

func findDeviceInCommand(deviceGroups []service.KasaDeviceGroup, body string) (service.KasaDeviceGroup, error) {
	bestDistance := math.MaxInt
	var best service.KasaDeviceGroup
	for _, deviceGroup := range deviceGroups {
		for _, name := range deviceGroup.Names {
			dist := util.Levenshtein([]rune(body), []rune(name))
			if dist < bestDistance {
				bestDistance = dist
				best = deviceGroup
			}
		}
	}
	if best.PreferredName == "" {
		return service.KasaDeviceGroup{}, core.NewFeedbackError("could not find a matching device")
	}

	return best, nil
}

package intent

import (
	"fmt"
	"math"
	"strings"

	"github.com/johnjones4/hal-9000/hal9000/core"
	"github.com/johnjones4/hal-9000/hal9000/service"
	"github.com/johnjones4/hal-9000/hal9000/util"
)

type Lights struct {
	Service *service.Kasa
}

func (c *Lights) SupportedComandsForState(s core.State) []string {
	if s.State != core.StateDefault {
		return []string{}
	}
	return []string{
		"lights",
	}
}

func (c *Lights) Execute(req core.Request) (core.Response, error) {
	deviceGroup, err := findDeviceInCommand(c.Service.DeviceGroups(), req.Body)
	if err != nil {
		return core.Response{}, err
	}

	newState := determineStateRequest(req.Body)

	for _, device := range deviceGroup.Devices {
		err = c.Service.SetStatus(device, newState)
		if err != nil {
			return core.Response{}, err
		}
	}

	newStateStr := "on"
	if !newState {
		newStateStr = "off"
	}

	return core.Response{
		ResponseBody: core.ResponseBody{
			Message: fmt.Sprintf("\"%s\" is now %s.", deviceGroup.PreferredName, newStateStr),
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

func findDeviceInCommand(deviceGroups []service.DeviceGroup, body string) (service.DeviceGroup, error) {
	bestDistance := math.MaxInt
	var best service.DeviceGroup
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
		return service.DeviceGroup{}, core.NewFeedbackError("could not find a matching device")
	}

	return best, nil
}

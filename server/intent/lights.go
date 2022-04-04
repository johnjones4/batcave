package intent

import (
	"fmt"
	"main/core"
	"main/service"
	"strings"
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
	device, err := findDeviceInCommand(c.Service.DeviceNames(), req.Body)
	if err != nil {
		return core.Response{}, err
	}

	newState := determineStateRequest(req.Body)

	err = c.Service.SetStatus(device, newState)
	if err != nil {
		return core.Response{}, err
	}

	newStateStr := "on"
	if !newState {
		newStateStr = "off"
	}

	return core.Response{
		ResponseBody: core.ResponseBody{
			Message: fmt.Sprintf("I've turned the device \"%s\" %s.", device, newStateStr),
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

func findDeviceInCommand(devices []string, body string) (string, error) {
	lc := strings.ToLower(body)
	for _, device := range devices {
		if strings.Contains(lc, device) {
			return device, nil
		}
	}
	return "", core.NewFeedbackError("could not find a matching device")
}

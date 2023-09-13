package intents

import (
	"context"
	"fmt"
	"main/core"
	"main/services/homeassistant"
	"math"
	"strings"

	"github.com/agnivade/levenshtein"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/exp/slices"
)

type ToggleDevice struct {
	HomeAssistant *homeassistant.HomeAssistant
}

type toggleDeviceReceiver struct {
	OnOff string `json:"onOff"`
}

var (
	ResponseNoDevices = core.Response{
		Message: core.Message{
			Text: "No devices found for that request",
		},
	}
)

func (p *ToggleDevice) IntentLabel() string {
	return "toggle_device"
}

func (p *ToggleDevice) IntentParsePrompt(req *core.Request) string {
	return fmt.Sprintf("Determine if this is an \"on\" or \"off\" request: \"%s\" and return the result in JSON formated as: {\"onOff\":\"\"}", req.Message.Text)
}

func (p *ToggleDevice) IntentParseReceiver() any {
	return toggleDeviceReceiver{}
}

func (td *ToggleDevice) ActOnIntent(ctx context.Context, req *core.Request, md *core.IntentMetadata) (core.Response, error) {
	var info toggleDeviceReceiver
	err := mapstructure.Decode(md.IntentParseReceiver, &info)
	if err != nil {
		return core.ResponseEmpty, err
	}
	ids := td.deviceIdsForRequst(req)
	if len(ids) == 0 {
		return ResponseNoDevices, nil
	}

	isOnRequest := info.OnOff == "on"

	for _, id := range ids {
		err := td.HomeAssistant.ToggleDeviceState(id, isOnRequest)
		if err != nil {
			return core.Response{}, err
		}
	}

	return core.ResponseEmpty, nil
}

func (td *ToggleDevice) deviceIdsForRequst(req *core.Request) []string {
	queryLc := strings.ToLower(req.Message.Text)
	lowestIdx := -1
	lowestScore := math.MaxInt64
	for i, group := range td.HomeAssistant.Configuration.Groups {
		if len(group.ClientIds) > 0 && !slices.Contains(group.ClientIds, req.ClientID) {
			continue
		}
		for _, name := range group.Names {
			distance := levenshtein.ComputeDistance(queryLc, name)
			if distance < lowestScore {
				lowestScore = distance
				lowestIdx = i
			}
		}
	}

	if lowestIdx < 0 {
		return []string{}
	}

	return td.HomeAssistant.Configuration.Groups[lowestIdx].DeviceIds
}

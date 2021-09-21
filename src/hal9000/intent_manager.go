package hal9000

import (
	"fmt"
	"hal9000/intents"
	"hal9000/types"
)

func errorNotImplemented(intentType string) error {
	return fmt.Errorf("not implemented: %s", intentType)
}

func GetIntentForIncomingMessage(runtime *types.Runtime, caller *types.Person, m types.ParsedRequestMessage) (types.Intent, error) {
	if m.IntentLabel == "message" {
		return intents.NewMessageIntent(runtime, *caller, m)
	} else if m.IntentLabel == "control_on" {
		return intents.NewControlIntent(runtime, m, true)
	} else if m.IntentLabel == "control_off" {
		return intents.NewControlIntent(runtime, m, false)
	} else if m.IntentLabel == "display" {
		return intents.NewDisplayIntent(runtime, m)
	} else if m.IntentLabel == "weather" {
		return intents.NewWeatherIntent(m)
	} else if m.IntentLabel == "agenda" {
		return intents.NewCalendarAgendaIntent(m)
	} else if m.IntentLabel == "calendar_add" {
		return intents.NewCalendarAddIntent(m)
	} else if m.IntentLabel == "job" {
		return nil, errorNotImplemented(m.IntentLabel)
	} else {
		return nil, fmt.Errorf("no intent for %s", m.IntentLabel)
	}
}

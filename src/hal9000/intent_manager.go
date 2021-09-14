package hal9000

import (
	"fmt"
	"hal9000/intents"
	"hal9000/types"
)

func ErrorNotImplemented(intentType string) error {
	return fmt.Errorf("not implemented: %s", intentType)
}

func GetIntentForIncomingMessage(runtime types.Runtime, intentType string, caller types.Person, m types.ParsedRequestMessage) (types.Intent, error) {
	if intentType == "message" {
		return intents.NewMessageIntent(runtime, caller, m)
	} else if intentType == "control_on" {
		return intents.NewControlIntent(runtime, m, true)
	} else if intentType == "control_off" {
		return intents.NewControlIntent(runtime, m, false)
	} else if intentType == "display" {
		return intents.NewDisplayIntent(runtime, m)
	} else if intentType == "weather" {
		return intents.NewWeatherIntent(m)
	} else if intentType == "agenda" {
		return intents.NewCalendarAgendaIntent(m)
	} else if intentType == "calendar_add" {
		return intents.NewCalendarAddIntent(m)
	} else if intentType == "job" {
		return nil, ErrorNotImplemented(intentType) //TODO
	} else {
		return nil, fmt.Errorf("no intent for %s", intentType)
	}
}

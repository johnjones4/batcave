package hal9000

import (
	"fmt"
)

type Intent interface {
	Execute(lastState State) (State, Message, error)
}

func ErrorNotImplemented(intentType string) error {
	return fmt.Errorf("not implemented: %s", intentType)
}

func GetIntentForIncomingMessage(intentType string, m ParsedMessage) (Intent, error) {
	if intentType == "message" {
		return NewMessageIntent(m)
	} else if intentType == "control_on" {
		return NewControlIntent(m, true)
	} else if intentType == "control_off" {
		return NewControlIntent(m, false)
	} else if intentType == "display" {
		return NewDisplayIntent(m)
	} else if intentType == "weather" {
		return NewWeatherIntent(m)
	} else if intentType == "datetime" {
		return nil, ErrorNotImplemented(intentType) //TODO
	} else if intentType == "job" {
		return nil, ErrorNotImplemented(intentType) //TODO
	} else {
		return nil, fmt.Errorf("no intent for %s", intentType)
	}
}

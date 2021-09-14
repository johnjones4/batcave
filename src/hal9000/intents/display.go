package intents

import (
	"errors"
	"fmt"
	"hal9000/types"
	"hal9000/util"
)

type displayIntent struct {
	display types.Displayable
}

func NewDisplayIntent(runtime types.Runtime, m types.ParsedRequestMessage) (displayIntent, error) {
	display, err := runtime.Displays().FindDisplayableInString(m.Original.Message)
	if err != nil {
		return displayIntent{}, err
	}

	return displayIntent{display}, nil
}

func (i displayIntent) Execute(runtime types.Runtime, lastState types.State) (types.State, types.ResponseMessage, error) {
	if i.display.GetType() == util.DisplayTypeVideo && i.display.GetSource() == util.DisplaySourceGoogle {
		m := types.ResponseMessage{
			Text:  fmt.Sprintf("Here's the %s.", i.display.GetNames()[0]),
			URL:   i.display.GetURL(),
			Extra: i,
		}
		return lastState, m, nil
	}

	return nil, types.ResponseMessage{}, errors.New("unable to handle display type")
}

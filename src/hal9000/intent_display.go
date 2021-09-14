package hal9000

import (
	"errors"
	"fmt"
	"hal9000/util"
)

type DisplayIntent struct {
	Display Displayable
}

func NewDisplayIntent(m ParsedRequestMessage) (DisplayIntent, error) {
	display, err := FindDisplayableInString(m.Original.Message)
	if err != nil {
		return DisplayIntent{}, err
	}

	return DisplayIntent{Display: display}, nil
}

func (i DisplayIntent) Execute(lastState State) (State, util.ResponseMessage, error) {
	if i.Display.Type == DisplayTypeVideo && i.Display.Source == DisplaySourceGoogle {
		m := util.ResponseMessage{
			Text:  fmt.Sprintf("Here's the %s.", i.Display.Names[0]),
			URL:   i.Display.URL,
			Extra: i,
		}
		return lastState, m, nil
	}

	return nil, util.ResponseMessage{}, errors.New("unable to handle display type")
}

package hal9000

import (
	"errors"
	"fmt"
	"hal9000/service"
)

type DisplayIntent struct {
	Display           Display                            `json:"display"`
	GoogleRefreshInfo service.GoogleStreamRefreshRequest `json:"googleRefreshInfo"`
	LastURL           string                             `json:"lastUrl"`
}

func NewDisplayIntent(m ParsedMessage) (DisplayIntent, error) {
	display, err := FindDisplayInString(m.Original)
	if err != nil {
		return DisplayIntent{}, err
	}

	return DisplayIntent{Display: display}, nil
}

func (i DisplayIntent) Execute(lastState State) (State, Message, error) {
	if i.Display.Type == DisplayTypeVideo && i.Display.Source == DisplaySourceGoogle {
		var url string
		var refreshInfo service.GoogleStreamRefreshRequest
		var err error
		if i.GoogleRefreshInfo.StreamExtensionToken != "" && i.LastURL != "" {
			url, refreshInfo, err = service.RefreshGoogleVideoStreamURL(i.LastURL, i.Display.ID, i.GoogleRefreshInfo.StreamExtensionToken)
		} else {
			url, refreshInfo, err = service.GetGoogleVideoStreamURL(i.Display.ID)
		}
		if err != nil {
			return nil, Message{}, err
		}
		i.GoogleRefreshInfo = refreshInfo
		i.LastURL = url
		m := Message{
			Text:  fmt.Sprintf("Here's the %s.", i.Display.Names[0]),
			URL:   url,
			Extra: i,
		}
		return lastState, m, nil
	}

	return nil, Message{}, errors.New("unable to handle display type")
}

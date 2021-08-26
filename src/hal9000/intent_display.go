package hal9000

import (
	"errors"
	"fmt"
	"hal9000/service"
)

type DisplayIntent struct {
	Display           Displayable                        `json:"displayable"`
	GoogleRefreshInfo service.GoogleStreamRefreshRequest `json:"googleRefreshInfo"`
	LastURL           string                             `json:"lastUrl"`
}

func NewDisplayIntent(m ParsedRequestMessage) (DisplayIntent, error) {
	display, err := FindDisplayableInString(m.Original.Message)
	if err != nil {
		return DisplayIntent{}, err
	}

	return DisplayIntent{Display: display}, nil
}

func (i DisplayIntent) Execute(lastState State) (State, ResponseMessage, error) {
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
			return nil, ResponseMessage{}, err
		}
		i.GoogleRefreshInfo = refreshInfo
		i.LastURL = url
		m := ResponseMessage{
			Text:  fmt.Sprintf("Here's the %s.", i.Display.Names[0]),
			URL:   url,
			Extra: i,
		}
		fmt.Println(m)
		return lastState, m, nil
	}

	return nil, ResponseMessage{}, errors.New("unable to handle display type")
}

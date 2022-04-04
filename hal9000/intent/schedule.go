package intent

import (
	"fmt"
	"strings"
	"time"

	"github.com/johnjones4/hal-9000/hal9000/core"
	"github.com/johnjones4/hal-9000/hal9000/service"

	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
)

type Schedule struct {
	Service *service.Google
}

func (c *Schedule) SupportedComandsForState(s core.State) []string {
	if s.State != core.StateDefault {
		return []string{}
	}
	return []string{
		"schedule",
	}
}

func (c *Schedule) Execute(req core.Request) (core.Response, error) {
	w := when.New(nil)
	w.Add(en.All...)
	w.Add(common.All...)

	dateInfo, err := w.Parse(req.Body, time.Now()) //TODO better parsing
	if err != nil {
		return core.Response{}, err
	}

	event := service.Event{
		Name:  strings.TrimSpace(strings.ReplaceAll(req.Body, dateInfo.Text, "")),
		Start: dateInfo.Time,
		End:   dateInfo.Time.Add(time.Hour),
	}

	createdEvent, err := c.Service.CreateNewEvent(event)
	if err != nil {
		return core.Response{}, err
	}

	return core.Response{
		ResponseBody: core.ResponseBody{
			//TODO better formatting
			Message: fmt.Sprintf("Scheduled %s for %s (%s)", createdEvent.Summary, dateInfo.Time.Local().String(), createdEvent.HtmlLink),
		},
		State: req.State,
	}, nil
}

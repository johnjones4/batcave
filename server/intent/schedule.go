package intent

import (
	"main/core"
	"main/service"
	"strings"
	"time"

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

	dateInfo, err := w.Parse(req.Body, time.Now())
	if err != nil {
		return core.Response{}, err
	}

	event := service.Event{
		Name:  strings.ReplaceAll(req.Body, dateInfo.Text, ""),
		Start: dateInfo.Time,
		End:   dateInfo.Time.Add(time.Hour),
	}

	createdEvent, err := c.Service.CreateNewEvent(event)
	if err != nil {
		return core.Response{}, err
	}

	return core.Response{
		ResponseBody: core.ResponseBody{
			Message: "Event scheduled: " + createdEvent.HtmlLink,
		},
		State: req.State,
	}, nil
}

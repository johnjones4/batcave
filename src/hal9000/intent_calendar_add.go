package hal9000

import (
	"errors"
	"fmt"
	"hal9000/service"
	"hal9000/util"
	"strings"
	"time"
)

const DefaultNewEventLengthMinutes = 30

var SchedulerWords = []string{"schedule", "add", "to calendar", "agenda", "to my calendar", "on", "at"}

type CalendarAddIntent struct {
	Title string
	Start time.Time
	End   time.Time
}

func NewCalendarAddIntent(m ParsedRequestMessage) (CalendarAddIntent, error) {
	if m.DateInfo == nil {
		return CalendarAddIntent{}, errors.New("no date for event provided")
	}

	title := strings.ReplaceAll(m.Original.Message, m.DateInfo.Text, "")
	for _, s := range SchedulerWords {
		title = strings.ReplaceAll(title, s, "")
	}
	title = strings.Trim(title, " ")

	return CalendarAddIntent{
		Title: title,
		Start: m.DateInfo.Time,
		End:   m.DateInfo.Time.Add(DefaultNewEventLengthMinutes * time.Minute),
	}, nil
}

func (i CalendarAddIntent) Execute(lastState State) (State, util.ResponseMessage, error) {
	err := service.CreateNewEvent(i.Title, i.Start, i.End)
	if err != nil {
		return nil, util.ResponseMessage{}, err
	}
	response := fmt.Sprintf("Scheduled \"%s\" for %s at %s", i.Title, i.Start.Format("Monday, January 2"), i.Start.Format("3:04pm"))
	return lastState, util.ResponseMessage{Text: response}, nil
}

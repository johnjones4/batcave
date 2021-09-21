package intents

import (
	"errors"
	"fmt"
	"hal9000/types"
	"strings"
	"time"
)

const DefaultNewEventLengthMinutes = 30

var SchedulerWords = []string{"schedule", "add", "to calendar", "agenda", "to my calendar", "on", "at"}

type calendarAddIntent struct {
	title string
	start time.Time
	end   time.Time
}

func NewCalendarAddIntent(m types.ParsedRequestMessage) (calendarAddIntent, error) {
	if m.DateInfo == nil {
		return calendarAddIntent{}, errors.New("no date for event provided")
	}

	title := strings.ReplaceAll(m.Original.Message, m.DateInfo.Text, "")
	for _, s := range SchedulerWords {
		title = strings.ReplaceAll(title, s, "")
	}
	title = strings.Trim(title, " ")

	return calendarAddIntent{
		title: title,
		start: m.DateInfo.Time,
		end:   m.DateInfo.Time.Add(DefaultNewEventLengthMinutes * time.Minute),
	}, nil
}

func (ac calendarAddIntent) GetStartTime() time.Time {
	return ac.start
}

func (ac calendarAddIntent) GetEndTime() time.Time {
	return ac.end
}

func (ac calendarAddIntent) GetName() string {
	return ac.title
}

func (i calendarAddIntent) Execute(runtime *types.Runtime, lastState *types.State) (*types.State, types.ResponseMessage, error) {
	err := (*(*runtime).Google()).CreateNewEvent(runtime, i)
	if err != nil {
		return nil, types.ResponseMessage{}, err
	}
	response := fmt.Sprintf("Scheduled \"%s\" for %s at %s", i.title, i.start.Format("Monday, January 2"), i.start.Format("3:04pm"))
	return lastState, types.ResponseMessage{Text: response}, nil
}

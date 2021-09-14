package intents

import (
	"fmt"
	"hal9000/types"
	"time"
)

type calendarAgendaIntent struct {
	start time.Time
	end   time.Time
}

func NewCalendarAgendaIntent(m types.ParsedRequestMessage) (calendarAgendaIntent, error) {
	date := time.Now()
	if m.DateInfo != nil {
		date = m.DateInfo.Time
	}
	year := date.Year()
	month := date.Month()
	day := date.Day()
	beginningOfDay := time.Date(year, month, day, 0, 0, 0, 0, date.Location())
	endOfDay := beginningOfDay.Add(time.Hour * 24)
	return calendarAgendaIntent{beginningOfDay, endOfDay}, nil
}

func (i calendarAgendaIntent) Execute(runtime types.Runtime, lastState types.State) (types.State, types.ResponseMessage, error) {
	events, err := runtime.Agenda().GetAgendaForDateRange(i.start, i.end)
	if err != nil {
		return nil, types.ResponseMessage{}, err
	}
	s := ""
	if len(events) != 1 {
		s = "s"
	}
	message := fmt.Sprintf("you have %d appointment%s on your calendar on %s", len(events), s, i.start.Format("January 2"))
	if len(events) > 0 {
		message += ":\n"
	} else {
		message += "."
	}
	for i, e := range events {
		message += e.GetName()
		if i < len(events)-2 {
			message += "\n"
		}
	}
	return lastState, types.ResponseMessage{Text: message}, nil
}

package hal9000

import (
	"fmt"
	"time"
)

type CalendarAgendaIntent struct {
	Start time.Time
	End   time.Time
}

func NewCalendarAgendaIntent(m ParsedRequestMessage) (CalendarAgendaIntent, error) {
	date := time.Now()
	if m.DateInfo != nil {
		date = m.DateInfo.Time
	}
	year := date.Year()
	month := date.Month()
	day := date.Day()
	beginningOfDay := time.Date(year, month, day, 0, 0, 0, 0, date.Location())
	endOfDay := beginningOfDay.Add(time.Hour * 24)
	return CalendarAgendaIntent{beginningOfDay, endOfDay}, nil
}

func (i CalendarAgendaIntent) Execute(lastState State) (State, ResponseMessage, error) {
	events, err := GetAgendaForDateRange(i.Start, i.End)
	if err != nil {
		return nil, ResponseMessage{}, err
	}
	s := ""
	if len(events) != 1 {
		s = "s"
	}
	message := fmt.Sprintf("you have %d appointment%s on your calendar on %s", len(events), s, i.Start.Format("January 2"))
	if len(events) > 0 {
		message += ":\n"
	} else {
		message += "."
	}
	for i, e := range events {
		message += e.Name
		if i < len(events)-2 {
			message += "\n"
		}
	}
	return lastState, ResponseMessage{Text: message}, nil
}

package service

import (
	"encoding/json"
	"hal9000/types"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/apognu/gocal"
)

type AgendaCalendar struct {
	ICSURL string `json:"icsUrl"`
}

type agendaCalendarEvent struct {
	start time.Time
	end   time.Time
	name  string
}

func (ac agendaCalendarEvent) GetStartTime() time.Time {
	return ac.start
}

func (ac agendaCalendarEvent) GetEndTime() time.Time {
	return ac.end
}

func (ac agendaCalendarEvent) GetName() string {
	return ac.name
}

type calendarProviderConcrete struct {
	calendars []AgendaCalendar
}

func InitAgendaProvider() (types.AgendaProvider, error) {
	bytes, err := ioutil.ReadFile(os.Getenv("CALENDARS_MANIFEST_PATH"))
	if err != nil {
		return nil, err
	}
	var calendars []AgendaCalendar
	err = json.Unmarshal(bytes, &calendars)
	if err != nil {
		return nil, err
	}
	return calendarProviderConcrete{calendars}, nil
}

func (cp calendarProviderConcrete) GetAgendaForDateRange(start time.Time, end time.Time) ([]types.Event, error) {
	events := make([]types.Event, 0)
	for _, calendar := range cp.calendars {
		response, err := http.Get(calendar.ICSURL)
		if err != nil {
			return nil, err
		}
		c := gocal.NewParser(response.Body)
		c.Start, c.End = &start, &end
		err = c.Parse()
		if err != nil {
			return nil, err
		}
		for _, e := range c.Events {
			events = append(events, agendaCalendarEvent{start: *e.Start, end: *e.End, name: e.Summary})
		}
	}
	return events, nil
}

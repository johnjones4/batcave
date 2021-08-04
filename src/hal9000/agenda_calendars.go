package hal9000

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/apognu/gocal"
)

type AgendaCalendar struct {
	ICSURL string `json:"icsUrl"`
}

type AgendaCalendarEvent struct {
	Start time.Time
	End   time.Time
	Name  string
}

var calendars []AgendaCalendar

func InitCalendarSchedules() error {
	bytes, err := ioutil.ReadFile(os.Getenv("CALENDARS_MANIFEST_PATH"))
	if err != nil {
		return err
	}
	calendars = nil
	err = json.Unmarshal(bytes, &calendars)
	if err != nil {
		return err
	}
	return nil
}

func GetAgendaForDateRange(start time.Time, end time.Time) ([]AgendaCalendarEvent, error) {
	events := make([]AgendaCalendarEvent, 0)
	for _, calendar := range calendars {
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
			events = append(events, AgendaCalendarEvent{*e.Start, *e.End, e.Summary})
		}
	}
	return events, nil
}

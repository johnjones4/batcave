package intent

import (
	"fmt"
	"time"

	"github.com/johnjones4/hal-9000/server/hal9000/core"
	"github.com/johnjones4/hal-9000/server/hal9000/service"
	"github.com/johnjones4/hal-9000/server/hal9000/util"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
)

type Forecast struct {
	Service *service.NOAA
}

func (c *Forecast) SupportedComandsForState(s string) map[string]core.CommandInfo {
	if s != core.StateDefault {
		return map[string]core.CommandInfo{}
	}
	return map[string]core.CommandInfo{
		"forecast": {
			Description:  "Get the forecast for the given date/time.",
			RequiresBody: true,
		},
	}
}

func (c *Forecast) Execute(req core.Inbound) (core.Outbound, error) {
	weatherDate := time.Now()

	if req.Body != "" {
		w := when.New(nil)
		w.Add(en.All...)
		w.Add(common.All...)

		dateInfo, err := w.Parse(req.Body, time.Now()) //TODO better parsing
		if err != nil {
			return core.Outbound{}, err
		}
		if dateInfo != nil && !dateInfo.Time.IsZero() {
			weatherDate = dateInfo.Time
		} else {
			weatherDate = time.Now()
		}
	}

	info, err := c.Service.PredictWeather(req.Location)
	if err != nil {
		return core.Outbound{}, err
	}

	if len(info.Forecast) == 0 {
		return core.Outbound{}, core.NewFeedbackError("There is no forecast for your current area.")
	}

	var forecast service.NOAAWeatherForecastPeriod
	radar := ""
	for i, f := range info.Forecast {
		if weatherDate.After(f.StartTime) && weatherDate.Before(f.EndTime) {
			forecast = f
			if i == 0 {
				radar = info.RadarURL
			}
			break
		}
	}

	if forecast.DetailedForecast == "" {
		return core.Outbound{}, core.NewFeedbackError("No weather available")
	}

	resp := core.Outbound{
		OutboundBody: core.OutboundBody{
			Body:  fmt.Sprintf("Forecast for %s to %s: %s", util.FormatTime(forecast.StartTime), util.FormatTime(forecast.EndTime), forecast.DetailedForecast),
			Media: radar,
		},
		State: req.State,
	}

	return resp, nil
}

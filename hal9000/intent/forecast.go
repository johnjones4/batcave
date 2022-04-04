package intent

import (
	"github.com/johnjones4/hal-9000/hal9000/core"
	"github.com/johnjones4/hal-9000/hal9000/service"
)

type Forecast struct {
	Service *service.NOAA
}

func (c *Forecast) SupportedComandsForState(s core.State) []string {
	if s.State != core.StateDefault {
		return []string{}
	}
	return []string{
		"forecast",
	}
}

func (c *Forecast) Execute(req core.Request) (core.Response, error) {
	info, err := c.Service.PredictWeather(req.Location)
	if err != nil {
		return core.Response{}, err
	}

	if len(info.Forecast) == 0 {
		return core.Response{}, core.NewFeedbackError("There is no forecast for your current area.")
	}

	resp := core.Response{
		ResponseBody: core.ResponseBody{
			Message: info.Forecast[0].DetailedForecast, //TODO date, alerts?
			Media:   info.RadarURL,
		},
		State: req.State,
	}

	return resp, nil
}

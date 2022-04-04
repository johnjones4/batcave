package intent

import (
	"errors"
	"main/core"
	"main/service"
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
		return core.Response{}, errors.New("no forecast data returned")
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

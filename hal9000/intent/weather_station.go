package intent

import (
	"fmt"

	"github.com/johnjones4/hal-9000/hal9000/core"
	"github.com/johnjones4/hal-9000/hal9000/service"
)

type WeatherStation struct {
	Service *service.WeatherStation
}

func (w *WeatherStation) SupportedComandsForState(s core.State) []string {
	if s.State != core.StateDefault {
		return []string{}
	}
	return []string{
		"weather",
	}
}

func (w *WeatherStation) Execute(req core.Request) (core.Response, error) {
	info, err := w.Service.GetWeather()
	if err != nil {
		return core.Response{}, err
	}

	return core.Response{
		ResponseBody: core.ResponseBody{
			Message: fmt.Sprintf("Weather station report:\nTemperature: %0.2f° F\nAverage wind speed: %0.2f m/s\nRelative humidity: %0.2f\nPressure: %0.2f inhg.",
				cToF(info.Temperature),
				mpsToMph(info.AvgWindSpeed),
				info.RelativeHumidity,
				mbarToInHg(info.Pressure),
			),
		},
		State: req.State,
	}, nil
}

func cToF(c float64) float64 {
	return c*(9.0/5.0) + 32
}

func mpsToMph(m float64) float64 {
	return m * 2.23694
}

func mbarToInHg(m float64) float64 {
	return m / 33.863886666667
}

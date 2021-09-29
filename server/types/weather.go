package types

import "time"

type Weather struct {
	RadarURL string                              `json:"radarURL"`
	Forecast []NOAAWeatherForecastPeriod         `json:"forecast"`
	Alerts   []NOAAWeatherAlertFeatureProperties `json:"alerts"`
}

type NOAAWeatherPointProperties struct {
	ForecastURL  string `json:"forecast"`
	ForecastZone string `json:"forecastZone"`
	RadarStation string `json:"radarStation"`
}

type NOAAWeatherPointResponse struct {
	Properties NOAAWeatherPointProperties `json:"properties"`
}

type NOAAWeatherForecastPeriod struct {
	StartTime        time.Time `json:"startTime"`
	EndTime          time.Time `json:"endTime"`
	DetailedForecast string    `json:"detailedForecast"`
	Name             string    `json:"name"`
	Temperature      float64   `json:"temperature"`
	TemperatureUnit  string    `json:"temperatureUnit"`
	WindSpeed        string    `json:"windSpeed"`
	WindDirection    string    `json:"windDirection"`
	Icon             string    `json:"icon"`
	IsDaytime        bool      `json:"isDaytime"`
}

type NOAAWeatherForecastProperties struct {
	Periods []NOAAWeatherForecastPeriod `json:"periods"`
}

type NOAAWeatherForecastResponse struct {
	Properties NOAAWeatherForecastProperties `json:"properties"`
}

type NOAAWeatherAlertFeatureProperties struct {
	ID            string   `json:"id"`
	AffectedZones []string `json:"affectedZones"`
	Headline      string   `json:"headline"`
}

type NOAAWeatherAlertFeature struct {
	Properties NOAAWeatherAlertFeatureProperties `json:"properties"`
}

type NOAAWeatherAlertResponse struct {
	Features []NOAAWeatherAlertFeature `json:"features"`
}

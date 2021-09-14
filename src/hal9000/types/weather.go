package types

import "time"

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

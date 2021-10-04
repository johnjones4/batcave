export interface IFrame {
  url: string
  title: string
}

export interface WeatherForecastItem {
  detailedForecast: string
  name: string
  temperature: number
  temperatureUnit: string
  windSpeed: string
  windDirection: string
  icon: string
}


export interface Weather {
  radarURL: string
  forecast: [WeatherForecastItem]
}

export interface Response {
  iframes: [IFrame]
  weather: [Weather]
}
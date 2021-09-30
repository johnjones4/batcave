export interface Coordinate {
  latitude: number
  longitude: number
}

export interface IFrame {
  url: string
  title: string
}

export interface Configuration {
  iframes: [IFrame]
  weatherLocations: [Coordinate]
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

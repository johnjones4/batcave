import React, { Component } from 'react'
import { Weather } from '../Types'
import Widget from './Widget'
import './WeatherWidget.css'

interface WeatherWidgetProps {
  weather: Weather
}

interface WeatherWidgetState {
}

class WeatherWidget extends Component<WeatherWidgetProps, WeatherWidgetState> {
  constructor(props: WeatherWidgetProps) {
    super(props)
    this.state = {
    }
  }

  render() {
    return (
      <Widget classNameSuffix='Weather' title='Weather'>
        { this.props.weather.forecast.slice(0,3).map((forecast, i) => (
          <div key={i} className='Weather-Forecast-Item'>
            <h2 className='Weather-Forecast-Item-Title'>{forecast.name}</h2>
            <img className='Weather-Forecast-Item-Icon' src={forecast.icon} />
            <div className='Weather-Forecast-Item-Temp'>{forecast.temperature}&deg;</div>
            <div className='Weather-Forecast-Item-Wind'>
              <span className='Weather-Forecast-Item-Wind-Speed'>{forecast.windSpeed}</span>/
              <span className='Weather-Forecast-Item-Wind-Direction'>{forecast.windDirection}</span>
            </div>
            <div className='Weather-Forecast-Item-Description'>{forecast.detailedForecast}</div>
          </div>
        )) }
      </Widget>
    )
  }
}

export default WeatherWidget

import React, { Component } from 'react'
import { Weather } from '../Types'
import Widget from './Widget'

interface WeatherWidgetProps {
  weather: Weather
}

interface WeatherWidgetState {
  index: number
}

class WeatherWidget extends Component<WeatherWidgetProps, WeatherWidgetState> {
  constructor(props: WeatherWidgetProps) {
    super(props)
    this.state = {
      index: 0
    }
  }

  render() {
    return (
      <Widget classNameSuffix='Weather'>
        { this.props.weather.forecast.map(forecast => (
          <div className='Weather-Forecast-Item'>
            
          </div>
        )) }
      </Widget>
    )
  }
}

export default WeatherWidget

import React, { Component } from 'react'
import { Configuration, Weather } from './Types'
import { Loader } from './Loader'
import './Dashboard.css'
import IFrameWidget from './widgets/IFrameWidget'
import ImageWidget from './widgets/ImageWidget'
import WeatherWidget from './widgets/WeatherWidget'
interface DashboardProps {

}

interface DashboardState {
  configuration: Configuration | null
  weather: [Weather] | null
  error: any
}

export default class Dashboard extends Component<DashboardProps,DashboardState> {
  configurationLoader: Loader<Configuration>
  weatherLoader: Loader<[Weather]>

  constructor(props: DashboardProps) {
    super(props)
    this.state = {
      configuration: null,
      weather: null,
      error: null
    }
    this.configurationLoader = new Loader<Configuration>('/api/configuration')
    this.weatherLoader = new Loader<[Weather]>('/api/weather')
  }

  componentDidMount() {
    this.load()
  }

  async load() {
    try {
      this.setState({
        configuration: await this.configurationLoader.load(),
        weather: await this.weatherLoader.load()
      })
    } catch (e: any) {
      this.setState({
        error: e
      })
    }
  }

  render() {
    return (
      <div className='Dashboard'>
        { this.state.configuration && this.state.configuration.iframes.map(iframe => (<IFrameWidget iframe={iframe} key={iframe.url} />)) }
        { this.state.weather && this.state.weather.map(weather => (<ImageWidget src={weather.radarURL} key={weather.radarURL} />)) }
        { this.state.weather && this.state.weather.map(weather => (<WeatherWidget weather={weather} key={weather.radarURL} />)) }
      </div>
    )
  }
}

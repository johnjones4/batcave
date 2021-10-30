import React, { Component } from 'react'
import { Response, Weather } from './Types'
import { Loader } from './Loader'
import './Dashboard.css'
import IFrameWidget from './widgets/IFrameWidget'
import ImageWidget from './widgets/ImageWidget'
import WeatherWidget from './widgets/WeatherWidget'
import NewsWidget from './widgets/NewsWidget'
interface DashboardProps {

}

interface DashboardState {
  response: Response | null
  error: any
}

export default class Dashboard extends Component<DashboardProps,DashboardState> {
  responseLoader: Loader<Response>

  constructor(props: DashboardProps) {
    super(props)
    this.state = {
      response: null,
      error: null
    }
    this.responseLoader = new Loader<Response>('/api/data')
  }

  componentDidMount() {
    this.load()
  }

  async load() {
    try {
      this.setState({
        response: await this.responseLoader.load(),
      })
      setInterval(() => this.load(), 1000 * 60 * 5)
    } catch (e: any) {
      this.setState({
        error: e
      })
    }
  }

  render() {
    return (
      <div className='Dashboard'>
        { this.state.response && this.state.response.iframes.map(iframe => (<IFrameWidget iframe={iframe} key={iframe.url} />)) }
        {/* { this.state.response && this.state.response.weather.map(weather => (<ImageWidget src={weather.radarURL} key={weather.radarURL} title='Radar' />)) } */}
        { this.state.response && this.state.response.weather.map(weather => (<WeatherWidget weather={weather} key={weather.radarURL} />)) }
        { this.state.response && (<NewsWidget news={this.state.response.news} />) }
      </div>
    )
  }
}

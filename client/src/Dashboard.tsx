import React, { Component } from 'react'
import { Configuration } from './Configuration'
import { Loader } from './Loader'

interface DashboardProps {

}

interface DashboardState {
  configuration: Configuration | null
  error: Error | null
}

export default class Dashboard extends Component<DashboardProps,DashboardState> {
  configurationLoader: Loader<Configuration>

  constructor(props: DashboardProps) {
    super(props)
    this.state = {
      configuration: null,
      error: null
    }
    this.configurationLoader = new Loader<Configuration>('/api/config')
    this.reload()
  }

  async reload() {
    try {
      const config = await this.configurationLoader.load()
      this.setState({
        configuration: config
      })
    } catch (e) {
      this.setState({
        error: e
      })
    }
  }

  render() {
    
  }
}
import React, { Component } from 'react'
import './App.css'
import Chat, { LogItem } from './Chat'
import { HAL9000, HAL9000Response } from './HAL9000'
import Visual, { VisualItem } from './Visual'

const maxVisuals = 5

interface AppProps {

}

interface AppState {
  visuals: VisualItem[]
  log: LogItem[]
}

export default class App extends Component<AppProps,AppState> {
  private hal : HAL9000
  private chat?: Chat

  constructor(props: AppProps) {
    super(props)
    this.state = {
      visuals: [],
      log: []
    }
    this.hal = new HAL9000(`ws://${window.location.host}/ws?user=john&visual=true`, this)
  }

  handleInput(input: string) {
    const req = this.hal.send(input)
    if (req) {
      this.setState({
        log: this.state.log.concat([{
          request: req,
          response: null,
          error: null
        }])
      })
    }
  }

  handleError (e: any) {
    this.setState({
      log: this.state.log.concat([{
        request: null,
        response: null,
        error: e
      }])
    })
  }

  handleResponse (r: HAL9000Response) {
    let visuals = this.state.visuals

    if (r.url !== '') {
      const item = {
        response: r,
        pinned: false,
        added: new Date()   
      }
      if (this.state.visuals.length >= maxVisuals) {
        let oldestUnpinnedIndex = -1
        let oldestUnpinnedDate = new Date()
        this.state.visuals.forEach((v, i) => {
          if (!v.pinned && v.added.getTime() <= oldestUnpinnedDate.getTime()) {
            oldestUnpinnedIndex = i
            oldestUnpinnedDate = v.added
          }
        })
        if (oldestUnpinnedIndex >= 0) {
          visuals = visuals.map((v, i) => {
            if (i === oldestUnpinnedIndex) {
              return item
            }
            return v
          })
        }
      } else {
        visuals = visuals.concat([item])
      }
    }
    this.setState({
      log: this.state.log.concat([{
        request: null,
        response: r,
        error: null
      }]),
      visuals
    })
    if (this.chat) {
      this.chat.focus()
    }
  }

  render () {
    return (
      <div className="App">
        { this.state.visuals.map((v, i) => (
          <Visual
            visual={v}
            key={i}
            onPinnedToggle={(t) => {
              this.setState({
                visuals: this.state.visuals.map((v1, i1) => {
                  if (i !== i1) {
                    return v1
                  }
                  return {
                    ...v1,
                    pinned: t
                  }
                })
              })
            }}
          />
        )) }
        { Array.from(Array(maxVisuals - this.state.visuals.length).keys()).map(() => (<div className='Visual' />)) }
        <Chat
          log={this.state.log}
          onSubmit={(input: string) => this.handleInput(input)}
          onFocusable={(c) => {
            this.chat = c
          }}
        />
      </div>
    );
  }
}

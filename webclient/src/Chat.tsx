import { BaseSyntheticEvent } from "react";
import { Component } from "react";
import { HAL9000Request, HAL9000Response } from "./HAL9000";

export interface LogItem {
  request: HAL9000Request | null
  response: HAL9000Response | null
  error: any | null
}

interface ChatProps {
  log: LogItem[]
  onSubmit(input: string): void
  onFocusable(chat: Chat): void
}

interface ChatState {
  input: string
}

const senderForLogItem = (item: LogItem) : string => {
  if (item.error) {
    return 'Error'
  }
  if (item.request) {
    return 'You'
  }
  if (item.response) {
    return 'HAL'
  }
  return ''
}

const textForLogItem = (item: LogItem) : string => {
  if (item.error) {
    return item.error
  }
  if (item.request) {
    return item.request.message
  }
  if (item.response) {
    return item.response.text
  }
  return ''
}

export default class Chat extends Component<ChatProps,ChatState> {
  private inputField? : HTMLInputElement
  private logEl? : HTMLDivElement

  constructor(props: ChatProps) {
    super(props)
    this.state = {
      input: ''
    }
  }

  submitInput(e: BaseSyntheticEvent) {
    e.preventDefault()
    this.props.onSubmit(this.state.input)
    this.setState({
      input: ''
    })
  }

  componentDidMount() {
    if (this.inputField) {
      this.props.onFocusable(this)
    }
  }

  public focus() {
    this.inputField?.focus()
    this.logEl?.scrollTo(0,this.logEl.scrollHeight)
  }
  
  render () {
    return (
      <div className='Chat'>
        <div className='Chat-log' ref={(e) => {
          this.logEl = e as HTMLDivElement
        }}>
          <ol>
            { this.props.log.map((item, i) => (
              <li className={['Chat-item', 'Chat-item-' + senderForLogItem(item).toLowerCase()].join(' ')} key={i}>
                <span className='Chat-item-sender'>
                  { senderForLogItem(item) }:
                </span>
                <span className='Chat-item-text'>
                  { textForLogItem(item) }
                </span>
              </li>
            )) }
          </ol>
        </div>
        <form className='Chat-input' onSubmit={e => this.submitInput(e)}>
          <input
            type='text'
            value={this.state.input}
            onChange={(e) => this.setState({input: e.target.value})}
            className='Chat-input-text'
            ref={(e) => {
              this.inputField = e as HTMLInputElement
            }}
          />
          <button
            className='Chat-input-button'
            type='submit'
          >Send</button>
        </form>
      </div>
    )
  }
}
import React from 'react'
import { HAL9000Response } from './HAL9000'

export interface VisualItem {
  response: HAL9000Response
  pinned: boolean
  added: Date
}

export interface VisualProps {
  visual: VisualItem
  onPinnedToggle(on: boolean): void
}

const isImage = (url: string): boolean => {
  if (url.endsWith('.jpg') || url.endsWith('.jpeg') || url.endsWith('.gif') || url.endsWith('.png')) {
    return true
  }
  return false
}

const renderContent = (url: string) => {
  if (isImage(url)) {
    return (<div className='Visual-image' style={{backgroundImage: `url(${url})`}} />)
  }
  return (<iframe src={url} />)
}

const Visual = (props: VisualProps) => {
  return (
    <div className='Visual'>
      { renderContent(props.visual.response.url) }
      <input
        className='Visual-pin'
        type='checkbox'
        checked={props.visual.pinned}
        onChange={(e) => props.onPinnedToggle(e.target.checked)}
      />
    </div>
  )
}

export default Visual

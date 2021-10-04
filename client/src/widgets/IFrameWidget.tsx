import React from 'react'
import Widget from './Widget'
import './IFrameWidget.css'
import { IFrame } from '../Types'

interface IFrameProps {
  iframe: IFrame
}

const IFrameWidget = (props: IFrameProps) => {
  return (
    <Widget classNameSuffix='IFrame' title={props.iframe.title}>
      <iframe title={props.iframe.title} src={props.iframe.url} className='IFrame-iframe' />
    </Widget>
  )
}

export default IFrameWidget

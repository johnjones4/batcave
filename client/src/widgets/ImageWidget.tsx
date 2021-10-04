import React from 'react'
import Widget from './Widget'
import './ImageWidget.css'

interface ImageWidgetProps {
  src: string
  title: string
}

const ImageWidget = (props: ImageWidgetProps) => {
  return (
    <Widget classNameSuffix='Image' title={props.title}>
      <div className='ImageWidget-div' style={{backgroundImage: `url(${props.src})`}} />
    </Widget>
  )
}

export default ImageWidget

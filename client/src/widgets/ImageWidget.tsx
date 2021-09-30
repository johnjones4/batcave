import React from 'react'
import Widget from './Widget'
import './ImageWidget.css'

interface ImageWidgetProps {
  src: string
}

const ImageWidget = (props: ImageWidgetProps) => {
  return (
    <Widget classNameSuffix='Image'>
      <div className='ImageWidget-div' style={{backgroundImage: `url(${props.src})`}} />
    </Widget>
  )
}

export default ImageWidget

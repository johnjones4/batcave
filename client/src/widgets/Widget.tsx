import React from 'react'
import './Widget.css'

interface WidgetProps {
  children: any
  classNameSuffix: string
  title: string
}

const Widget = (props: WidgetProps) => {
  return (
    <div className={`Widget Widget-${props.classNameSuffix}Widget`}>
      <h1>{props.title}</h1>
      <div className='Widget-inner'>
        <div className='Widget-inner-inner'>
          { props.children }
        </div>
      </div>
    </div>
  )
}

export default Widget

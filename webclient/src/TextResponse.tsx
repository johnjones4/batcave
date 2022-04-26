import React from 'react'
import { Outbound } from './HAL9000'
import './TextResponse.css'

interface TextResponseProps {
  info: Outbound
}

export const TextResponse = (props: TextResponseProps) => {
  return (
    <div className='TextResponse'>
      <pre>
        {props.info.body}
      </pre>
    </div>
  )
}

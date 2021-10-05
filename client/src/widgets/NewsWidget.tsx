import React, { Component } from 'react'
import Widget from './Widget'
import './NewsWidget.css'
import { NewsItem } from '../Types'

interface NewsWidgetProps {
  news: NewsItem[]
}

const NewsWidget = (props: NewsWidgetProps) => {
  const divs = new Array<HTMLDivElement>(props.news.length)
  let current = 0
  setInterval(() => {
    if (divs[current]) {
      const parent = divs[current].parentElement!.parentElement! as HTMLDivElement
      parent.scrollTo({
        top: divs[current].offsetTop - 20,
        behavior: 'smooth'
      })
      current++
      if (current >= divs.length) {
        current = 0
      }
    }
  }, 5000)
  return (
    <Widget classNameSuffix='News' title='News'>
      { props.news.map((n, i) => (
        <div className='NewsWidget-Item' key={i} ref={(e) => {
          divs[i] = e as HTMLDivElement
        }}>
          <h3>{n.headline}</h3>
          <h4>{n.source}</h4>
          <p>{n.description}</p>
        </div>
      )) }
    </Widget>
  )
}

export default NewsWidget

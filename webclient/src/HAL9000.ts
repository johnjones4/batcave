import CryptoJS from 'crypto-js'

interface Message {
  body: string
}

interface Coordinate {
  latitude: number
  longitude: number
}

export interface AudioSample {
  mimeType: string
  data: string
  gzipped: boolean
}

export interface Inbound extends Message {
  location: Coordinate
  audio: AudioSample
}

export interface Outbound extends Message {
  body: string
  media: string
  url: string
}

interface Pong {
  pong: boolean
}

interface HALError extends Error {
  error: String
  code?: number
}

export default class HAL9000 {
  private clientId: string
  private key: string

  constructor(clientId: string, key: string) {
    this.clientId = clientId
    this.key = key
  }

  private async request<V>(method: string, path: string, body?: any): Promise<V> {
    let params = {
      method,
      headers: new Headers({})
    }

    if (body) {
      (params as any).body = JSON.stringify(body)
      params.headers.set('Content-type', 'application/json')
    }

    const reqTime = new Date().toISOString()
    const message = `${this.clientId}:${reqTime}`

    let hash = CryptoJS.HmacSHA256(message, this.key)
    let sig = CryptoJS.enc.Hex.stringify(hash)

    params.headers.set('X-Request-Time', reqTime)
    params.headers.set('User-Agent', this.clientId)
    params.headers.set('X-Signature', sig)
    
    const res = await fetch(path, params)
    const info = await res.json()
    if (res.status >= 300) {
      throw info as HALError
    }    
    return info as V
  }

  async send(message: Inbound): Promise<Outbound> {
    return this.request("POST", "/api/request", message)
  }

  async ping(): Promise<Pong> {
    return this.request("GET", "/api/ping", null)
  }
}
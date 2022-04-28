const CryptoJS = require('crypto-js')
const fetch = require('node-fetch')

class Client {
  clientId
  key

  constructor(urlRoot, clientId, key) {
    this.urlRoot = urlRoot
    this.clientId = clientId
    this.key = key
  }

  async request(method, path, body) {
    let params = {
      method,
      headers: {}
    }

    if (body) {
      params.body = JSON.stringify(body)
      params.headers['Content-type'] = 'application/json'
    }

    const reqTime = new Date().toISOString()
    const message = `${this.clientId}:${reqTime}`

    let hash = CryptoJS.HmacSHA256(message, this.key)
    let sig = CryptoJS.enc.Hex.stringify(hash)

    params.headers['X-Request-Time'] = reqTime
    params.headers['User-Agent'] = this.clientId
    params.headers['X-Signature'] = sig
    
    const res = await fetch(this.urlRoot+path, params)
    const info = await res.json()
    if (res.status >= 300) {
      throw new Error(info)
    }    
    return info
  }

  async send(message) {
    return this.request("POST", "/api/request", message)
  }

  async sendText(location, body) {
    return this.send({
      body,
      location,
      audio: {
        mimeType: '',
        data: '',
        gzipped: false
      }
    })
  }

  async ping() {
    return this.request("GET", "/api/ping", null)
  }
}

module.exports = Client

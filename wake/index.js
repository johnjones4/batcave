const { app, BrowserWindow } = require('electron')
const Listener = require('./lib/Listener')
const Client = require('./lib/Client')
const {urlRoot, halId, halKey, location} = require('./res/priv/credentials')
const path = require('path')
const fs = require('fs/promises')

let win
const buffer = []

const pushToBuffer = async (command) => {
  buffer.push(command)
  if (buffer.length >= 5000) {
    buffer.shift()
  }
  await fs.writeFile(process.env.LOG, buffer.join('\n'), 'utf-8')
}

const createWindow = () => {
  win = new BrowserWindow({
    width: 1600,
    height: 900,
    webPreferences: {
      preload: path.join(__dirname, 'preload.js')
    }
  })

  win.loadFile('index.html')
}

const start = async () => {
  try {
    const client = new Client(urlRoot, halId, halKey)
    await client.ping()
    const l = new Listener('./res/priv/model.tflite', './res/priv/huge-vocabulary.scorer', 'computer') 
    l.on('activate', () => {
      win.webContents.send('listening', true)
      console.log('listening')
    })
    l.on('deactivate', () => {
      win.webContents.send('listening', false)
      console.log('not listening')
    })
    l.on('command', async (result) => {
      try {
        win.webContents.send('sending', result)
        console.log(result)
        const response = await client.sendText(location, result)
        win.webContents.send('response', response)
        console.log(response)
        await pushToBuffer(result)
      } catch (e) {
        console.error(e)
        win.webContents.send('error', e)
      }
    })
    l.start()

    app.whenReady().then(() => {
      createWindow()
    
      app.on('activate', () => {
        if (BrowserWindow.getAllWindows().length === 0) createWindow()
      })
    })
  } catch (e) {
    console.error(e)
  }
}
start()

const { contextBridge, ipcRenderer } = require('electron')

contextBridge.exposeInMainWorld('hal', {
  onListening: (callback) => ipcRenderer.on('listening', callback), 
  onSending: (callback) => ipcRenderer.on('sending', callback), 
  onResponse: (callback) => ipcRenderer.on('response', callback),
  onError: (callback) => ipcRenderer.on('error', callback)
})

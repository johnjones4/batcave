import React, { useEffect, useState } from 'react';
import './App.css';
import { halId, halKey } from './credentials';
import HAL9000, { Inbound } from './HAL9000';
import { blobToBase64 } from './util';

function App() {
  const hal = new HAL9000(halId, halKey)
  let mediaRecorder: MediaRecorder|null = null
  let buffer: Blob[] = []
  let mimeType: string
  let [dataInfo, setDataInfo] = useState('')

  const ping = async () => {
    try {
      await hal.ping()
    } catch (e) {
      alert(e)
    }
  }

  const sendAudio = async () => {
    try {
      const data = await blobToBase64(new Blob(buffer))
      const inbound: Inbound = {
        location: {
          latitude: 0.0,
          longitude: 0.0
        },
        body: '',
        audio: {
          mimeType,
          data,
          gzipped: false
        }
      }
      const info = await hal.send(inbound)
      setDataInfo(JSON.stringify(info))
    } catch (e) {
      alert(e)
    }
  }

  const setupRecorder = async () => {
    const ms = await navigator.mediaDevices.getUserMedia({audio: true})
    const mr = new MediaRecorder(ms, {
      mimeType: 'audio/ogg',
      audioBitsPerSecond: 16000
    });
    console.log(mr.audioBitsPerSecond)
    mr.onstart = () => {
      buffer = []
    }
    mr.ondataavailable = e => {
      buffer.push(e.data)
      mimeType = e.data.type
    }
    mr.onstop = e => {
      sendAudio()      
    }
    mediaRecorder = mr
  }

  const registerRecordTracker = () => {
    console.log('init')
    document.onkeydown = e => {
      if (e.key === ' ' && mediaRecorder && mediaRecorder.state !== 'recording') {
        console.log('press')
        mediaRecorder.start()
      }
    }
    document.onkeyup = e => {
      if (e.key === ' ' && mediaRecorder) {
        mediaRecorder.stop()
      }
    }
  }

  useEffect(() => {
    ping()
    setupRecorder()
    registerRecordTracker()
  }, [])
  return (
    <div className="App">
      <div className="App-response"></div>
      {dataInfo}
    </div>
  );
}

export default App;

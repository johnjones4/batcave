const STT = require('stt');
const VAD = require('node-vad');
const mic = require('mic');

const SILENCE_THRESHOLD = 200

class Listener {
  constructor(modelPath, scorerPath, wakeWord) {
    this.wakeWord = wakeWord
    this.englishModel = new STT.Model(modelPath)
    this.englishModel.enableExternalScorer(scorerPath)
    this.vad = new VAD(VAD.Mode.VERY_AGGRESSIVE)
    this.recordedChunks = 0
    this.silenceStart = null
    this.recordedAudioLength = 0
    this.endTimeout = null
    this.silenceBuffers = []
    this.listeners = {}
    this.awake = false
  }

  start() {
    if (this.microphone) {
      this.microphone.stop()
    }
    this.createStream()
    this.microphone = mic({
      rate: '16000',
      channels: '1',
      debug: false,
      fileType: 'wav'
    })
    this.microphone.getAudioStream().on('data', data => {
      this.processAudioStream(data)
    })
    this.microphone.start()
  }

  processAudioStream(data) {
    // console.log(data.length)
    this.vad.processAudio(data, 16000).then((res) => {
      switch (res) {
        case VAD.Event.ERROR:
          console.log("VAD ERROR")
          break;
        case VAD.Event.NOISE:
          console.log("VAD NOISE")
          break;
        case VAD.Event.SILENCE:
          this.processSilence(data)
          break;
        case VAD.Event.VOICE:
          console.log('voice')
          this.processVoice(data)
          break;
        default:
          console.log('default', res)
      }
    });
    
    // timeout after 1s of inactivity
    clearTimeout(this.endTimeout);
    this.endTimeout = setTimeout(() => {
      console.log('timeout')
      this.resetAudioStream()
    },SILENCE_THRESHOLD*3)
  }

  processVoice(data) {
    this.silenceStart = null
    this.recordedChunks++
    data = this.addBufferedSilence(data)
    this.feedAudioContent(data)
  }

  processSilence(data) {
    if (this.recordedChunks > 0) {      
      this.feedAudioContent(data)
      
      if (this.silenceStart === null) {
        this.silenceStart = new Date().getTime()
      } else {
        let now = new Date().getTime()
        if (now - this.silenceStart > SILENCE_THRESHOLD) {
          this.silenceStart = null
          let results = this.intermediateDecode()
          if (results && results.text.toLowerCase().indexOf(this.wakeWord +' ') === 0 && results.text.length > this.wakeWord.length + 1 && this.listeners['command']) {
            if (this.listeners['deactivate']) {
              this.listeners['deactivate']()
            }
            this.listeners['command'](results.text.substring(this.wakeWord.length + 1))
          }
        }
      }
    } else {
      this.bufferSilence(data)
    }
  }

  resetAudioStream() {
    clearTimeout(this.endTimeout)
    this.intermediateDecode()
    this.recordedChunks = 0
    this.silenceStart = null
    this.awake = false
    const callback = this.listeners['deactivate']
    if (callback) {
      callback()
    }
  }
  
  addBufferedSilence(data) {
    let audioBuffer
    if (this.silenceBuffers.length) {
      this.silenceBuffers.push(data)
      let length = 0;
      this.silenceBuffers.forEach(buf => {
        length += buf.length
      })
      audioBuffer = Buffer.concat(this.silenceBuffers, length)
      this.silenceBuffers = []
    } else {
      audioBuffer = data
    }
    return audioBuffer
  }
  
  feedAudioContent(chunk) {
    if (!this.awake && this.modelStream.intermediateDecode().toLowerCase().indexOf(this.wakeWord) === 0) {
      this.awake = true
      const callback = this.listeners['activate']
      if (callback) {
        callback()
      }
    }
    this.recordedAudioLength += (chunk.length / 2) * (1 / 16000) * 1000
    this.modelStream.feedAudioContent(chunk)
  }

  intermediateDecode() {
    let results = this.finishStream()
    this.createStream()
    return results
  }

  bufferSilence(data) {
    // VAD has a tendency to cut the first bit of audio data from the start of a recording
    // so keep a buffer of that first bit of audio and in addBufferedSilence() reattach it to the beginning of the recording
    this.silenceBuffers.push(data)
    if (this.silenceBuffers.length >= 3) {
      this.silenceBuffers.shift()
    }
  }

  createStream() {
    this.modelStream = this.englishModel.createStream();
    this.recordedChunks = 0
    this.recordedAudioLength = 0
  }

  finishStream() {
    if (this.modelStream) {
      let start = new Date();
      let text = this.modelStream.finishStream()
      if (text) {
        let recogTime = new Date().getTime() - start.getTime()
        return {
          text,
          recogTime,
          audioLength: Math.round(this.recordedAudioLength)
        };
      }
    }
    this.silenceBuffers = []
    this.modelStream = null
  }

  on(key, callback) {
    this.listeners[key] = callback
  }
}

module.exports = Listener

const OutputText = ({text}) => {
  return preact.h('pre', {
    className: 'hal-output-text'
  }, text)
}

const isImageURL = (url) => {
  const urlObj = new URL(url)
  return urlObj.pathname.endsWith('.jpeg') || urlObj.pathname.endsWith('.jpg') || urlObj.pathname.endsWith('.gif') || urlObj.pathname.endsWith('.png')
}

const OutputMedia = ({url}) => {
  if (!url) {
    return null
  }
  if (isImageURL(url)) {
    return preact.h('div', {
      className: 'hal-output-image',
      style: {
        backgroundImage: `url(${url})`
      }
    }, null)
  } else {
    return preact.h('iframe', {
      className: 'hal-output-iframe',
      src: url,
    }, null)
  }
}

const Output = (props) => {
  return preact.h('div', {
    className: 'hal-output'
  }, [
    preact.h('div', {
      className: 'hal-output-inner hal-display'
    }, [
      props.error || (props.response && props.response.media === '') ? (
        preact.h(OutputText, {
          text: props.error ? `${props.error}` : (props.response && props.response.body)
        })
      ) : (
        preact.h(OutputMedia, {
          url: props.response && props.response.media
        })
      )
    ])
  ])
}

const Input = (props) => {
  return preact.h('div', {
    className: 'hal-input'
  }, [
    preact.h('div', {
      className: 'hal-input-inner hal-display'
    }, props.input)
  ])
}

class App extends preact.Component {
  constructor(props) {
    super(props)
    this.state = {
      listening: false,
      sending: null,
      lastSent: null,
      response: null,
      error: null
    }
  }

  componentDidMount() {
    window.hal.onListening((_, listening) => {
      this.setState({listening})
    })
    window.hal.onSending((_, sending) => {
      this.setState({
        sending,
        lastSent: sending
      })
    })
    window.hal.onResponse((_, response) => {
      this.setState({
        response,
        sending: null,
        listening: false,
        error: null
      })
    })
    window.hal.onError((_, error) => {
      this.setState({
        error,
        sending: null,
        listening: false,
      })
    })
  }

  render() {
    return preact.h('div', {
      className: [
        'hal',
        this.state.listening ? 'hal-listening' : null,
        this.state.sending ? 'hal-sending' : null,
        this.state.error ? 'hal-error' : null
      ].filter(c => !!c).join(' ')
    }, [
      preact.h(Output, {
        response: this.state.response,
        error: this.state.error
      }),
      preact.h(Input, {
        input: this.state.lastSent
      })
    ])
  }
}

preact.render(preact.h(App), document.getElementById('main'))

package core

const (
	SignalTypeToggleOn Signal = iota
	SignalTypeToggleOff
	SignalTypeEsc
)

const (
	StatusLightError StatusLight = iota
	StatusLightWorking
	StatusLightListening
)

func (l StatusLight) String() string {
	switch l {
	case StatusLightError:
		return "ERROR"
	case StatusLightWorking:
		return "WORKING"
	case StatusLightListening:
		return "LISTENING"
	}
	return ""
}

const (
	MediaTypeAudioStream = "audio_stream"
	MediaTypeImage       = "image"
)

const (
	ActionPlay = "play"
	ActionStop = "stop"
)

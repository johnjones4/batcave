package learning

import (
	"bytes"
	"io"

	"github.com/asticode/go-asticoqui"
	"github.com/cryptix/wav"
)

//TODO steps https://github.com/asticode/go-asticoqui

type VoiceTranscriberConfiguration struct {
	ModelPath string
}

type VoiceTranscriber struct {
	model *asticoqui.Model
}

func NewVoiceTranscriber(configuration VoiceTranscriberConfiguration) (*VoiceTranscriber, error) {
	model, err := asticoqui.New(configuration.ModelPath)
	if err != nil {
		return nil, err
	}
	return &VoiceTranscriber{model: model}, nil
}

func (t *VoiceTranscriber) Transcribe(mp3Bytes []byte) (string, error) {
	data, err := getRawAudio(mp3Bytes)
	if err != nil {
		return "", err
	}

	return t.model.SpeechToText(data)
}

func getRawAudio(wavBytes []byte) ([]int16, error) {
	r, err := wav.NewReader(bytes.NewReader(wavBytes), int64(len(wavBytes)))
	if err != nil {
		return nil, err
	}
	var d []int16
	for {
		s, err := r.ReadSample()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		d = append(d, int16(s))
	}
	return d, nil
}

package learning

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"

	"github.com/asticode/go-asticoqui"
	"github.com/cryptix/wav"
	"github.com/hraban/opus"
	"github.com/johnjones4/hal-9000/server/hal9000/core"
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

func (t *VoiceTranscriber) Transcribe(audio core.Audio) (string, error) {
	data, err := parseToRawAudio(audio)
	if err != nil {
		return "", err
	}

	return t.model.SpeechToText(data)
}

func parseToRawAudio(audio core.Audio) ([]int16, error) {
	audioReader := base64.NewDecoder(base64.StdEncoding, bytes.NewBuffer([]byte(audio.Data)))
	var err error
	if audio.GZipped {
		audioReader, err = gzip.NewReader(audioReader)
		if err != nil {
			return nil, err
		}
	}

	encodedBytes, err := ioutil.ReadAll(audioReader)
	if err != nil {
		return nil, err
	}

	switch audio.MimeType {
	case "audio/wav":
		return wavToRawAudio(encodedBytes)
	case "audio/ogg; codecs=opus":
		//TODO sudo apt-get install pkg-config libopus-dev libopusfile-dev
		return oggOpusToRawAudio(encodedBytes)
	}

	return nil, errors.New("unsupported codec")
}

func wavToRawAudio(wavBytes []byte) ([]int16, error) {
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

func oggOpusToRawAudio(opusBytes []byte) ([]int16, error) {
	channels := 1
	sampleRate := 16000
	dec, err := opus.NewDecoder(sampleRate, channels)
	if err != nil {
		return nil, err
	}
	frameSizeMs := 60
	frameSize := channels * frameSizeMs * sampleRate / 1000
	pcm := make([]int16, int(frameSize))
	n, err := dec.Decode(opusBytes, pcm)
	if err != nil {
		return nil, err
	}
	return pcm[:n*channels], nil
}

package learning

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"sync"

	"github.com/asticode/go-asticoqui"
	"github.com/cryptix/wav"
	"github.com/johnjones4/hal-9000/server/hal9000/core"
	"github.com/xfrr/goffmpeg/transcoder"
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
		return oggToRawAudio(encodedBytes)
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

func oggToRawAudio(oggBytes []byte) ([]int16, error) {
	var pcmData []int16
	var rErr1 error
	var rErr2 error
	trans := new(transcoder.Transcoder)
	err := trans.InitializeEmptyTranscoder()
	if err != nil {
		return nil, err
	}
	trans.MediaFile().SetSkipVideo(true)
	trans.MediaFile().SetAudioRate(16000)
	trans.MediaFile().SetAudioCodec("pcm_s16le")
	trans.MediaFile().SetAudioChannels(1)
	w, err := trans.CreateInputPipe()
	if err != nil {
		return nil, err
	}
	r, err := trans.CreateOutputPipe("s16le")

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer r.Close()
		defer wg.Done()
		pcmData = make([]int16, 0)
		for {
			buf := make([]byte, 2)
			_, err := r.Read(buf)
			if err == io.EOF {
				break
			} else if err != nil {
				rErr1 = err
				break
			}
			var sample int16 = int16(buf[0]) + int16(buf[1])<<8
			pcmData = append(pcmData, sample)
		}
	}()

	go func() {
		defer w.Close()
		_, rErr2 = w.Write(oggBytes)
	}()

	done := trans.Run(false)
	err = <-done
	if err != nil {
		return nil, err
	}

	wg.Wait()

	if rErr1 != nil {
		return nil, rErr1
	}

	if rErr2 != nil {
		return nil, rErr2
	}

	return pcmData, nil
}

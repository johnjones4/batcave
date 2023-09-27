package stt

import (
	"github.com/asticode/go-asticoqui"
	"github.com/go-audio/audio"
)

type AstiCoqui struct {
	model *asticoqui.Model
}

func NewAstiCoqui(modelPath string, scorerPath string) (*AstiCoqui, error) {
	model, err := asticoqui.New(modelPath)
	if err != nil {
		return nil, err
	}
	err = model.EnableExternalScorer(scorerPath)
	if err != nil {
		return nil, err
	}
	return &AstiCoqui{
		model: model,
	}, nil
}

func (a *AstiCoqui) SpeechToText(speech *audio.IntBuffer) (string, error) {
	ints := make([]int16, len(speech.Data))
	for i, ii := range speech.Data {
		ints[i] = int16(ii)
	}
	return a.model.SpeechToText(ints)
}

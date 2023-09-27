package stt

import (
	"github.com/go-audio/audio"
	"github.com/sirupsen/logrus"
)

type DummyStt struct {
	Log logrus.FieldLogger
}

func (d *DummyStt) SpeechToText(speech *audio.IntBuffer) (string, error) {
	d.Log.Infof("Speech is %d ffame long", speech.NumFrames())
	return "stop", nil
}

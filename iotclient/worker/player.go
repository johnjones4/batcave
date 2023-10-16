package worker

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/ebitengine/oto/v3"
	"github.com/kvark128/minimp3"
	"github.com/sirupsen/logrus"
)

type Player struct {
	workerConcrete
	player  *oto.Player
	context *oto.Context
}

func NewPlayer(log logrus.FieldLogger) *Player {
	return &Player{
		workerConcrete: newWorkerConcrete(log),
	}
}

func (p *Player) Play(ctx context.Context, url string) {
	err := p.Stop()
	if err != nil {
		p.errors <- err
		return
	}

	res, err := http.Get(url)
	if err != nil {
		p.errors <- err
		return
	}

	var reader io.Reader
	var opts oto.NewContextOptions

	switch res.Header.Get("Content-type") {
	case "audio/mpeg":
		decoder := minimp3.NewDecoder(res.Body)
		decoder.Read([]byte{})

		opts = oto.NewContextOptions{
			SampleRate:   decoder.SampleRate(),
			ChannelCount: decoder.Channels(),
			Format:       oto.FormatSignedInt16LE,
		}
		reader = decoder
	default:
		p.errors <- errors.New("usupported media")
		return
	}

	c, ready, err := oto.NewContext(&opts)
	if err != nil {
		p.errors <- err
		return
	}
	<-ready
	p.context = c
	p.player = c.NewPlayer(reader)
	p.player.Play()
}

func (p *Player) Stop() error {
	if p.player == nil || !p.player.IsPlaying() {
		return nil
	}
	p.player.Pause()
	err := p.player.Close()
	if err != nil {
		return err
	}
	return p.context.Suspend()
}

func (p *Player) Setup(errors chan error) error {
	p.errors = errors
	return nil
}

func (p *Player) Teardown() error {
	return nil
}

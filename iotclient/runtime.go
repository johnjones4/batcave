package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"main/core"
	"main/util"
	"main/worker"
	"os"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type closeable interface {
	Close() error
}

type runtime struct {
	log           logrus.FieldLogger
	controller    core.Controller
	display       core.Display
	lights        core.StatusLightsControl
	voiceWorker   *worker.VoiceWorker
	commandWorker *worker.CommandWorker
	player        *worker.Player
	cfg           util.ServerConfig
	pendingEvents map[string]bool
}

func (r *runtime) start(ctx context.Context) error {
	r.log.Debug("Starting up")
	r.pendingEvents = make(map[string]bool)
	errs := make(chan error, 128)
	workers := []core.Worker{
		r.voiceWorker,
		r.commandWorker,
		r.player,
	}
	closables := []closeable{
		r.controller,
		r.display,
		r.lights,
	}

	r.log.Debug("Initializing workers")
	for _, w := range workers {
		err := w.Setup(errs)
		if err != nil {
			return err
		}
	}

	r.log.Debug("Starting workers")
	go r.controller.Start(ctx)
	go r.commandWorker.Start(ctx)

	r.log.Debug("Beginning event loop")
	go func() {
		for {
			var terminate bool
			var err error
			select {
			case <-ctx.Done():
				terminate = true
			case signal := <-r.controller.SignalChannel():
				terminate, err = r.handleControlSignal(ctx, signal)
			case voiceData := <-r.voiceWorker.Chan():
				err = r.handleVoiceData(ctx, voiceData)
			case res := <-r.commandWorker.Chan():
				err = r.handleReponse(ctx, res)
			case err1 := <-errs:
				err = err1
			}
			if err != nil {
				r.log.Errorf("runtime error: %s", err)
				r.lights.SetModeStatusLight(ctx, core.StatusLightError, true)
			}
			if terminate {
				for _, c := range closables {
					err := c.Close()
					if err != nil {
						r.log.Errorf("close error: %s", err)
					}
				}
				for _, w := range workers {
					err = w.Teardown()
					if err != nil {
						r.log.Errorf("teardown error: %s", err)
					}
				}
				os.Exit(0)
			}
		}
	}()

	r.log.Debug("Beginning main application loop")
	r.display.Start(ctx)

	return nil
}

func (r *runtime) handleControlSignal(ctx context.Context, signal core.Signal) (bool, error) {
	r.log.Debugf("Got signal: %d", signal)
	switch signal {
	case core.SignalTypeEsc:
		return true, nil
	case core.SignalTypeToggleOff:
		go r.voiceWorker.Stop()
		return false, r.lights.SetModeStatusLight(ctx, core.StatusLightListening, false)
	case core.SignalTypeToggleOn:
		go r.voiceWorker.Start(ctx)
		return false, r.lights.SetModeStatusLight(ctx, core.StatusLightListening, true)
	}
	return false, nil
}

func (r *runtime) handleVoiceData(ctx context.Context, voiceData []byte) error {
	r.log.Debugf("Got %d bytes of audio to send", len(voiceData))
	var req core.Request
	req.ClientID = r.cfg.ClientId
	req.EventId = uuid.NewString()
	req.Message.Audio.Data = base64.StdEncoding.EncodeToString(voiceData)
	r.commandWorker.SendChan() <- req
	r.pendingEvents[req.EventId] = true
	return r.lights.SetModeStatusLight(ctx, core.StatusLightWorking, true)
}

func (r *runtime) handleReponse(ctx context.Context, res core.Response) error {
	var responseBody *core.ResponseBody
	if res.Response != nil {
		responseBody = res.Response
		delete(r.pendingEvents, res.Response.EventId)
		if len(r.pendingEvents) == 0 {
			err := r.lights.SetModeStatusLight(ctx, core.StatusLightWorking, false)
			if err != nil {
				return err
			}
		}
	} else if res.PushMessage != nil {
		responseBody = res.PushMessage
	}
	if responseBody != nil {
		if responseBody.Message.Audio.Data != "" {
			data, err := base64.StdEncoding.DecodeString(responseBody.Message.Audio.Data)
			if err != nil {
				return err
			}
			r.player.PlayBuffer(ctx, "audio/wav", bytes.NewReader(data))
		}

		switch responseBody.Action {
		case core.ActionStop:
			err := r.player.Stop()
			if err != nil {
				return err
			}
		case core.ActionPlay:
			switch responseBody.Media.Type {
			case core.MediaTypeAudioStream:
				go r.player.PlayURL(ctx, responseBody.Media.URL)
			}
		}
	}
	r.display.Display(ctx, res)
	return nil
}

package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"main/core"
	"main/util"
	"main/worker"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type runtime struct {
	log           logrus.FieldLogger
	controller    core.Controller
	display       core.Display
	lights        core.StatusLightsControl
	voiceWorker   *worker.VoiceWorker
	commandWorker *worker.CommandWorker
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
	}

	r.log.Debug("Initializing workers")
	for _, w := range workers {
		err := w.Setup(errs)
		if err != nil {
			return err
		}
	}
	defer func() {
		err := r.controller.Close()
		if err != nil {
			r.log.Errorf("controller close error: %s", err)
		}

		err = r.display.Close()
		if err != nil {
			r.log.Errorf("display close error: %s", err)
		}

		err = r.lights.Close()
		if err != nil {
			r.log.Errorf("lights close error: %s", err)
		}

		for _, w := range workers {
			err = w.Teardown()
			if err != nil {
				r.log.Errorf("teardown error: %s", err)
			}
		}
	}()

	r.log.Debug("Starting workers")
	go r.controller.Start(ctx)
	go r.commandWorker.Start(ctx)

	r.log.Debug("Beginning main application loop")
	for {
		var terminate bool
		var err error
		select {
		case <-ctx.Done():
			return ctx.Err()
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
			return nil
		}
	}
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
	switch res.Type {
	case "request":
		return r.display.Write(ctx, fmt.Sprintf("You: %s", res.Request.Message.Text))
	case "response", "push":
		var resBody core.ResponseBody
		if res.Response != nil {
			resBody = *res.Response
		} else if res.PushMessage != nil {
			resBody = *res.PushMessage
		}
		if resBody.EventId != "" {
			delete(r.pendingEvents, resBody.EventId)
			if len(r.pendingEvents) == 0 {
				err := r.lights.SetModeStatusLight(ctx, core.StatusLightWorking, false)
				if err != nil {
					return err
				}
			}

			return r.display.Write(ctx, fmt.Sprintf("HAL: %s", resBody.Message.Text))
		}
	}
	return nil
}

package mode

import (
	"context"
	"encoding/base64"
	"fmt"
	"main/core"
	"main/worker"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ModeCommand struct {
	modeConcrete
	voiceWorker   *worker.VoiceWorker
	commandWorker *worker.CommandWorker
	on            bool
	clientId      string
}

func NewModeCommand(log logrus.FieldLogger, voiceWorker *worker.VoiceWorker, commandWorker *worker.CommandWorker, clientId string) *ModeCommand {
	return &ModeCommand{
		modeConcrete:  newModeConcrete(log),
		voiceWorker:   voiceWorker,
		commandWorker: commandWorker,
		clientId:      clientId,
	}
}

// TODO status lights
func (m *ModeCommand) Start(ctx context.Context, displayCtx core.DisplayContext) {
	m.log.Debug("Starting command mode")
	go m.commandWorker.Start(ctx)
	defer func() {
		m.stopped <- true
	}()
	defer displayCtx.Close()
	for {
		select {
		case <-ctx.Done():
			return
		case <-m.stop:
			return
		case voiceData := <-m.voiceWorker.Chan():
			m.log.Debugf("Got %d bytes of audio to send", len(voiceData))
			var req core.Request
			req.ClientID = m.clientId
			req.EventId = uuid.NewString()
			req.Message.Audio.Data = base64.StdEncoding.EncodeToString(voiceData)
			m.commandWorker.SendChan() <- req
		case res := <-m.commandWorker.Chan():
			switch res.Type {
			case "request":
				displayCtx.Write(ctx, fmt.Sprintf("You: %s", res.Request.Message.Text))
			case "response", "push":
				var resBody core.ResponseBody
				if res.Response != nil {
					resBody = *res.Response
				} else if res.PushMessage != nil {
					resBody = *res.PushMessage
				}
				if resBody.EventId != "" {
					displayCtx.Write(ctx, fmt.Sprintf("HAL: %s", resBody.Message.Text))
				}
			}
		}
	}
}

func (m *ModeCommand) Stop() error {
	m.stop <- true
	<-m.stopped
	return m.commandWorker.Stop()
}

func (m *ModeCommand) Toggle(ctx context.Context) error {
	if !m.on {
		go m.voiceWorker.Start(ctx)
		m.on = true
	} else {
		m.on = false
		return m.voiceWorker.Stop()
	}
	return nil
}

package worker

import (
	"context"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

const (
	Channels   = 1
	SampleRate = 16000
	BitDepth   = 16
)

type VoiceWorker struct {
	workerConcrete
	queue  chan []byte
	cancel context.CancelFunc
}

func NewVoiceWorker(log logrus.FieldLogger) *VoiceWorker {
	return &VoiceWorker{
		workerConcrete: newWorkerConcrete(log),
		queue:          make(chan []byte),
	}
}

func (v *VoiceWorker) Setup(errors chan error) error {
	v.workerConcrete.errors = errors
	return nil
}

func (v *VoiceWorker) Teardown() error {
	return nil
}

func (v *VoiceWorker) Stop() {
	v.log.Debug("Stopping voice")
	if v.cancel != nil {
		v.cancel()
	}
	go func() {
		file, err := os.ReadFile("/tmp/hal.wav")
		if err != nil && err != context.Canceled {
			v.errors <- err
			return
		}
		v.Chan() <- file
	}()
	v.log.Debug("Stopped voice")
}

func (v *VoiceWorker) Start(ctx context.Context) {
	v.log.Debug("Starting voice")
	cancellable, cancel := context.WithCancel(ctx)
	v.cancel = cancel
	cmd := exec.CommandContext(cancellable, "sox", "-d", "/tmp/hal.wav")
	err := cmd.Run()
	if err != nil && err != context.Canceled {
		v.errors <- err
		return
	}
}

func (v *VoiceWorker) Chan() chan []byte {
	return v.queue
}

package worker

import (
	"context"
	"errors"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

type BufferPlayer struct {
	workerConcrete
	cancel context.CancelFunc
}

func NewBufferPlayer(log logrus.FieldLogger) *BufferPlayer {
	return &BufferPlayer{
		workerConcrete: newWorkerConcrete(log),
	}
}

func (p *BufferPlayer) PlayWAV(ctx context.Context, bytes []byte) {
	p.log.Debug("Starting player")
	err := os.WriteFile("/tmp/output.wav", bytes, 0777)
	if err != nil {
		p.errors <- err
		return
	}
	cancellable, cancel := context.WithCancel(ctx)
	p.cancel = cancel
	cmd := exec.CommandContext(cancellable, "play", "/tmp/output.wav")
	err = cmd.Run()
	if err != nil && !errors.As(err, &exitError) {
		p.errors <- err
		return
	}
}

func (p *BufferPlayer) Stop() error {
	p.log.Debug("Stopping player")
	if p.cancel != nil {
		p.cancel()
	}
	return nil
}

func (p *BufferPlayer) Setup(errors chan error) error {
	p.errors = errors
	return nil
}

func (p *BufferPlayer) Teardown() error {
	return nil
}

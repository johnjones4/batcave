package mode

import (
	"context"
	"main/core"
	"main/worker"

	"github.com/sirupsen/logrus"
)

type ModeLog struct {
	modeConcrete
	logWorker *worker.LogWorker
}

func NewModeLog(log logrus.FieldLogger, logWorker *worker.LogWorker) *ModeLog {
	return &ModeLog{
		modeConcrete: newModeConcrete(log),
		logWorker:    logWorker,
	}
}

// TODO status lights
func (m *ModeLog) Start(ctx context.Context, displayCtx core.DisplayContext) {
	m.log.Debug("Starting log mode")
	go m.logWorker.Start(ctx)
	defer func() {
		displayCtx.Close()
		m.stopped <- true
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case <-m.stop:
			return
		case msg := <-m.logWorker.Chan():
			displayCtx.Write(ctx, msg)
		}
	}
}

func (m *ModeLog) Stop() error {
	m.stop <- true
	<-m.stopped
	return m.logWorker.Stop()
}

func (m *ModeLog) Toggle(ctx context.Context) error {
	return nil
}

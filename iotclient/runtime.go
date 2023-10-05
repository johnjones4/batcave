package main

import (
	"context"
	"main/core"

	"github.com/sirupsen/logrus"
)

type Runtime struct {
	log                   logrus.FieldLogger
	modes                 []core.Mode
	controller            core.Controller
	displayContextFactory core.DisplayContextFactory
	activeMode            int
}

func New(log logrus.FieldLogger, modes []core.Mode, controller core.Controller, displayContextFactory core.DisplayContextFactory) *Runtime {
	return &Runtime{
		activeMode:            0,
		log:                   log,
		modes:                 modes,
		controller:            controller,
		displayContextFactory: displayContextFactory,
	}
}

func (r *Runtime) startNextMode(ctx context.Context) {
	dctx, err := r.displayContextFactory()
	if err != nil {
		r.log.Errorf("Error getting display context: %s", err)
		return
	}
	go r.modes[r.activeMode].Start(ctx, dctx)
}

func (r *Runtime) Start(ctx context.Context) {
	go r.controller.Start(ctx)
	r.startNextMode(ctx)
	for event := range r.controller.SignalChannel() {
		switch event.Type {
		case core.SignalTypeMode:
			if event.Mode < len(r.modes) {
				err := r.modes[r.activeMode].Stop()
				if err != nil {
					r.log.Errorf("Error stopping current mode: %s", err)
				}
				r.activeMode = event.Mode
				r.startNextMode(ctx)
			}
		case core.SignalTypeToggle:
			err := r.modes[r.activeMode].Toggle(ctx)
			if err != nil {
				r.log.Errorf("Error sending toggle to mode: %s", err)
			}
		case core.SignalTypeEsc:
			r.log.Debug("Stopping all")
			err := r.modes[r.activeMode].Stop()
			if err != nil {
				r.log.Errorf("Error stopping current mode: %s", err)
			}
			return
		}
	}
}

package api

import (
	"context"
	"main/core"
)

func (a *API) bundledHandler(ctx context.Context, req *core.Request) (core.Response, error) {
	err := a.prepareRequest(ctx, req)
	if err != nil {
		return core.ResponseEmpty, err
	}

	res, err := a.coreHandler(ctx, req)
	if err != nil {
		return core.ResponseEmpty, err
	}

	return res, nil
}

func (a *API) prepareRequest(ctx context.Context, req *core.Request) error {
	for _, proc := range a.RequestProcessors {
		err := proc(ctx, req)
		if err != nil {
			return err
		}
	}
	//TODO undo request insert on failure
	return nil
}

func (a *API) coreHandler(ctx context.Context, req *core.Request) (core.Response, error) {
	if req.Message.Text == "" {
		return core.ResponseEmpty, nil
	}

	i, md, err := a.IntentMatcher.Match(ctx, req)
	if err != nil {
		return core.Response{}, err
	}

	res, err := i.ActOnIntent(ctx, req, &md)
	if err != nil {
		return core.Response{}, err
	}

	res.EventId = req.EventId

	for _, proc := range a.ResponseProcessors {
		err = proc(ctx, req, &res)
		if err != nil {
			return core.Response{}, err
		}
	}

	return res, nil
}

package api

import (
	"context"
	"main/core"
)

func (a *apiConcrete) coreHandler(ctx context.Context, req core.Request) (core.Response, error) {
	for _, proc := range a.RequestProcessors {
		err := proc(ctx, &req)
		if err != nil {
			return core.Response{}, err
		}
	}

	i, md, err := a.IntentMatcher.Match(ctx, &req)
	if err != nil {
		return core.Response{}, err
	}

	res, err := i.ActOnIntent(ctx, &req, &md)
	if err != nil {
		return core.Response{}, err
	}

	res.EventId = req.EventId

	for _, proc := range a.ResponseProcessors {
		err = proc(ctx, &req, &res)
		if err != nil {
			return core.Response{}, err
		}
	}

	return res, nil
}

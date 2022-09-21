package core

import (
	"context"
	"errors"
)

type ContextKey string

const (
	coordinateKey ContextKey = "coordinate"
)

func ContextWithCoordinates(ctx context.Context, c Coordinate) context.Context {
	return context.WithValue(ctx, coordinateKey, c)
}

func CoordinatesInContext(ctx context.Context) (Coordinate, error) {
	coord, ok := ctx.Value(coordinateKey).(Coordinate)
	if !ok {
		return Coordinate{}, errors.New("no coordinates in context")
	}
	return coord, nil
}

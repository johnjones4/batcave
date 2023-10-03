package pgstore

import (
	"errors"
	"main/mocks"

	"go.uber.org/mock/gomock"
)

var errorTestError = errors.New("test")

func newStore(ctrl *gomock.Controller) *PGStore {
	return &PGStore{
		pool: mocks.NewMockDatabase(ctrl),
		log:  mocks.NewMockFieldLogger(ctrl),
	}
}

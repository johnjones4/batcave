package api

import (
	"errors"
	"main/core"
	"main/mocks"

	"go.uber.org/mock/gomock"
)

var (
	errorTestError = errors.New("test error")
)

func mockAPI(ctrl *gomock.Controller) *API {
	intentMatcher := mocks.NewMockIntentMatcher(ctrl)
	requestProcessors := []core.RequestProcessor{}
	responseProcessors := []core.ResponseProcessor{}
	logger := mocks.NewMockHookableLogger(ctrl)
	telegram := mocks.NewMockTelegramSender(ctrl)
	clientRegistry := mocks.NewMockClientRegistry(ctrl)
	socketSender := mocks.NewMockSocketSender(ctrl)

	logger.EXPECT().AddHook(gomock.Any()).Times(1)

	return New(APIParams{
		IntentMatcher:      intentMatcher,
		RequestProcessors:  requestProcessors,
		ResponseProcessors: responseProcessors,
		Log:                logger,
		Telegram:           telegram,
		ClientRegistry:     clientRegistry,
		SocketSender:       socketSender,
	})
}

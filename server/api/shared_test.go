package api

import (
	"encoding/json"
	"errors"
	"main/core"
	"main/mocks"
	"net/http"

	"go.uber.org/mock/gomock"
)

var (
	testError = errors.New("test error")
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

type requestMatcher struct {
	a *http.Request
}

func newRequestMatcher(a *http.Request) *requestMatcher {
	return &requestMatcher{a}
}

func (m *requestMatcher) Matches(x interface{}) bool {
	b := x.(*http.Request)
	if m.a.URL.String() != b.URL.String() || m.a.Method != b.Method || len(m.a.Header) != len(b.Header) {
		return false
	}
	for k, v := range m.a.Header {
		if b.Header.Get(k) != v[0] {
			return false
		}
	}
	return true
}

func (m *requestMatcher) String() string {
	return "http request matcher"
}

type jsonMatcher struct {
	a string
}

func newJsonMatcher(a any) *jsonMatcher {
	ab, _ := json.Marshal(a)
	return &jsonMatcher{string(ab)}
}

func (m *jsonMatcher) Matches(b interface{}) bool {
	bb, _ := json.Marshal(b)
	return m.a == string(bb)
}

func (m *jsonMatcher) String() string {
	return m.a
}

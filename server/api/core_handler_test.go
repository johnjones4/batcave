package api

import (
	"context"
	"main/core"
	"main/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type prepareRequestTest struct {
	err error
}

func (t *prepareRequestTest) Process(ctx context.Context, req *core.Request) error {
	return t.err
}

func TestPrepareRequest(t *testing.T) {
	cases := []prepareRequestTest{
		{
			err: errorTestError,
		},
		{
			err: nil,
		},
	}

	for _, c := range cases {
		a := &API{
			APIParams: APIParams{
				RequestProcessors: []core.RequestProcessor{
					c.Process,
				},
			},
		}
		err := a.prepareRequest(context.Background(), &core.Request{})
		assert.Equal(t, c.err, err)
	}
}

type coreHandlerTest struct {
	request                    core.Request
	response                   core.Response
	intentMatcherError         error
	intentActorError           error
	responseProcessorError     error
	responseProcessorDidRun    bool
	responseProcessorShouldRun bool
}

func (t *coreHandlerTest) Process(ctx context.Context, req *core.Request, resp *core.Response) error {
	t.responseProcessorDidRun = true
	return t.responseProcessorError
}

func TestCoreHandler(t *testing.T) {
	ctrl := gomock.NewController(t)

	api := mockAPI(ctrl)
	intentMatcher := api.IntentMatcher.(*mocks.MockIntentMatcher)

	cases := []coreHandlerTest{
		{
			request: core.Request{
				EventId: "event id",
				Message: core.Message{
					Text: "hello",
				},
				ClientID: "123",
			},
			response: core.Response{
				OutboundMessage: core.OutboundMessage{
					EventId: "event id",
					Message: core.Message{
						Text: "world",
					},
				},
			},
			responseProcessorShouldRun: true,
		},
		{
			request: core.Request{
				EventId: "event id",
				Message: core.Message{
					Text: "hello",
				},
				ClientID: "123",
			},
			intentMatcherError: errorTestError,
		},
		{
			request: core.Request{
				EventId: "event id",
				Message: core.Message{
					Text: "hello",
				},
				ClientID: "123",
			},
			intentActorError: errorTestError,
		},
		{
			request: core.Request{
				EventId: "event id",
				Message: core.Message{
					Text: "hello",
				},
				ClientID: "123",
			},
			responseProcessorError:     errorTestError,
			responseProcessorShouldRun: true,
		},
	}

	for _, c := range cases {
		api.ResponseProcessors = []core.ResponseProcessor{c.Process}

		actor := mocks.NewMockIntentActor(ctrl)

		if c.request.Message.Text != "" {
			metadata := core.IntentMetadata{
				IntentParseReceiver: map[string]any{},
			}
			imCall := intentMatcher.EXPECT().Match(gomock.Any(), gomock.Any())
			if c.intentMatcherError == nil {
				imCall.Return(actor, metadata, nil)

				iaCall := actor.EXPECT().ActOnIntent(gomock.Any(), gomock.Any(), gomock.Any())
				if c.intentActorError == nil {
					iaCall.Return(c.response, nil)
				} else {
					iaCall.Return(core.Response{}, c.intentActorError)
				}
			} else {
				imCall.Return(nil, core.IntentMetadata{}, c.intentMatcherError)
			}
		}

		response, err := api.coreHandler(context.Background(), &c.request)

		if c.intentMatcherError != nil {
			assert.Equal(t, c.intentMatcherError, err)
		} else if c.intentActorError != nil {
			assert.Equal(t, c.intentActorError, err)
		} else if c.responseProcessorError != nil {
			assert.Equal(t, c.responseProcessorError, err)
		}

		if err == nil {
			assert.Equal(t, c.response.EventId, c.request.EventId)
		}
		assert.Equal(t, c.response.Message.Text, response.Message.Text)
		assert.Equal(t, c.response.Message.Audio.Data, response.Message.Audio.Data)
		assert.Equal(t, c.response.Media.Type, response.Media.Type)
		assert.Equal(t, c.response.Media.URL, response.Media.URL)

		assert.Equal(t, c.responseProcessorShouldRun, c.responseProcessorDidRun)
	}
}

package api

import (
	"bytes"
	"encoding/json"
	"main/core"
	"main/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type messageTestCase struct {
	request             core.Request
	response            core.Response
	breakBeforeHandling bool
	clientId            string
	responseCode        int
	clientRegistryError error
}

func TestMessage(t *testing.T) {
	ctrl := gomock.NewController(t)

	api := mockAPI(ctrl)
	clientRegistryMock := api.ClientRegistry.(*mocks.MockClientRegistry)
	intentMatcher := api.IntentMatcher.(*mocks.MockIntentMatcher)
	log := api.Log.(*mocks.MockHookableLogger)

	cases := []messageTestCase{
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
			clientId:     "123",
			responseCode: http.StatusOK,
		},
		{
			request: core.Request{
				EventId: "event id",
				Message: core.Message{
					Text: "hello",
				},
				ClientID: "123",
			},
			response:            core.ResponseEmpty,
			clientId:            "123",
			responseCode:        http.StatusInternalServerError,
			clientRegistryError: testError,
		},
		{
			request: core.Request{
				EventId: "event id",
				Message: core.Message{
					Text: "hello",
				},
				ClientID: "",
			},
			response:            core.ResponseEmpty,
			clientId:            "123",
			responseCode:        http.StatusBadRequest,
			breakBeforeHandling: true,
		},
		{
			request: core.Request{
				EventId: "event id",
				Message: core.Message{
					Text: "hello",
				},
				ClientID: "fadfssdfsdf",
			},
			response:            core.ResponseEmpty,
			clientId:            "123",
			responseCode:        http.StatusBadRequest,
			breakBeforeHandling: true,
		},
		{
			request: core.Request{
				EventId: "",
				Message: core.Message{
					Text: "hello",
				},
				ClientID: "123",
			},
			response:            core.ResponseEmpty,
			clientId:            "123",
			responseCode:        http.StatusBadRequest,
			breakBeforeHandling: true,
		},
	}

	clientKey := "456"

	for _, c := range cases {
		if c.responseCode != http.StatusOK {
			log.EXPECT().Error(gomock.Any())
		}

		body, _ := json.Marshal(c.request)

		req, _ := http.NewRequest("POST", "/api/client/message", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Client-Id", c.clientId)
		req.Header.Set("X-API-Key", clientKey)

		crCall := clientRegistryMock.EXPECT().Client(gomock.Any(), "api", c.clientId, gomock.Any())
		if c.clientRegistryError == nil {
			crCall.Return(core.Client{
				Id: c.clientId,
				Info: directAPIClientInfo{
					Key: clientKey,
				},
			}, nil)

			if !c.breakBeforeHandling {
				actor := mocks.NewMockIntentActor(ctrl)
				intentMatcher.EXPECT().Match(gomock.Any(), gomock.Any()).Return(actor, core.IntentMetadata{}, nil)
				actor.EXPECT().ActOnIntent(gomock.Any(), gomock.Any(), gomock.Any()).Return(c.response, nil)
			}
		} else {
			crCall.Return(core.Client{}, c.clientRegistryError)
		}

		log.EXPECT().Print(gomock.Any())

		recorder := httptest.NewRecorder()

		api.mux.ServeHTTP(recorder, req)

		assert.Equal(t, c.responseCode, recorder.Code)

		if c.responseCode == http.StatusOK {
			var resp core.Response
			json.Unmarshal(recorder.Body.Bytes(), &resp)

			assert.Equal(t, c.response.EventId, c.request.EventId)
			assert.Equal(t, c.response.Message.Text, resp.Message.Text)
			assert.Equal(t, c.response.Message.Audio.Data, resp.Message.Audio.Data)
			assert.Equal(t, c.response.Media.Type, resp.Media.Type)
			assert.Equal(t, c.response.Media.URL, resp.Media.URL)
		}

	}
}

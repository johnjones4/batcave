package api

import (
	"bytes"
	"encoding/json"
	"main/core"
	"main/mocks"
	"main/services/telegram"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type telegramTestCase struct {
	update                 telegram.Update
	response               core.Response
	IsClientPermittedError error
	IsClientPermittedOk    bool
	responseCode           int
}

func TestTelegram(t *testing.T) {
	ctrl := gomock.NewController(t)

	api := mockAPI(ctrl)
	telegramSender := api.Telegram.(*mocks.MockTelegramSender)
	intentMatcher := api.IntentMatcher.(*mocks.MockIntentMatcher)
	log := api.Log.(*mocks.MockHookableLogger)

	cases := []telegramTestCase{
		{
			update: telegram.Update{
				UpdateId: 1,
				Message: telegram.IncomingMessage{
					Chat: telegram.Chat{
						Id:   1,
						Type: "private",
					},
					From: telegram.User{
						Id: 1,
					},
					Message: telegram.Message{
						Text: "hello",
					},
				},
			},
			response: core.Response{
				OutboundMessage: core.OutboundMessage{
					Message: core.Message{
						Text: "world",
					},
				},
			},
			IsClientPermittedOk: true,
			responseCode:        http.StatusOK,
		},
		{
			update: telegram.Update{
				UpdateId: 1,
				Message: telegram.IncomingMessage{
					Chat: telegram.Chat{
						Id:   1,
						Type: "private",
					},
					From: telegram.User{
						Id: 1,
					},
					Message: telegram.Message{
						Text: "hello",
					},
				},
			},
			IsClientPermittedOk: true,
			responseCode:        http.StatusOK,
		},
		{
			update: telegram.Update{
				UpdateId: 1,
				Message: telegram.IncomingMessage{
					Chat: telegram.Chat{
						Id:   1,
						Type: "private",
					},
					From: telegram.User{
						Id: 1,
					},
					Message: telegram.Message{
						Text: "hello",
					},
				},
			},
			IsClientPermittedError: errorTestError,
			responseCode:           http.StatusUnauthorized,
		},
		{
			update: telegram.Update{
				UpdateId: 1,
				Message: telegram.IncomingMessage{
					Chat: telegram.Chat{
						Id:   1,
						Type: "private",
					},
					From: telegram.User{
						Id: 1,
					},
					Message: telegram.Message{
						Text: "hello",
					},
				},
			},
			responseCode: http.StatusOK,
		},
	}

	for _, c := range cases {
		log.EXPECT().Debug(gomock.Any())
		log.EXPECT().Print(gomock.Any())
		if c.IsClientPermittedError != nil {
			log.EXPECT().Error(gomock.Any())
		}

		reqBody, _ := json.Marshal(c.update)
		req, _ := http.NewRequest("POST", "/api/telegram", bytes.NewBuffer(reqBody))

		telegramSender.EXPECT().IsClientPermitted(gomock.Any(), gomock.AssignableToTypeOf(&http.Request{}), c.update.Message.From.Id, c.update.Message.Message.Text, c.update.Message.Chat.Type).Return(c.IsClientPermittedOk, c.IsClientPermittedError)

		if c.IsClientPermittedError == nil && c.IsClientPermittedOk {
			actor := mocks.NewMockIntentActor(ctrl)
			intentMatcher.EXPECT().Match(gomock.Any(), gomock.Any()).Return(actor, core.IntentMetadata{}, nil)
			actor.EXPECT().ActOnIntent(gomock.Any(), gomock.Any(), gomock.Any()).Return(c.response, nil)
			if c.response.Message.Text != "" {
				telegramSender.EXPECT().SendOutbound(gomock.Any(), c.update.Message.Chat.Id, gomock.AssignableToTypeOf(core.OutboundMessage{}))
			}
		}

		recorder := httptest.NewRecorder()
		api.mux.ServeHTTP(recorder, req)
		assert.Equal(t, c.responseCode, recorder.Code)
	}
}

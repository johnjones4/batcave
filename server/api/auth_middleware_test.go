package api

import (
	"main/core"
	"main/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type authMiddlewareTestCase struct {
	clientId        string
	clientKey       string
	headerClientId  string
	headerClientKey string
	response        int
}

func TestAuthMiddlewareSuccess(t *testing.T) {
	api := mockAPI(gomock.NewController(t))

	cases := []authMiddlewareTestCase{
		{
			clientId:        "clientId",
			clientKey:       "clientKey",
			headerClientId:  "clientId",
			headerClientKey: "clientKey",
			response:        http.StatusOK,
		},
		{
			clientId:        "",
			clientKey:       "",
			headerClientId:  "",
			headerClientKey: "",
			response:        http.StatusUnauthorized,
		},
		{
			clientId:        "clientId",
			clientKey:       "",
			headerClientId:  "clientId",
			headerClientKey: "",
			response:        http.StatusUnauthorized,
		},
		{
			clientId:        "clientId",
			clientKey:       "clientKey",
			headerClientId:  "clientId",
			headerClientKey: "",
			response:        http.StatusUnauthorized,
		},
		{
			clientId:        "clientId",
			clientKey:       "",
			headerClientId:  "clientId",
			headerClientKey: "clientKey",
			response:        http.StatusUnauthorized,
		},
		{
			clientId:        "clientId1",
			clientKey:       "clientKey",
			headerClientId:  "clientId2",
			headerClientKey: "clientKey",
			response:        http.StatusUnauthorized,
		},
		{
			clientId:        "clientId1",
			clientKey:       "clientKey1",
			headerClientId:  "clientId2",
			headerClientKey: "clientKey2",
			response:        http.StatusUnauthorized,
		},
		{
			clientId:        "clientId",
			clientKey:       "clientKey1",
			headerClientId:  "clientId",
			headerClientKey: "clientKey2",
			response:        http.StatusUnauthorized,
		},
	}

	for _, c := range cases {
		clientRegistryMock := api.ClientRegistry.(*mocks.MockClientRegistry)
		log := api.Log.(*mocks.MockHookableLogger)

		req, _ := http.NewRequest("GET", "/api/client", nil)
		req.Header.Set("X-Client-Id", c.headerClientId)
		req.Header.Set("X-Api-Key", c.headerClientKey)

		if c.headerClientId != "" {
			clientRegistryMock.EXPECT().Client(gomock.Any(), "api", c.headerClientId, gomock.Any()).Return(core.Client{
				Id: c.clientId,
				Info: directAPIClientInfo{
					Key: c.clientKey,
				},
			}, nil)
		}

		recorder := httptest.NewRecorder()

		log.EXPECT().Print(gomock.Any())

		if c.response != http.StatusOK {
			log.EXPECT().Error(gomock.Any())
		}

		api.mux.ServeHTTP(recorder, req)

		assert.Equal(t, c.response, recorder.Code)
	}
}

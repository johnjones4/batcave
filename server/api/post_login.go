package api

import (
	"context"
	"encoding/json"
	"main/storage"
	"net/http"

	"github.com/swaggest/usecase"
)

func makeLoginHandler(userStore *storage.UserStore, tm *TokenManager) usecase.Interactor {
	type loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	return usecase.NewIOI(new(loginRequest), new(token), func(ctx context.Context, input, output interface{}) error {
		in := input.(*loginRequest)

		user, err := userStore.Login(in.Username, in.Password)
		if err != nil {
			handleError(w, err, http.StatusForbidden)
			return
		}

		token, err := tm.NewToken(user)
		if err != nil {
			handleInternalError(w, err)
			return
		}

		responseBody, err := json.Marshal(token)
		if err != nil {
			handleInternalError(w, err)
			return
		}
	})
}

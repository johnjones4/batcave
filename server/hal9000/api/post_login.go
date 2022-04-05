package api

import (
	"context"

	"github.com/johnjones4/hal-9000/server/hal9000/core"
	"github.com/johnjones4/hal-9000/server/hal9000/security"
	"github.com/johnjones4/hal-9000/server/hal9000/storage"

	"github.com/swaggest/usecase"
)

func makeLoginHandler(userStore *storage.UserStore, tm *security.TokenManager) usecase.Interactor {
	return usecase.NewIOI(new(core.LoginRequest), new(core.Token), func(ctx context.Context, input, output interface{}) error {
		in := input.(*core.LoginRequest)

		user, err := userStore.Login(in.Username, in.Password)
		if err != nil {
			return wrappedError(err, core.ErrorCodeUsernamePassword)
		}

		t, err := tm.NewToken(user)
		if err != nil {
			return wrappedError(err, core.ErrorCodeToken)
		}

		out := output.(*core.Token)
		*out = t

		return nil
	})
}

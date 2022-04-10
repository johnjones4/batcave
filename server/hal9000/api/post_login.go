package api

import (
	"context"

	"github.com/johnjones4/hal-9000/server/hal9000/runtime"

	"github.com/johnjones4/hal-9000/server/hal9000/core"

	"github.com/swaggest/usecase"
)

func makeLoginHandler(r *runtime.Runtime) usecase.Interactor {
	return usecase.NewIOI(new(core.LoginRequest), new(core.Token), func(ctx context.Context, input, output interface{}) error {
		in := input.(*core.LoginRequest)

		user, err := r.UserStore.Login(in.Username, in.Password)
		if err != nil {
			return wrappedError(err, core.ErrorCodeUsernamePassword)
		}

		t, err := r.TokenManager.NewToken(user)
		if err != nil {
			return wrappedError(err, core.ErrorCodeToken)
		}

		out := output.(*core.Token)
		*out = t

		return nil
	})
}

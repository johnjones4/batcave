package pgstore

import (
	"context"
	"main/core"
)

func (s *PGStore) UserForClient(ctx context.Context, source, clientId string) (core.User, error) {
	var user core.User
	err := s.pool.QueryRow(ctx, "SELECT user_id FROM clients_registry WHERE source = $1 AND client_id = $2", source, clientId).Scan(&user.Id)
	if err != nil {
		return core.User{}, err
	}

	err = s.pool.QueryRow(ctx, "SELECT name FROM users_registry WHERE user_id = $1", user.Id).Scan(&user.Name)
	if err != nil {
		return core.User{}, err
	}

	return user, nil
}

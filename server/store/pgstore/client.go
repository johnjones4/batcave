package pgstore

import (
	"context"
	"main/core"

	"github.com/jackc/pgx/v4"
)

func (s *PGStore) Client(ctx context.Context, source, clientId string, infoParser func(client *core.Client, info string) error) (core.Client, error) {
	var info string
	var client core.Client
	err := s.pool.QueryRow(ctx, "SELECT source, client_id, latitude, longitude, info FROM clients_registry WHERE source = $1 AND client_id = $2", source, clientId).Scan(&client.Source, &client.Id, &client.DefaultLocation.Latitude, &client.DefaultLocation.Longitude, &info)
	if err == pgx.ErrNoRows {
		return core.Client{}, nil
	}
	if err != nil {
		return core.Client{}, err
	}
	if infoParser != nil {
		err = infoParser(&client, info)
		if err != nil {
			return core.Client{}, err
		}
	}
	return client, nil
}

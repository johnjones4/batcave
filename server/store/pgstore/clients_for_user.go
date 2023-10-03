package pgstore

import (
	"context"
	"main/core"
)

func (s *PGStore) ClientsForUser(ctx context.Context, userId string, infoParser func(client *core.Client, info string) error) ([]core.Client, error) {
	rows, err := s.pool.Query(ctx, "SELECT source, client_id, user_id, latitude, longitude, info FROM clients_registry WHERE user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	clients := make([]core.Client, 0)

	for rows.Next() {
		var info string
		var client core.Client
		err = rows.Scan(&client.Source, &client.Id, &client.UserId, &client.DefaultLocation.Latitude, &client.DefaultLocation.Longitude, &info)
		if err != nil {
			return nil, err
		}
		if infoParser != nil {
			err = infoParser(&client, info)
			if err != nil {
				return nil, err
			}
		}
		clients = append(clients, client)
	}

	return clients, nil
}

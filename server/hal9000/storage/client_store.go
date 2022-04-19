package storage

import (
	"errors"

	"github.com/johnjones4/hal-9000/server/hal9000/core"
)

var (
	errorClientNotFound = errors.New("client not found")
)

type ClientStoreRecord struct {
	Client core.Client `json:"client"`
	Key    string      `json:"key"`
}

type ClientStore struct {
	store
	clients []ClientStoreRecord
}

func NewClientStore(configuration Configuration) *ClientStore {
	return &ClientStore{
		store: store{
			path: configuration.ClientsPath,
		},
		clients: make([]ClientStoreRecord, 0),
	}
}

func (s *ClientStore) GetClient(id string) (ClientStoreRecord, error) {
	for _, c := range s.clients {
		if c.Client.ID == id {
			return c, nil
		}
	}
	return ClientStoreRecord{}, errorClientNotFound
}

func (s *ClientStore) Load() error {
	return s.load(&s.clients)
}

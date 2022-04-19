package storage

import (
	"errors"

	"github.com/johnjones4/hal-9000/server/hal9000/core"
)

type UserStoreRecord struct {
	User core.User `json:"user"`
}

type UserStore struct {
	store
	users []UserStoreRecord
}

var (
	errorInvalidPassword = errors.New("invalid password")
	errorUserNotFound    = errors.New("user not found")
)

func NewUserStore(configuration Configuration) *UserStore {
	return &UserStore{
		store: store{
			path: configuration.UsersPath,
		},
		users: make([]UserStoreRecord, 0),
	}
}

func (us *UserStore) GetUser(username string) (UserStoreRecord, error) {
	for _, ur := range us.users {
		if ur.User.Name == username {
			return ur, nil
		}
	}
	return UserStoreRecord{}, errorUserNotFound
}

func (us *UserStore) Load() error {
	return us.load(&us.users)
}

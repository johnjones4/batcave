package storage

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"main/core"
	"sync"
)

type UserStoreRecord struct {
	User     core.User `json:"user"`
	Password string    `json:"password"`
}

type UserStore struct {
	store
	users []UserStoreRecord
}

func NewUserStore(path string) *UserStore {
	return &UserStore{
		store: store{
			path: path,
			lock: sync.Mutex{},
		},
		users: make([]UserStoreRecord, 0),
	}
}

func (us *UserStore) Login(username, password string) (core.User, error) {
	userRecord, err := us.GetUser(username)
	if err != nil {
		return core.User{}, err
	}
	hashed := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
	if hashed != userRecord.Password {
		return core.User{}, errors.New("invalid password")
	}
	return userRecord.User, nil
}

func (us *UserStore) GetUser(username string) (UserStoreRecord, error) {
	us.lock.Lock()
	defer us.lock.Unlock()
	for _, ur := range us.users {
		if ur.User.Name == username {
			return ur, nil
		}
	}
	return UserStoreRecord{}, errors.New("user not found")
}

func (us *UserStore) Load() error {
	return us.load(&us.users)
}
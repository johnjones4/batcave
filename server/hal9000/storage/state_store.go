package storage

import (
	"sync"

	"github.com/johnjones4/hal-9000/hal9000/core"
)

type StateStore struct {
	store
	states map[string]string
}

func NewStateStore(path string) *StateStore {
	return &StateStore{
		store: store{
			path: path,
			lock: sync.Mutex{},
		},
		states: make(map[string]string),
	}
}

func (ss *StateStore) Load() error {
	return ss.load(&ss.states)
}

func (ss *StateStore) Save() error {
	return ss.save(ss.states)
}

func (ss *StateStore) GetStateForUser(user core.User) (core.State, error) {
	ss.lock.Lock()
	defer ss.lock.Unlock()

	if state, ok := ss.states[user.Name]; ok {
		return core.State{
			State: state,
			User:  user,
		}, nil
	}

	return core.State{
		State: core.StateDefault,
		User:  user,
	}, nil
}

func (ss *StateStore) SetStateForUSer(state core.State) error {
	ss.lock.Lock()
	ss.states[state.User.Name] = state.State
	ss.lock.Unlock()
	return ss.Save()
}

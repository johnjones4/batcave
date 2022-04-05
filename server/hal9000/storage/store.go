package storage

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type store struct {
	path string
	lock sync.Mutex
}

func (s *store) load(dest interface{}) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	contents, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	err = json.Unmarshal(contents, dest)
	if err != nil {
		return err
	}

	return nil
}

func (s *store) save(source interface{}) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	contents, err := json.Marshal(source)
	if err != nil {
		return err
	}

	err = os.WriteFile(s.path, contents, 0777)
	if err != nil {
		return err
	}

	return nil
}

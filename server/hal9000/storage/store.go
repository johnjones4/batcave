package storage

import (
	"encoding/json"
	"errors"
	"os"
)

type store struct {
	path string
}

func (s *store) load(dest interface{}) error {
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

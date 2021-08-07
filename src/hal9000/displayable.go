package hal9000

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

const (
	DisplaySourceGoogle = "google"
	DisplayTypeVideo    = "video"
)

type Displayable struct {
	Names  []string `json:"names"`
	ID     string   `json:"id"`
	Type   string   `json:"type"`
	Source string   `json:"source"`
}

var displayables []Displayable

var ErrorDisplayNotFound = errors.New("display not found")

func InitDisplayables() error {
	bytes, err := ioutil.ReadFile(os.Getenv("DISPLAYABLES_MANIFEST_PATH"))
	if err != nil {
		return err
	}
	displayables = nil
	err = json.Unmarshal(bytes, &displayables)
	if err != nil {
		return err
	}
	return nil
}

func FindDisplayableInString(str string) (Displayable, error) {
	lcStr := strings.ToLower(str)
	for _, displayable := range displayables {
		for _, name := range displayable.Names {
			lcName := strings.ToLower(name)
			if strings.Contains(lcStr, lcName) {
				return displayable, nil
			}
		}
	}
	return Displayable{}, ErrorDisplayNotFound
}

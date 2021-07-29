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

type Display struct {
	Names  []string `json:"names"`
	ID     string   `json:"id"`
	Type   string   `json:"type"`
	Source string   `json:"source"`
}

var displays []Display

var ErrorDisplayNotFound = errors.New("display not found")

func InitDisplays() error {
	bytes, err := ioutil.ReadFile(os.Getenv("DISPLAY_MANIFEST_PATH"))
	if err != nil {
		return err
	}
	displays = nil
	err = json.Unmarshal(bytes, &displays)
	if err != nil {
		return err
	}
	return nil
}

func FindDisplayInString(str string) (Display, error) {
	lcStr := strings.ToLower(str)
	for _, display := range displays {
		for _, name := range display.Names {
			lcName := strings.ToLower(name)
			if strings.Contains(lcStr, lcName) {
				return display, nil
			}
		}
	}
	return Display{}, ErrorDisplayNotFound
}

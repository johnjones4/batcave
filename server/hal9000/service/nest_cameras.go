package service

import (
	"encoding/json"
	"os"

	"github.com/johnjones4/hal-9000/server/hal9000/core"
)

type NestCamera struct {
	NamesStrs []string `json:"names"`
	URLstr    string   `json:"url"`
}

func (c NestCamera) Names() []string {
	return c.NamesStrs
}

func (c NestCamera) URL(core.Inbound) (string, error) {
	return c.URLstr, nil
}

type NestCamerasConfiguration struct {
	CamerasPath string
}

type NestCameras struct {
	cameras []NestCamera
}

func NewNestCameras(configuration NestCamerasConfiguration) (*NestCameras, error) {
	nc := NestCameras{}

	contents, err := os.ReadFile(configuration.CamerasPath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(contents, &nc.cameras)
	if err != nil {
		return nil, err
	}

	return &nc, nil
}

func (nc *NestCameras) Displays() []Displayable {
	d := make([]Displayable, len(nc.cameras))
	for i, c := range nc.cameras {
		d[i] = c
	}
	return d
}

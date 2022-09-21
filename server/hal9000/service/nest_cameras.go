package service

import (
	"context"
	"encoding/json"
	"os"
)

type NestCamera struct {
	NamesStrs []string `json:"names"`
	URLstr    string   `json:"url"`
}

func (c NestCamera) Names() []string {
	return c.NamesStrs
}

func (c NestCamera) URL(context.Context) (string, error) {
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

func (nc *NestCameras) Name() string {
	return "nest"
}

func (nc *NestCameras) Info(c context.Context) (interface{}, error) {
	displays := nc.Displays()
	info := make(map[string]string)
	for _, d := range displays {
		names := d.Names()
		url, err := d.URL(c)
		if err != nil {
			return nil, err
		}
		info[names[0]] = url
	}
	return info, nil
}

func (nc *NestCameras) Displays() []Displayable {
	d := make([]Displayable, len(nc.cameras))
	for i, c := range nc.cameras {
		d[i] = c
	}
	return d
}

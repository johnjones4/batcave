package service

import "context"

type InfoService interface {
	Name() string
	Info(context.Context) (interface{}, error)
}

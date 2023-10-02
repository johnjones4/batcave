package core

import (
	"github.com/sirupsen/logrus"
)

type Env struct {
	HttpHost      string `env:"HTTP_HOST"`
	ServiceConfig string `env:"SERVICE_CONFIG"`
	DatabaseURL   string `env:"DATABASE_URL"`
}

type HookableLogger interface {
	logrus.FieldLogger
	AddHook(hook logrus.Hook)
}

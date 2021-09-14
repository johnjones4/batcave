package types

import "time"

type KVStore interface {
	GetString(key string, defaultVal string) string
	GetBytes(key string, defaultVal []byte) []byte
	GetFloat(key string, defaultVal float64) float64
	GetInt(key string, defaultVal int) int
	Set(key string, value interface{}, expiration time.Time) error
	SetBytes(key string, value []byte, expiration time.Time) error
}

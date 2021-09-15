package types

import "time"

type KVStore interface {
	GetString(key string, defaultVal string) string
	GetInterface(key string, iface interface{}) error
	GetBytes(key string, defaultVal []byte) []byte
	GetFloat(key string, defaultVal float64) float64
	GetInt(key string, defaultVal int) int
	Set(key string, value interface{}, expiration time.Time) error
	SetBytes(key string, value []byte, expiration time.Time) error
	SetInterface(key string, iface interface{}, expiration time.Time) error
}

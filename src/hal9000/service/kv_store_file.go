package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hal9000/types"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

func InitFileKVStore() (types.KVStore, error) {
	bytes, err := ioutil.ReadFile(os.Getenv("KV_FILE_PATH"))
	if err != nil {
		return nil, err
	}

	var data map[string]string
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}

	return &FileKVStore{data}, nil
}

type FileKVStore struct {
	data map[string]string
}

func (store *FileKVStore) GetString(key string, defaultVal string) string {
	if val, ok := store.data[key]; ok {
		return val
	}
	return defaultVal
}

func (store *FileKVStore) GetBytes(key string, defaultVal []byte) []byte {
	if val, ok := store.data[key]; ok {
		bytes, err := base64.StdEncoding.DecodeString(val)
		if err == nil {
			return bytes
		}
	}
	return defaultVal
}

func (store *FileKVStore) GetFloat(key string, defaultVal float64) float64 {
	if val, ok := store.data[key]; ok {
		fval, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return defaultVal
		}
		return fval
	}
	return defaultVal
}

func (store *FileKVStore) GetInt(key string, defaultVal int) int {
	if val, ok := store.data[key]; ok {
		ival, err := strconv.Atoi(val)
		if err != nil {
			return defaultVal
		}
		return ival
	}
	return defaultVal
}

func (store *FileKVStore) Set(key string, value interface{}, expiration time.Time) error {
	store.data[key] = fmt.Sprint(value)

	bytes, err := json.Marshal(store.data)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(os.Getenv("KV_FILE_PATH"), bytes, 0777)
	if err != nil {
		return err
	}

	return nil
}

func (store *FileKVStore) SetBytes(key string, value []byte, expiration time.Time) error {
	return store.Set(key, base64.StdEncoding.EncodeToString([]byte(value)), expiration)
}

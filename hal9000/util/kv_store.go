package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

var kvStore map[string]string

func InitKVStore() error {
	bytes, err := ioutil.ReadFile(os.Getenv("KV_FILE_PATH"))
	if err != nil {
		return err
	}

	kvStore = nil
	err = json.Unmarshal(bytes, &kvStore)
	if err != nil {
		return err
	}

	return nil
}

func GetKVValueString(key string, defaultVal string) string {
	if val, ok := kvStore[key]; ok {
		return val
	}
	return defaultVal
}

func GetKVValueFloat(key string, defaultVal float64) float64 {
	if val, ok := kvStore[key]; ok {
		fval, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return defaultVal
		}
		return fval
	}
	return defaultVal
}

func GetKVValueInt(key string, defaultVal int) int {
	if val, ok := kvStore[key]; ok {
		ival, err := strconv.Atoi(val)
		if err != nil {
			return defaultVal
		}
		return ival
	}
	return defaultVal
}

func SetKVValue(key string, value interface{}) error {
	kvStore[key] = fmt.Sprint(value)

	bytes, err := json.Marshal(kvStore)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(os.Getenv("KV_FILE_PATH"), bytes, 0777)
	if err != nil {
		return err
	}

	return nil
}

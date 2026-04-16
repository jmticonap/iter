package utils

import (
	"os"
	"strconv"
)

func GetEnvOrDefault[T string | int | bool](key string, defaultValue T) T {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	var result any
	var err error

	switch any(defaultValue).(type) {
	case string:
		result = value
	case int:
		result, err = strconv.Atoi(value)
	case bool:
		result, err = strconv.ParseBool(value)
	}

	if err != nil {
		return defaultValue
	}

	return result.(T)
}

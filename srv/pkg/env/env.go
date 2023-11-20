package env

import (
	"os"
	"strconv"
)

func GetAsString(key string, value string) string {
	s, ok := os.LookupEnv(key)
	if !ok {
		return value
	}
	return s
}

func GetAsInt[T ~int | ~int32 | ~int8 | ~int64](key string, value T) T {
	s, ok := os.LookupEnv(key)
	if !ok {
		return value
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return value
	}
	return T(v)
}

func GetAsBool(key string, value bool) bool {
	s, ok := os.LookupEnv(key)
	if !ok {
		return value
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		return value
	}
	return v
}

package utils

import (
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

func GetEnv(key string, defaultValue ...string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	} else {
		return ""
	}
}

func GetEnvInt(key string, defaultValue ...int) int {
	if value, exists := os.LookupEnv(key); exists {
		parse, err := strconv.Atoi(value)
		if err == nil {
			return parse
		}
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	} else {
		return 0
	}
}

func GetEnvBool(key string, defaultValue ...bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		parse, err := strconv.ParseBool(value)
		if err == nil {
			return parse
		}
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	} else {
		return false
	}
}

func MustGetEnv(key string, desc ...string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	log.Panic().Str("env", key).Strs("desc", desc).Msg("missing env")

	return ""
}

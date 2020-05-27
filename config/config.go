package config

import (
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

func GetString(key, fallback string) string {
	var result string
	if os.Getenv(key) != "" {
		result = os.Getenv(key)
	} else {
		result = fallback
	}

	log.Debug().Msgf("Successfully loaded env var: %s=%v", key, result)

	return result
}

func GetInt(key string, fallback int) int {
	var result int
	if os.Getenv(key) != "" {
		intVal, err := strconv.ParseInt(os.Getenv(key), 10, 64)
		if err != nil {
			log.Panic().Err(err).Msgf("error parsing int from env var: %s", key)
		}
		result = int(intVal)
	} else {
		result = fallback
	}

	log.Debug().Msgf("Successfully loaded env var: %s=%v", key, result)

	return result
}

func GetBool(key string, fallback bool) bool {
	var result bool
	if os.Getenv(key) != "" {
		value, err := strconv.ParseBool(os.Getenv(key))
		if err != nil {
			log.Panic().Err(err).Msgf("error parsing bool from env var: %s", key)
		}
		result = value
	} else {
		result = fallback
	}

	log.Debug().Msgf("Successfully loaded env var: %s=%v", key, result)

	return result
}

func GetStringOrDie(key string) string {
	value := os.Getenv(key)

	if value == "" {
		log.Panic().Msgf("Could not find env var: %s", key)
	}

	log.Debug().Msgf("Successfully loaded env var: %s=%s", key, value)
	return value
}

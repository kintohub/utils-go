package utils

import (
	_ "github.com/joho/godotenv/autoload"
	"os"
	"strconv"
)

func Get(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

func GetInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			return int(intVal)
		}
	}
	return fallback
}

func GetFloat(key string, fallback float64) float64 {
	if value, ok := os.LookupEnv(key); ok {
		floatVal, err := strconv.ParseFloat(value, 64)
		if err == nil {
			return floatVal
		}
	}
	return fallback
}

func GetBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		boolVal, err := strconv.ParseBool(value)
		if err == nil {
			return boolVal
		}
	}
	return fallback
}

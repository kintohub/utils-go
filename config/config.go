package config

import (
	"github.com/kintohub/utils-go/klog"
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

	klog.Debugf("Successfully loaded env var: %s=%v", key, result)

	return result
}

func GetInt(key string, fallback int) int {
	var result int
	if os.Getenv(key) != "" {
		intVal, err := strconv.Atoi(os.Getenv(key))
		if err != nil {
			klog.PanicfWithError(err, "error parsing int from env var: %s", key)
		}
		result = int(intVal)
	} else {
		result = fallback
	}

	klog.Debugf("Successfully loaded env var: %s=%v", key, result)

	return result
}

func GetIntOrDie(key string) int {
	value, err := strconv.Atoi(GetStringOrDie(key))
	if err != nil {
		klog.PanicfWithError(err, "error parsing int from env var: %s", key)
	}

	klog.Debugf("Successfully loaded env var: %s=%d", key, value)

	return value
}

func GetBool(key string, fallback bool) bool {
	var result bool
	if os.Getenv(key) != "" {
		value, err := strconv.ParseBool(os.Getenv(key))
		if err != nil {
			klog.PanicfWithError(err, "error parsing bool from env var: %s", key)
		}
		result = value
	} else {
		result = fallback
	}

	klog.Debugf("Successfully loaded env var: %s=%v", key, result)

	return result
}

func GetBoolOrDie(key string) bool {
	value, err := strconv.ParseBool(GetStringOrDie(key))
	if err != nil {
		klog.PanicfWithError(err, "error parsing bool from env var: %s", key)
	}

	klog.Debugf("Successfully loaded env var: %s=%t", key, value)

	return value
}

func GetStringOrDie(key string) string {
	value := os.Getenv(key)

	if value == "" {
		klog.Panicf("Could not find env var: %s", key)
	}

	klog.Debugf("Successfully loaded env var: %s=%s", key, value)
	return value
}

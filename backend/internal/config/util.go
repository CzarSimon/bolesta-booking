package config

import (
	"os"
	"strconv"
)

func getBoolEnvVar(key string, defaultValue bool) bool {
	str := os.Getenv(key)
	if str == "" {
		return defaultValue
	}

	val, err := strconv.ParseBool(str)
	if err != nil {
		log.Fatalf("failed to parse %s as boolean. Key=%s: %w", str, key, err)
	}

	return val
}

func getIntEnvVar(key string, defaultValue int) int {
	str := os.Getenv(key)
	if str == "" {
		return defaultValue
	}

	val, err := strconv.Atoi(str)
	if err != nil {
		log.Fatalf("failed to parse %s as int. Key=%s: %w", str, key, err)
	}

	return val
}

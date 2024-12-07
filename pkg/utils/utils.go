package utils

import "os"

func GetEnvWithDefault(envKey, defaultValue string) string {
	value, ok := os.LookupEnv(envKey)
	if ok {
		return value
	}
	return defaultValue
}

package utils

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var userHome string

func GetEnvWithDefault(envKey, defaultValue string) string {
	value, ok := os.LookupEnv(envKey)
	if ok {
		return value
	}
	return defaultValue
}

// GetUserHome gets the user's home directory and caches it for future calls.
func GetUserHome() string {
	if userHome == "" {
		var err error
		userHome, err = os.UserHomeDir()
		cobra.CheckErr(err)
	}
	return userHome
}

// ExpandUserHome expands any leading "~" in given path to the user's home directory.
func ExpandUserHome(path string) string {
	if path == "~" {
		return GetUserHome()
	}
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(GetUserHome(), path[2:])
	}
	return path
}

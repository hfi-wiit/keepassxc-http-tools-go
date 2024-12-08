//go:build linux
// +build linux

package utils

import "path"

// GetConfigDir returns the path to the users config dir.
func GetConfigDir() string {
	return GetEnvWithDefault("XDG_CONFIG_HOME", path.Join(GetUserHome(), ".config"))
}
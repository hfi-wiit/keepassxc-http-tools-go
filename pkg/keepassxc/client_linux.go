//go:build linux
// +build linux

package keepassxc

import (
	"errors"
	"fmt"
	"keepassxc-http-tools-go/pkg/utils"
	"net"
	"os"
	"path"
)

// SocketPath tries to find the path to the socket of the keepassxc http api - Linux version
func SocketPath() (string, error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	lookupPaths := []string{
		path.Join(userHome, utils.LinuxSnapCommonSubdir),
		utils.GetEnvWithDefault(utils.LinuxEnvXdgRuntimeDir, fmt.Sprintf("/run/user/%d/", os.Getuid())),
		utils.GetEnvWithDefault(utils.LinuxEnvTmpDir, utils.LinuxEnvTmpDirDefault),
	}

	var filename string
	for _, base := range lookupPaths {
		filename = path.Join(base, utils.SocketFileName)
		if _, err := os.Stat(filename); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return "", fmt.Errorf("%w, error on file %s: %w", utils.ErrKeepassxcSocketNotFound, filename, err)
		}
		break
	}

	if filename == "" {
		return "", utils.ErrKeepassxcSocketNotFound
	}

	return filename, nil
}

// connect implements the os specific socket connection action - Linux version
func connect(socketPath string) (net.Conn, error) {
	return net.DialUnix("unix", nil, &net.UnixAddr{Name: socketPath, Net: "unix"})
}

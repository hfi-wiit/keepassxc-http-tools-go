//go:build darwin
// +build darwin

package keepassxc

import (
	"fmt"
	"keepassxc-http-tools-go/pkg/utils"
	"net"
	"os"
	"path/filepath"
)

// SocketPath finds the path to the socket file descriptor of the keepassxc http api.
func SocketPath() (string, error) {
	tmpDir, ok := os.LookupEnv(utils.DarwinEnvTmpDir)
	if !ok {
		return "", fmt.Errorf("%w, $%s not set", utils.ErrKeepassxcSocketNotFound, utils.DarwinEnvTmpDir)
	}

	path := filepath.Join(tmpDir, utils.SocketFileName)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", utils.ErrKeepassxcSocketNotFound
	}
	return path, nil
}

// connect implements the os specific socket connection action.
func connect(socketPath string) (net.Conn, error) {
	return net.DialUnix("unix", nil, &net.UnixAddr{Name: socketPath, Net: "unix"})
}

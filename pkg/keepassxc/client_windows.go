//go:build windows
// +build windows

package keepassxc

import (
	"fmt"
	"keepassxc-http-tools-go/pkg/utils"
	"net"
	"os"

	"github.com/Microsoft/go-winio"
)

// SocketPath tries to find the path to the socket of the keepassxc http api - Windows version
func SocketPath() (string, error) {
	return fmt.Sprintf(`\\.\pipe\%s_%s`, utils.SocketFileName, os.Getenv(utils.WindowsEnvVarUsername)), nil
}

// connect implements the os specific socket connection action - Windows version
func connect(socketPath string) (net.Conn, error) {
	return winio.DialPipe(socketPath, nil)
}

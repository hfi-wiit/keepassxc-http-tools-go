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

func SocketPath() (string, error) {
	return fmt.Sprintf(`\\.\pipe\%s_%s`, utils.SocketFileName, os.Getenv(utils.WindowsEnvVarUsername)), nil
}

func connect(socketPath string) (net.Conn, error) {
	return winio.DialPipe(socketPath, nil)
}

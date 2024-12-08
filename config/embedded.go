package config

import (
	"embed"

	"github.com/spf13/cobra"
)

//go:embed *.yaml
var embeddedFs embed.FS

func LoadFile(path string) string {
	data, err := embeddedFs.ReadFile(path)
	cobra.CheckErr(err)
	return string(data)
}

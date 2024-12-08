/*
Copyright © 2024 Heiko Finzel heiko.finzel@wiit.cloud
*/
package cmd

import (
	"fmt"
	"keepassxc-http-tools-go/config"
	"keepassxc-http-tools-go/pkg/utils"
	"path"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Args:  cobra.NoArgs,
	Run:   configCmdRun,
	Short: "Print an example config to stdout",
	Long: fmt.Sprintf(`Print an example config to stdout.
	
The default config file is %s.
Most of the config is optional, see the example config from this command for details.
The "assoc" profile will be created and saved automatically on first connection to the database.`,
		path.Join(utils.GetConfigDir(), utils.ConfigFileNameDefault),
	),
	Example: fmt.Sprintf("  %s config", utils.ApplicationNameShort),
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func configCmdRun(cmd *cobra.Command, args []string) {
	fmt.Printf("%s\n", config.LoadFile("kpht.yaml"))
}

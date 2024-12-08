/*
Copyright © 2024 Heiko Finzel heiko.finzel@wiit.cloud
*/
package cmd

import (
	"fmt"
	"keepassxc-http-tools-go/pkg/utils"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// this version string is set at compile time in the Makefile
var Version = "dev"

// global flags storage
type GlobalFlags struct {
	// path to the config file
	ConfigFile string
}

// global flags storage
var globalFlags = GlobalFlags{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     utils.ApplicationNameShort,
	Version: Version,
	Args:    cobra.NoArgs,
	Short:   "A command line client to interact with keepassxc's http api.",
	Long: fmt.Sprintf(`A command line client to interact with keepassxc's http api.
	
To learn more about the config use "%s config" or "%s config -h".`,
		utils.ApplicationNameShort,
		utils.ApplicationNameShort,
	),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&globalFlags.ConfigFile, "config", "c",
		path.Join(utils.GetConfigDir(), utils.ConfigFileNameDefault), "the config file")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetDefault(utils.ConfigKeypathEntryIdentifier, []string{"%s (%s)", "name", "login"})
	viper.SetDefault(utils.ConfigKeypathClipDefaultCopy, []string{utils.ConfigDefaultClipDefaultCopy})
	viper.SetDefault(utils.ConfigKeypathScriptIndicatorUrl, utils.ConfigDefaultScriptIndicatorUrl)
	viper.SetConfigFile(utils.ExpandUserHome(globalFlags.ConfigFile))
	// read in environment variables that match, but only with KGHT_ prefix
	viper.SetEnvPrefix(utils.ConfigEnvPrefix)
	viper.AutomaticEnv()
	// If a config file is found, read it in.
	viper.ReadInConfig()
}

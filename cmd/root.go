/*
Copyright © 2024 Heiko Finzel heiko.finzel@wiit.cloud
*/
package cmd

import (
	"keepassxc-http-tools-go/pkg/utils"
	"os"
	"path"

	"github.com/kevinburke/nacl"
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
	Long:    "A command line client to interact with keepassxc's http api.",
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
	viper.SetConfigFile(utils.ExpandUserHome(globalFlags.ConfigFile))
	// read in environment variables that match, but only with KGHT_ prefix
	viper.SetEnvPrefix(utils.ConfigEnvPrefix)
	viper.AutomaticEnv()
	// If a config file is found, read it in.
	viper.ReadInConfig()
}

// Implements the KeepassxcClientProfile interface for viper config
type ViperKeepassxcProfile struct{}

func (p ViperKeepassxcProfile) GetAssocName() string {
	return viper.GetString(utils.ConfigKeypathAssocName)
}

func (p ViperKeepassxcProfile) GetAssocKey() nacl.Key {
	b64String := viper.GetString(utils.ConfigKeypathAssocKey)
	if b64String == "" {
		return nil
	}
	return utils.B64ToNaclKey(b64String)
}

func (p ViperKeepassxcProfile) SetAssoc(name string, key nacl.Key) error {
	viper.Set(utils.ConfigKeypathAssocName, name)
	viper.Set(utils.ConfigKeypathAssocKey, utils.NaclKeyToB64(key))
	return viper.WriteConfig()
}

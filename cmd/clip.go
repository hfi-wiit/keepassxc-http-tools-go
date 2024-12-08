/*
Copyright © 2024 Heiko Finzel heiko.finzel@wiit.cloud
*/
package cmd

import (
	"fmt"
	"keepassxc-http-tools-go/pkg/keepassxc"
	"keepassxc-http-tools-go/pkg/utils"
	"strings"
	"time"

	fzf "github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	clip "golang.design/x/clipboard"
)

// clipCmd represents the clip command
var clipCmd = &cobra.Command{
	Use:   "clip [namefilter]",
	Args:  cobra.ArbitraryArgs,
	Run:   clipCmdRun,
	Short: "Copy data from an entry to clipboard.",
	Long: `Copy data from an entry to clipboard.

TODO more details...`,
}

func init() {
	rootCmd.AddCommand(clipCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clipCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clipCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func clipCmdRun(cmd *cobra.Command, args []string) {
	client, err := keepassxc.NewClient(utils.ViperKeepassxcProfile{})
	cobra.CheckErr(err)
	defer client.Disconnect()
	entries, err := client.GetLogins(utils.ScriptIndicatorUrl)
	cobra.CheckErr(err)

	filter := utils.ScriptIndicatorUrl
	if len(args) > 0 {
		filter = strings.Join(args, " ")
		entries = entries.FilterByName(args...)
	}
	var selectedEntry *keepassxc.Entry
	switch len(entries) {
	case 0:
		cobra.CheckErr(fmt.Errorf("No logins match the search criteria: %s", filter))
	case 1:
		selectedEntry = entries[0]
	default:
		idx, err := fzf.Find(entries, func(i int) string {
			return entries[i].GetCombined(viper.GetStringSlice(utils.ConfigKeypathEntryIdentifier))
		})
		cobra.CheckErr(err)
		selectedEntry = entries[idx]
	}

	overrideMap := viper.GetStringMapStringSlice(utils.ConfigKeypathClipCopy)
	copyKeys, ok := overrideMap[selectedEntry.Uuid]
	if !ok {
		copyKeys = viper.GetStringSlice(utils.ConfigKeypathClipDefaultCopy)
	}
	copyValue := selectedEntry.GetCombined(copyKeys)

	err = clip.Init()
	cobra.CheckErr(err)
	clip.Write(clip.FmtText, []byte(copyValue))
	// it seems we need at least some (~5?) milliseconds to be sure the value is copied into clipboard
	time.Sleep(100 * time.Millisecond)

	fmt.Printf("Copied %s from %s\n",
		utils.GetCombinedKeys(copyKeys),
		selectedEntry.GetCombined(viper.GetStringSlice(utils.ConfigKeypathEntryIdentifier)))
}

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

// clip flags storage
type ClipFlags struct {
	CopyLogin    bool
	CopyPassword bool
	CopyTotp     bool
	CopyUuid     bool
}

// clip flags storage
var clipFlags = ClipFlags{}

// clipCmd represents the clip command
var clipCmd = &cobra.Command{
	Use:   "clip [namefilter...]",
	Args:  cobra.ArbitraryArgs,
	Run:   clipCmdRun,
	Short: "Copy data from an entry to clipboard",
	Long: `Copy data from an entry to clipboard.

TODO more details...`,
}

func init() {
	rootCmd.AddCommand(clipCmd)
	clipCmd.Flags().BoolVarP(&clipFlags.CopyLogin, "login", "l", false,
		"Copy login instead of the field specified in config.")
	clipCmd.Flags().BoolVarP(&clipFlags.CopyPassword, "password", "p", false,
		"Copy password instead of the field specified in config.")
	clipCmd.Flags().BoolVarP(&clipFlags.CopyTotp, "totp", "t", false,
		"Copy totp instead of the field specified in config.")
	clipCmd.Flags().BoolVarP(&clipFlags.CopyUuid, "uuid", "u", false,
		"Copy uuid instead of the field specified in config.")
	clipCmd.MarkFlagsMutuallyExclusive("login", "password", "totp", "uuid")
}

func clipCmdRun(cmd *cobra.Command, args []string) {
	// get entries from keepassxc
	client, err := keepassxc.NewClient(utils.ViperKeepassxcProfile{})
	cobra.CheckErr(err)
	defer client.Disconnect()
	scriptIndicatorUrl := viper.GetString(utils.ConfigKeypathScriptIndicatorUrl)
	entries, err := client.GetLogins(scriptIndicatorUrl)
	cobra.CheckErr(err)

	// filter entries by configured groups
	groups := viper.GetStringSlice(utils.ConfigKeypathClipFilterGroups)
	if len(groups) > 0 {
		entries = entries.FilterByGroup(groups...)
	}

	// filter entries by optional name filter arguments
	filter := scriptIndicatorUrl
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
		// and if multiple are left, chose one per fuzzy finder
		idx, err := fzf.Find(entries, func(i int) string {
			return entries[i].GetCombined(viper.GetStringSlice(utils.ConfigKeypathEntryIdentifier))
		})
		cobra.CheckErr(err)
		selectedEntry = entries[idx]
	}

	// select the value(s) to copy from the selected entry, either from flag
	var copyKeys []string
	if clipFlags.CopyTotp {
		copyKeys = []string{"totp"}
	} else if clipFlags.CopyPassword {
		copyKeys = []string{"password"}
	} else if clipFlags.CopyLogin {
		copyKeys = []string{"login"}
	} else if clipFlags.CopyUuid {
		copyKeys = []string{"uuid"}
	} else {
		// or from config
		var ok bool
		overrideMap := viper.GetStringMapStringSlice(utils.ConfigKeypathClipCopy)
		copyKeys, ok = overrideMap[selectedEntry.Uuid]
		if !ok {
			copyKeys = viper.GetStringSlice(utils.ConfigKeypathClipDefaultCopy)
		}
	}
	copyValue := selectedEntry.GetCombined(copyKeys)

	// copy that value to clipboard
	err = clip.Init()
	cobra.CheckErr(err)
	clip.Write(clip.FmtText, []byte(copyValue))
	// it seems we need at least some (~5?) milliseconds to be sure the value is copied into clipboard
	time.Sleep(100 * time.Millisecond)

	fmt.Printf("Copied %s from %s\n",
		utils.GetCombinedKeys(copyKeys),
		selectedEntry.GetCombined(viper.GetStringSlice(utils.ConfigKeypathEntryIdentifier)))
}

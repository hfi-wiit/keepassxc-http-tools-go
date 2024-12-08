/*
Copyright © 2024 Heiko Finzel heiko.finzel@wiit.cloud
*/
package cmd

import (
	"fmt"
	"keepassxc-http-tools-go/pkg/keepassxc"
	"keepassxc-http-tools-go/pkg/utils"
	"time"

	fzf "github.com/ktr0731/go-fuzzyfinder"
	"github.com/spf13/cobra"
	clip "golang.design/x/clipboard"
)

// clipCmd represents the clip command
var clipCmd = &cobra.Command{
	Use:   "clip",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: clipCmdRun,
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
	client, err := keepassxc.NewClient(ViperKeepassxcProfile{})
	cobra.CheckErr(err)
	defer client.Disconnect()
	entries, err := client.GetLogins(utils.ScriptIndicatorUrl)
	cobra.CheckErr(err)

	filter := utils.ScriptIndicatorUrl
	if len(args) > 0 {
		filter = args[0]
		entries = entries.FilterByName(filter)
	}
	var selectedEntry *keepassxc.Entry
	switch len(entries) {
	case 0:
		cobra.CheckErr(fmt.Errorf("No logins match the search criteria: %s", filter))
	case 1:
		selectedEntry = entries[0]
	default:
		idx, err := fzf.Find(entries, func(i int) string { return entries[i].Name })
		cobra.CheckErr(err)
		selectedEntry = entries[idx]
	}

	fmt.Printf("%+v\n", selectedEntry)
	err = clip.Init()
	cobra.CheckErr(err)
	clip.Write(clip.FmtText, []byte(selectedEntry.Password.Plaintext()))
	// it seems we need at least some (~5?) milliseconds to be sure the value is copied into clipboard
	time.Sleep(100 * time.Millisecond)
}

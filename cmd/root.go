package cmd

import (
	"gogcli/sdk"

	"github.com/spf13/cobra"
)

var debugMode bool
var cookieFile string
var sdkInst sdk.Sdk

var rootCmd = &cobra.Command{
	Use:   "gogcli",
	Short: "A Client to Interact with the GOG.com API",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		sdkInst = sdk.NewSdk(cookieFile)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cookieFile, "cookiefile", "c", "cookie", "Path were to read the user provided cookie file")
	rootCmd.PersistentFlags().BoolVarP(&debugMode, "debug", "d", false, "Provide additional more detailed ouputs to help troubleshoot the tool")

	rootCmd.AddCommand(generateGameDetailsCmd())
	rootCmd.AddCommand(generateOwnedGamesCmd())
	rootCmd.AddCommand(generateUserInfoCmd())
}

func Execute() error {
	return rootCmd.Execute()
}

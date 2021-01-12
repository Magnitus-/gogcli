package cmd

import (
	"fmt"
	"gogcli/sdk"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var debugMode bool
var cookieFile string
var sdkPtr *sdk.Sdk

var rootCmd = &cobra.Command{
	Use:   "gogcli",
	Short: "A Client to Interact with the GOG.com API",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		logger := log.New(os.Stdout, "gogcli sdk: ", log.Lshortfile)
		sdkPtr, err = sdk.NewSdk(cookieFile, logger)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cookieFile, "cookiefile", "c", "cookie", "Path were to read the user provided cookie file")
	rootCmd.MarkPersistentFlagFilename("cookiefile")
	rootCmd.PersistentFlags().BoolVarP(&debugMode, "debug", "d", false, "Provide additional more detailed ouputs to help troubleshoot the tool")

	rootCmd.AddCommand(generateGameDetailsCmd())
	rootCmd.AddCommand(generateOwnedGamesCmd())
	rootCmd.AddCommand(generateUserInfoCmd())
	rootCmd.AddCommand(generateManifestGenerationCmd())
}

func Execute() error {
	return rootCmd.Execute()
}

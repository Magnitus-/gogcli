package cmd

import (
	"fmt"
	"gogcli/logging"
	"gogcli/sdk"
	"os"

	"github.com/spf13/cobra"
)

var logLevel string
var cookieFile string
var cookieFileType string
var sdkPtr *sdk.Sdk
var logSource *logging.Source

var rootCmd = &cobra.Command{
	Use:   "gogcli",
	Short: "A Client to Interact with the GOG.com API",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		logSource = logging.CreateSource(logLevel)
		cookie, err := sdk.ReadCookie(cookieFile, cookieFileType)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		sdkPtr = sdk.NewSdk(cookie, logSource)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cookieFile, "cookiefile", "c", "cookie", "Path were to read the user provided cookie file")
	rootCmd.MarkPersistentFlagFilename("cookiefile")
	rootCmd.PersistentFlags().StringVarP(&cookieFileType, "cookiefile-type", "y", "default", "The type of cookie file. Can either be 'default' or 'netscape'")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "g", "info", "Logs below this level of significance won't be displayed. Possible values are: debug, info and warning")

	rootCmd.AddCommand(generateUpdateCmd())
	rootCmd.AddCommand(generateGogApiCmd())
	rootCmd.AddCommand(generateManifestCmd())
	rootCmd.AddCommand(generateMetadataCmd())
	rootCmd.AddCommand(generateStorageCmd())
	rootCmd.AddCommand(generateVersionCmd())
	rootCmd.AddCommand(generateActionsCmd())
}

func Execute() error {
	return rootCmd.Execute()
}

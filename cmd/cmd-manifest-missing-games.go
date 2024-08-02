package cmd

import (
	"fmt"
	"os"

	"gogcli/manifest"
	"gogcli/sdk"

	"github.com/spf13/cobra"
)

func generateManifestMissingGamesCmd() *cobra.Command {
	var m manifest.Manifest
	var manifestPath string
	var missingGamesFile string
	var terminalOutput bool
	var concurrency int
	var pause int

	manifestMissingGamesCmd := &cobra.Command{
		Use:   "missing-games",
		Short: "Command to retrieve a list of missing games from gog.com that your manifest doesn't have",
		PreRun: func(cmd *cobra.Command, args []string) {
			var err error
			m, err = loadManifestFromFile(manifestPath)
			processError(err)
		},
		Run: func(cmd *cobra.Command, args []string) {
			o, err := sdkPtr.GetAllOwnedGamesPagesSync("", concurrency, pause)
			if err != nil {
				fmt.Println("Could not retrieve owned games from gog.com: ", err)
				os.Exit(1)
			}

			missingGames := sdk.GetMissingGames(o, &m)

			processSerializableOutput(missingGames, []error{}, terminalOutput, missingGamesFile)
		},
	}

	manifestMissingGamesCmd.Flags().StringVarP(&manifestPath, "manifest", "m", "manifest.json", "Manifest file to check for missing games")
	manifestMissingGamesCmd.Flags().StringVarP(&missingGamesFile, "missing-games-file", "f", "missing-games.json", "File to output the missing games in if in json format")
	manifestMissingGamesCmd.Flags().BoolVarP(&terminalOutput, "terminal", "t", true, "If set to true and json format is used, the missing games will be output on the terminal instead of in a file")
	manifestMissingGamesCmd.Flags().IntVarP(&concurrency, "concurrency", "r", 4, "Maximum number of concurrent requests that will be made on the GOG api")
	manifestMissingGamesCmd.Flags().IntVarP(&pause, "pause", "s", 200, "Number of milliseconds to wait between batches of api calls")

	return manifestMissingGamesCmd
}
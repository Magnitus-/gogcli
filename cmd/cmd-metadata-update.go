package cmd

import (
	"encoding/json"
	"fmt"
	"gogcli/metadata"
	"gogcli/gameupdates"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

func generateMetadataUpdateCmd() *cobra.Command {
	var u *gameupdates.Updates
	var m metadata.Metadata
	var gameIds []int64
	var updateFile string
	var metadataFile string
	var progressFile string
	var warningFile string
	var concurrency int
	var pause int
	var tolerateDangles bool

	metadataUpdateCmd := &cobra.Command{
		Use:   "update",
		Short: "Update a metadata file based on changes from gog api",
		PreRun: func(cmd *cobra.Command, args []string) {
			var err error
			m, err = loadMetadataFromFile(metadataFile)
			processError(err)

			if len(gameIds) == 0 && updateFile == "" {
				fmt.Println("You either need to pass ids of games to update or pass an updates file", err)
				os.Exit(1)
			}

			if updateFile != "" {
				var bs []byte
				bs, err = ioutil.ReadFile(updateFile)
				if err != nil {
					fmt.Println("Could not load the updates: ", err)
					os.Exit(1)
				}

				u = &gameupdates.Updates{}
				err = json.Unmarshal(bs, u)
				if err != nil {
					fmt.Println("Updates file doesn't appear to contain valid json: ", err)
					os.Exit(1)
				}
			}

			CleanupFile(warningFile)
		},
		Run: func(cmd *cobra.Command, args []string) {
			ids := gameIds
			if u != nil {
				ids = append(ids, (*u).NewGames...)
				ids = append(ids, (*u).UpdatedGames...)
			}

			progressFn := PersistMetadataProgress(progressFile)
			writer := metadata.NewMetadataGamesWriter(
				metadata.NewMetadataGamesWriterState(m.Filter, ids, m.SkipImages),
				logSource,
			)
			errs := writer.Write( 
				sdkPtr.GenerateMetadataGameGetter(concurrency, pause, tolerateDangles),
				progressFn,
			)
			uMetadata, warnings := writer.State.Metadata, writer.State.Warnings
			
			if len(warnings) > 0 {
				warningsOutput := Errors{make([]string, len(warnings))}
				for idx, _ := range warnings {
					warningsOutput.Errors[idx] = warnings[idx]
				}
				processSerializableOutput(warningsOutput, []error{}, false, warningFile)
			}
			processErrors(errs)

			m.OverwriteGames(uMetadata.Games)

			processSerializableOutput(m, []error{}, false, metadataFile)

			CleanupFile(progressFile)
		},
	}

	metadataUpdateCmd.Flags().Int64SliceVarP(&gameIds, "id", "i", []int64{}, "Optional ids of games to update")
	metadataUpdateCmd.Flags().IntVarP(&concurrency, "concurrency", "r", 4, "Maximum number of concurrent requests that will be made on the GOG api")
	metadataUpdateCmd.Flags().IntVarP(&pause, "pause", "s", 200, "Number of milliseconds to wait between batches of api calls")
	metadataUpdateCmd.Flags().StringVarP(&metadataFile, "metadata-file", "f", "metadata.json", "Metadata file to update")
	metadataUpdateCmd.MarkFlagFilename("metadata-file")
	metadataUpdateCmd.Flags().StringVarP(&progressFile, "progress-file", "z", "metadata-update-progress.json", "File to save transient progress for the manifest update in")
	metadataUpdateCmd.Flags().StringVarP(&updateFile, "update", "u", "", "Optional update file containing new and updated game ids")
	metadataUpdateCmd.Flags().BoolVarP(&tolerateDangles, "tolerate-dangles", "d", true, "If set to true, undownloadable dangling files (ie, 404 code on download url) will be tolerated and will not prevent manifest generation")
	metadataUpdateCmd.Flags().StringVarP(&warningFile, "warning-file", "w", "metadata-update-warnings.json", "Warnings from files whose download url return 404 will be listed in this file. Will only be generated if tolerate-dangles is set to true")
	return metadataUpdateCmd
}

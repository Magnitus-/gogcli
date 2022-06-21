package cmd

import (
	"encoding/json"
	"fmt"
	"gogcli/manifest"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

func generateManifestUpdateCmd() *cobra.Command {
	var u *manifest.Updates
	var m manifest.Manifest
	var gameIds []int64
	var updateFile string
	var manifestFile string
	var warningFile string
	var duplicatesFile string
	var concurrency int
	var pause int
	var tolerateDangles bool
	var tolerateBadFileMetadata bool

	manifestUpdateCmd := &cobra.Command{
		Use:   "update",
		Short: "Update a manifest file based on changes from gog api",
		PreRun: func(cmd *cobra.Command, args []string) {
			var err error
			m, err = loadManifestFromFile(manifestFile)
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

				u = &manifest.Updates{}
				err = json.Unmarshal(bs, u)
				if err != nil {
					fmt.Println("Updates file doesn't appear to contain valid json: ", err)
					os.Exit(1)
				}
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			ids := gameIds
			if u != nil {
				ids = append(ids, (*u).NewGames...)
				ids = append(ids, (*u).UpdatedGames...)
			}

			uManifest, errs, warnings := sdkPtr.GetManifestFromIds(m.Filter, ids, concurrency, pause, tolerateDangles, tolerateBadFileMetadata)
			m.OverwriteGames(uManifest.Games)
			duplicates := m.Finalize()

			processErrors(errs)
			processSerializableOutput(m, []error{}, false, manifestFile)

			if len(duplicates) > 0 {
				processSerializableOutput(duplicates, []error{}, false, duplicatesFile)
			}
			if len(warnings) > 0 {
				warningsOutput := Errors{make([]string, len(warnings))}
				for idx, _ := range warnings {
					warningsOutput.Errors[idx] = warnings[idx].Error()
				}
				processSerializableOutput(warningsOutput, []error{}, false, warningFile)
			}
		},
	}

	manifestUpdateCmd.Flags().Int64SliceVarP(&gameIds, "id", "i", []int64{}, "Optional ids of games to update")
	manifestUpdateCmd.Flags().IntVarP(&concurrency, "concurrency", "r", 4, "Maximum number of concurrent requests that will be made on the GOG api")
	manifestUpdateCmd.Flags().IntVarP(&pause, "pause", "s", 200, "Number of milliseconds to wait between batches of api calls")
	manifestUpdateCmd.Flags().StringVarP(&manifestFile, "manifest-file", "f", "manifest.json", "Manifest file to update")
	manifestUpdateCmd.MarkFlagFilename("manifest")
	manifestUpdateCmd.Flags().StringVarP(&updateFile, "update", "u", "", "Optional update file containing new and updated game ids")
	manifestUpdateCmd.Flags().BoolVarP(&tolerateDangles, "tolerate-dangles", "d", true, "If set to true, undownloadable dangling files (ie, 404 code on download url) will be tolerated and will not prevent manifest generation")
	manifestUpdateCmd.Flags().StringVarP(&warningFile, "warning-file", "w", "manifest-update-warnings.json", "Warnings from files whose download url return 404 will be listed in this file. Will only be generated if tolerate-dangles is set to true")
	manifestUpdateCmd.Flags().StringVarP(&duplicatesFile, "duplicates-file", "l", "manifest-update-duplicates.json", "Files that had duplicate filenames within the same game and had to be renamed will be listed in this file")
	manifestUpdateCmd.Flags().BoolVarP(&tolerateBadFileMetadata, "tolerate-bad-metadata", "b", true, "Tolerate files for which metadata cannot be retrieved. The checksum will be infered by performing a throwaway file download instead.")
	return manifestUpdateCmd
}

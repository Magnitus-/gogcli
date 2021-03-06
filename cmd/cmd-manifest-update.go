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
	var concurrency int
	var pause int
	var tolerateDangles bool

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

			uManifest, errs, errs404 := sdkPtr.GetManifestFromIds(m.Filter, ids, concurrency, pause, tolerateDangles)
			m.OverwriteGames(uManifest.Games)
			m.Finalize()
			processSerializableOutput(m, errs, false, manifestFile)
			if len(errs404) > 0 {
				errs404Output := Errors{make([]string, len(errs404))}
				for idx, _ := range errs404 {
					errs404Output.Errors[idx] = errs404[idx].Error()
				}
				processSerializableOutput(errs404Output, []error{}, false, warningFile)
			}
		},
	}

	manifestUpdateCmd.Flags().Int64SliceVarP(&gameIds, "id", "i", []int64{}, "Optional ids of games to update")
	manifestUpdateCmd.Flags().IntVarP(&concurrency, "concurrency", "r", 10, "Maximum number of concurrent requests that will be made on the GOG api")
	manifestUpdateCmd.Flags().IntVarP(&pause, "pause", "s", 200, "Number of milliseconds to wait between batches of api calls")
	manifestUpdateCmd.Flags().StringVarP(&manifestFile, "manifest", "f", "manifest.json", "Manifest file to update")
	manifestUpdateCmd.MarkFlagFilename("manifest")
	manifestUpdateCmd.Flags().StringVarP(&updateFile, "update", "u", "", "Optional update file containing new and updated game ids")
	manifestUpdateCmd.Flags().BoolVarP(&tolerateDangles, "tolerate-dangles", "g", true, "If set to true, undownloadable dangling files (ie, 404 code on download url) will be tolerated and will not prevent manifest generation")
	manifestUpdateCmd.Flags().StringVarP(&warningFile, "warning-file", "w", "manifest-404-warnings.json", "Warnings from files whose download url return 404 will be listed in this file. Will only be generated if tolerate-dangles is set to true")
	return manifestUpdateCmd
}

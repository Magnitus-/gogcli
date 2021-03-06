package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gogcli/manifest"
	"gogcli/storage"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

func generateStoragePlanCmd() *cobra.Command {
	var m manifest.Manifest
	var manifestPath string
	var file string
	var terminalOutput bool
	var path string
	var storageType string

	show := func(a *manifest.GameActions) {
		var buf bytes.Buffer
		var output []byte
		
		output, _ = json.Marshal((*a))
		json.Indent(&buf, output, "", "  ")
		output = buf.Bytes()

		if terminalOutput {
			fmt.Println(string(output))
		} else {
			err := ioutil.WriteFile(file, output, 0644)
			if err != nil {
				fmt.Println("Could not write output to file:", err)
				os.Exit(1)
			}
		}
	}

	storagePlanCmd := &cobra.Command{
		Use:   "plan",
		Short: "Generate a plan of the actions that would be executed if a given manifest was applied to the storage",
		PreRun: func(cmd *cobra.Command, args []string) {
			var err error
			m, err = loadManifestFromFile(path)
			processError(err)
		},
		Run: func(cmd *cobra.Command, args []string) {
			var actions *manifest.GameActions

			gamesStorage, _ := getStorage(path, storageType, debugMode, "")

			err := storage.EnsureInitialization(gamesStorage)
			processError(err)

			actions, err = storage.PlanManifest(&m, storage.GetFileSystem(path, debugMode, ""))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			show(actions)
		},
	}

	storagePlanCmd.Flags().StringVarP(&manifestPath, "manifest", "m", "manifest.json", "Path were the manifest you want to apply is")
	storagePlanCmd.MarkFlagFilename("manifest")
	storagePlanCmd.Flags().StringVarP(&file, "file", "f", "actions.json", "File to output the plan in")
	storagePlanCmd.Flags().BoolVarP(&terminalOutput, "terminal", "t", false, "If set to true, the plan will be output on the terminal instead of in a file")
	storagePlanCmd.Flags().StringVarP(&path, "path", "p", "games", "Path to your games' storage (directory if it is of type fs, json configuration file if it is of type s3)")
	storagePlanCmd.Flags().StringVarP(&storageType, "storage", "k", "fs", "The type of storage you are using. Can be 'fs' (for file system) or 's3' (for s3 store)")

	return storagePlanCmd
}

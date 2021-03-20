package cmd

import (
	"gogcli/manifest"
	"gogcli/storage"

	"github.com/spf13/cobra"
)

func generateStorageRepairCmd() *cobra.Command {
	var m manifest.Manifest
	var manifestPath string
	var path string
	var storageType string
	var concurrency int

	storageRepairCmd := &cobra.Command{
		Use:   "repair",
		Short: "Upload the given manifest in the storage and generate the actions required to make the storage's game files an accurate reflection of the manifest",
		PreRun: func(cmd *cobra.Command, args []string) {
			var err error
			m, err = loadManifestFromFile(manifestPath)
			processError(err)
		},
		Run: func(cmd *cobra.Command, args []string) {
			gamesStorage, _ := getStorage(path, storageType, logSource, "")
			err := storage.Repair(&m, gamesStorage, storage.Source{Type: "gog"}, concurrency)
			processError(err)
		},
	}

	storageRepairCmd.Flags().StringVarP(&manifestPath, "manifest", "m", "manifest.json", "Manifest that you want to upload in your storage")
	storageRepairCmd.MarkFlagFilename("manifest")
	storageRepairCmd.Flags().StringVarP(&path, "path", "p", "games", "Path to your games' storage (directory if it is of type fs, json configuration file if it is of type s3)")
	storageRepairCmd.Flags().StringVarP(&storageType, "storage", "k", "fs", "The type of storage you are using. Can be 'fs' (for file system) or 's3' (for s3 store)")
	storageRepairCmd.Flags().IntVarP(&concurrency, "concurrency", "r", 10, "Number of manifest games that should be processed at the same time")

	return storageRepairCmd
}

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
	var verifyChecksum bool
	var useFileManifest bool

	storageRepairCmd := &cobra.Command{
		Use:   "repair",
		Short: "Scan the storage's game files and generate the remedial actions required to ensure the storage's game files are an accurate reflection of the manifest",
		PreRun: func(cmd *cobra.Command, args []string) {
			var err error

			if useFileManifest {
				m, err = loadManifestFromFile(manifestPath)
				processError(err)
				return
			}

			gamesStorage, _ := getStorage(path, storageType, logSource, "")
			storageManifest, err := gamesStorage.LoadManifest()
			m = (*storageManifest)
		},
		Run: func(cmd *cobra.Command, args []string) {
			gamesStorage, _ := getStorage(path, storageType, logSource, "")
			err := storage.Repair(&m, gamesStorage, storage.Source{Type: "gog"}, concurrency, verifyChecksum)
			processError(err)
		},
	}

	storageRepairCmd.Flags().BoolVarP(&useFileManifest, "file-manifest", "f", false, "If set to true, a specified manifest file will be used to repair the storage, otherwise the storage's manifest will be used")
	storageRepairCmd.Flags().StringVarP(&manifestPath, "manifest", "m", "manifest.json", "Manifest file that you want to use to repair the storage")
	storageRepairCmd.MarkFlagFilename("manifest")
	storageRepairCmd.Flags().StringVarP(&path, "path", "p", "games", "Path to your games' storage (directory if it is of type fs, json configuration file if it is of type s3)")
	storageRepairCmd.Flags().StringVarP(&storageType, "storage", "k", "fs", "The type of storage you are using. Can be 'fs' (for file system) or 's3' (for s3 store)")
	storageRepairCmd.Flags().IntVarP(&concurrency, "concurrency", "r", 10, "Number of manifest games that should be processed at the same time")
    storageRepairCmd.Flags().BoolVarP(&verifyChecksum, "verify-checksum", "v", false, "If set to true, checksum comparison of files against the manifest checksum value will be performed")

	return storageRepairCmd
}

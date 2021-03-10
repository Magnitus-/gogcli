package cmd

import (
	"gogcli/manifest"
	"gogcli/storage"

	"github.com/spf13/cobra"
)

func generateStorageUpdateActionsCmd() *cobra.Command {
	var m manifest.Manifest
	var manifestPath string
	var path string
	var storageType string
	var allowEmptyCheckum bool

	storageUpdateActionsCmd := &cobra.Command{
		Use:   "update-actions",
		Short: "Update the manifest and uncompleted actions in a storage given a newer manifest",
		PreRun: func(cmd *cobra.Command, args []string) {
			var err error
			m, err = loadManifestFromFile(manifestPath)
			processError(err)
		},
		Run: func(cmd *cobra.Command, args []string) {
			gamesStorage, _ := getStorage(path, storageType, debugMode, "")
			err := storage.UpdateActions(&m, gamesStorage, allowEmptyCheckum)
			processError(err)
		},
	}

	storageUpdateActionsCmd.Flags().StringVarP(&manifestPath, "manifest", "m", "manifest.json", "Path were the manifest you want to update the storage's actions with is")
	storageUpdateActionsCmd.MarkFlagFilename("manifest")
	storageUpdateActionsCmd.Flags().StringVarP(&path, "path", "p", "games", "Path to your games' storage (directory if it is of type fs, json configuration file if it is of type s3)")
	storageUpdateActionsCmd.Flags().StringVarP(&storageType, "storage", "k", "fs", "The type of storage you are using. Can be 'fs' (for file system) or 's3' (for s3 store)")
	storageUpdateActionsCmd.Flags().BoolVarP(&allowEmptyCheckum, "empty-checksum", "s", false, "If set to true, manifest files with empty checksums will count as already uploaded if everything else matches")

	return storageUpdateActionsCmd
}
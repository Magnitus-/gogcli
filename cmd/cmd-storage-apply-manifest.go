package cmd

import (
	"gogcli/manifest"
	"gogcli/storage"

	"github.com/spf13/cobra"
)

func generateStorageApplyManifestCmd() *cobra.Command {
	var m manifest.Manifest
	var manifestPath string
	var path string
	var storageType string
	var allowEmptyCheckum bool

	storageApplyManifestCmd := &cobra.Command{
		Use:   "manifest",
		Short: "Applies a given manifest into a storage, generating (and potentially modify existing) actions which will need to be executed to make the game files in the storage like the manifest",
		PreRun: func(cmd *cobra.Command, args []string) {
			var err error
			m, err = loadManifestFromFile(manifestPath)
			processError(err)
		},
		Run: func(cmd *cobra.Command, args []string) {
			gamesStorage, _ := getStorage(path, storageType, logSource, "")

			err := storage.EnsureInitialization(gamesStorage)
			processError(err)

			err = storage.ApplyManifest(&m, gamesStorage, storage.Source{Type: "gog"}, allowEmptyCheckum)
			processError(err)			
		},
	}

	storageApplyManifestCmd.Flags().StringVarP(&manifestPath, "manifest", "m", "manifest.json", "Path were the manifest you want to apply is")
	storageApplyManifestCmd.MarkFlagFilename("manifest")
	storageApplyManifestCmd.Flags().StringVarP(&path, "path", "p", "games", "Path to the directory where game files should be stored")
	storageApplyManifestCmd.Flags().StringVarP(&storageType, "storage", "k", "fs", "The type of storage you are using. Can be 'fs' (for file system) or 's3' (for s3 store)")
	storageApplyManifestCmd.Flags().BoolVarP(&allowEmptyCheckum, "empty-checksum", "s", false, "If set to true, manifest files with empty checksums will count as already uploaded if everything else matches")
	return storageApplyManifestCmd
}
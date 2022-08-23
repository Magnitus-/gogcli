package cmd

import (
	"errors"
	"fmt"
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
	var allowGameDeletions bool

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

			err = storage.ImprintProtectedFiles(&m, gamesStorage)
			processError(err)

			if !allowGameDeletions {
				checksumValidation := manifest.ChecksumValidation
				if allowEmptyCheckum {
					checksumValidation = manifest.ChecksumValidationIfPresent
				}

				actions, err := storage.PlanManifest(&m, gamesStorage, checksumValidation)
				processError(err)
				summary := actions.GetSummary()
				if summary.GameDeletions > 0 {
					processError(errors.New(fmt.Sprintf("Executing the action would result in the deletion of %d games, aborting.", summary.GameDeletions)))
				}
			}

			err = storage.ApplyManifest(&m, gamesStorage, storage.Source{Type: "gog"}, allowEmptyCheckum)
			processError(err)
		},
	}

	storageApplyManifestCmd.Flags().StringVarP(&manifestPath, "manifest", "m", "manifest.json", "Path were the manifest you want to apply is")
	storageApplyManifestCmd.MarkFlagFilename("manifest")
	storageApplyManifestCmd.Flags().StringVarP(&path, "path", "p", "games", "Path to the directory where game files should be stored")
	storageApplyManifestCmd.Flags().StringVarP(&storageType, "storage", "k", "fs", "The type of storage you are using. Can be 'fs' (for file system) or 's3' (for s3 store)")
	storageApplyManifestCmd.Flags().BoolVarP(&allowEmptyCheckum, "empty-checksum", "s", false, "If set to true, manifest files with empty checksums will count as already uploaded if everything else matches")
	storageApplyManifestCmd.Flags().BoolVarP(&allowGameDeletions, "allow-game-deletions", "d", false, "If set to true, an actions file that contain game deletion actions will be allowed, otherwise the command will abort if this would be the result")
	return storageApplyManifestCmd
}

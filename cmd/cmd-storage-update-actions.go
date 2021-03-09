package cmd

import (
	"github.com/spf13/cobra"
)

func generateStorageUpdateActionsCmd() *cobra.Command {
	var m manifest.Manifest
	var manifestPath string
	
	storageUpdateActionsCmd := &cobra.Command{
		Use:   "update-actions",
		Short: "Update the manifest and uncompleted actions in a storage given a newer manifest",
		PreRun: func(cmd *cobra.Command, args []string) {
			var err error
			m, err = loadManifestFromFile(manifestPath)
			processError(err)
		},
		Run: func(cmd *cobra.Command, args []string) {
		}
	}

	return storageUpdateActionsCmd
}
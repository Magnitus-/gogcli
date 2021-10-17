package cmd

import (
	"errors"
	"github.com/spf13/cobra"
)

func generateStorageDownloadManifestCmd() *cobra.Command {
	var path string
	var storageType string
	var file string
	var terminalOutput bool

	storageDownloadManifestCmd := &cobra.Command{
		Use:   "manifest",
		Short: "Commands to download the manifest file from the storage",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			callPersistentPreRun(cmd, args)
		},
		Run: func(cmd *cobra.Command, args []string) {
			gamesStorage, _ := getStorage(path, storageType, logSource, "")

			exists, err := gamesStorage.Exists()
			processError(err)
			if !exists {
				processError(errors.New("Specified storage doesn't exist"))
			}

			has, hasErr := gamesStorage.HasManifest()
			processError(hasErr)
			if !has {
				processError(errors.New("Specified storage doesn't have a manifest"))
			}

			m, mErr := gamesStorage.LoadManifest()
			processError(mErr)
			processSerializableOutput(m, []error{}, terminalOutput, file)
		},
	}

	storageDownloadManifestCmd.Flags().StringVarP(&file, "file", "f", "manifest.json", "File to output the manifest in")
	storageDownloadManifestCmd.Flags().BoolVarP(&terminalOutput, "terminal", "t", false, "If set to true, the manifest will be output on the terminal instead of in a file")
	storageDownloadManifestCmd.Flags().StringVarP(&path, "path", "p", "games", "Path to your games' storage (directory if it is of type fs, json configuration file if it is of type s3)")
	storageDownloadManifestCmd.Flags().StringVarP(&storageType, "storage", "k", "fs", "The type of storage you are using. Can be 'fs' (for file system) or 's3' (for s3 store)")

	return storageDownloadManifestCmd
}

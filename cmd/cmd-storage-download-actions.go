package cmd

import (
	"errors"
	"github.com/spf13/cobra"
)

func generateStorageDownloadActionsCmd() *cobra.Command {
	var path string
	var storageType string
	var file string
	var terminalOutput bool

	storageDownloadActionsCmd := &cobra.Command{
		Use:   "actions",
		Short: "Commands to download the actions file from the storage",
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

			has, hasErr := gamesStorage.HasActions()
			processError(hasErr)
			if !has {
				processError(errors.New("Specified storage doesn't have actions"))
			}

			a, mErr := gamesStorage.LoadActions()
			processError(mErr)
			processSerializableOutput(a, []error{}, terminalOutput, file)
		},
	}

	storageDownloadActionsCmd.Flags().StringVarP(&file, "file", "f", "actions.json", "File to output the actions in")
	storageDownloadActionsCmd.Flags().BoolVarP(&terminalOutput, "terminal", "t", false, "If set to true, the actions will be output on the terminal instead of in a file")
	storageDownloadActionsCmd.Flags().StringVarP(&path, "path", "p", "games", "Path to your games' storage (directory if it is of type fs, json configuration file if it is of type s3)")
	storageDownloadActionsCmd.Flags().StringVarP(&storageType, "storage", "k", "fs", "The type of storage you are using. Can be 'fs' (for file system) or 's3' (for s3 store)")

	return storageDownloadActionsCmd
}

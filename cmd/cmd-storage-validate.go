package cmd

import (
	"fmt"
	"gogcli/storage"
	"os"

	"github.com/spf13/cobra"
)

func generateStorageValidateFsCmd(concurrency *int) *cobra.Command {
	var path string

	storageValidateFsCmd := &cobra.Command{
		Use:   "fs",
		Short: "Use a file system storage",
		Run: func(cmd *cobra.Command, args []string) {
			errs := storage.ValidateManifest(storage.GetFileSystem(path, debugMode, ""), (*concurrency))
			if len(errs) > 0 {
				for _, err := range errs {
					fmt.Println(err)
				}
				os.Exit(1)
			}
		},
	}

	storageValidateFsCmd.Flags().StringVarP(&path, "path", "p", "games", "Path to the directory where game files are stored")
	return storageValidateFsCmd
}

func generateStorageValidateCmd() *cobra.Command {
	var concurrency int

	storageValidateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate that all the game files in the storage match the size and checksum values in the manifest",
	}

	storageValidateCmd.PersistentFlags().IntVarP(&concurrency, "concurrency", "r", 10, "Number of downloads that should be attempted at the same time")

	storageValidateCmd.AddCommand(generateStorageValidateFsCmd(&concurrency))

	return storageValidateCmd
}
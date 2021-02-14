package cmd

import (
	"fmt"
	"gogcli/storage"
	"os"

	"github.com/spf13/cobra"
)

func generateStorageValidateCmd() *cobra.Command {
	var concurrency int
	var path string
	var storageType string

	storageValidateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate that all the game files in the storage match the size and checksum values in the manifest",
		Run: func(cmd *cobra.Command, args []string) {
			gamesStorage, _ := getStorage(path, storageType, debugMode, "")
			
			errs := storage.ValidateManifest(gamesStorage, concurrency)
			if len(errs) > 0 {
				for _, err := range errs {
					fmt.Println(err)
				}
				os.Exit(1)
			}
		},
	}

	storageValidateCmd.Flags().IntVarP(&concurrency, "concurrency", "r", 10, "Number of downloads that should be attempted at the same time")
	storageValidateCmd.Flags().StringVarP(&path, "path", "p", "games", "Path to your games' storage (directory if it is of type fs, json configuration file if it is of type s3)")
	storageValidateCmd.Flags().StringVarP(&storageType, "storage", "k", "fs", "The type of storage you are using. Can be 'fs' (for file system) or 's3' (for s3 store)")

	return storageValidateCmd
}
package cmd

import (
	"gogcli/storage"

	"github.com/spf13/cobra"
)

func generateStorageCopyCmd() *cobra.Command {
	var concurrency int
	var sourcePath string
	var destinationPath string
	var sourceStorage string
	var destinationStorage string

	storageCopyCmd := &cobra.Command{
		Use:   "copy",
		Short: "Copy the game files from one storage to another",
		Run: func(cmd *cobra.Command, args []string) {
			source, downloader := getStorage(sourcePath, sourceStorage, debugMode, "SOURCE")
			destination, _ := getStorage(destinationPath, destinationStorage, debugMode, "DESTINATION")

			errs := storage.Copy(source, destination, downloader, concurrency)
			processErrors(errs)
		},
	}

	storageCopyCmd.Flags().IntVarP(&concurrency, "concurrency", "r", 10, "Number of downloads that should be attempted at the same time")
	storageCopyCmd.Flags().StringVarP(&sourcePath, "source-path", "s", "games", "Path to the source of your games (directory if it is of type fs, json configuration file if it is of type s3)")
	storageCopyCmd.Flags().StringVarP(&sourceStorage, "source-storage", "t", "fs", "Kind of storage your source is. Can be 'fs' (for file system) or 's3' (for s3 store)")
	storageCopyCmd.Flags().StringVarP(&destinationPath, "destination-path", "n", "games-copy", "Path to the destination of your games (directory if it is of type fs, json configuration file if it is of type s3)")
	storageCopyCmd.Flags().StringVarP(&destinationStorage, "destination-storage", "o", "fs", "Kind of storage your destination is. Can be 'fs' (for file system) or 's3' (for s3 store)")

	return storageCopyCmd
}

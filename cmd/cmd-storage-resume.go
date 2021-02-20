package cmd

import (
	"gogcli/sdk"
	"gogcli/storage"

	"github.com/spf13/cobra"
)

func generateStorageResumeCmd() *cobra.Command {
	var path string
	var storageType string
	var gamesMax int
	var concurrency int

	storageResumeCmd := &cobra.Command{
		Use:   "resume",
		Short: "Resume a manifest apply operation that didn't complete on a storage",
		Run: func(cmd *cobra.Command, args []string) {
			var downloader storage.Downloader
			gamesStorage, _ := getStorage(path, storageType, debugMode, "DESTINATION")
			
			source, err := gamesStorage.LoadSource()
			if err != nil {
				processError(err)
			}
			if source.Type == "gog" {
				downloader = sdk.Downloader{sdkPtr}
			} else if source.Type == "fs" {
				fs, sourceErr := storage.GetFileSystemFromSource(*source ,debugMode, "SOURCE")
				if sourceErr != nil {
					processError(sourceErr)
				}
				downloader = storage.FileSystemDownloader{fs}
			} else {
				s3, sourceErr := storage.GetS3StoreFromSource(*source ,debugMode, "SOURCE")
				if sourceErr != nil {
					processError(sourceErr)
				}
				downloader = storage.S3StoreDownloader{s3}
			}
			
			errs := storage.ResumeActions(gamesStorage, concurrency, downloader, gamesMax)
			processErrors(errs)
		},
	}

	storageResumeCmd.Flags().StringVarP(&path, "path", "p", "games", "Path to your games' storage (directory if it is of type fs, json configuration file if it is of type s3)")
	storageResumeCmd.Flags().StringVarP(&storageType, "storage", "k", "fs", "The type of storage you are using. Can be 'fs' (for file system) or 's3' (for s3 store)")
	storageResumeCmd.Flags().IntVarP(&gamesMax, "maximum", "x", -1, "The maximum number of games to upload into storage.")
	storageResumeCmd.Flags().IntVarP(&concurrency, "concurrency", "r", 10, "Number of downloads that should be attempted at the same time")

	return storageResumeCmd
}

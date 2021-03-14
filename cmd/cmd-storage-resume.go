package cmd

import (
	"gogcli/manifest"
	"gogcli/sdk"
	"gogcli/storage"

	"github.com/spf13/cobra"
)

func generateStorageResumeCmd() *cobra.Command {
	var path string
	var storageType string
	var gamesMax int
	var concurrency int
	var preferredGameIds []int64
	var sortCriterion string
	var sortAscending bool

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
			
			sort := manifest.NewActionIteratorSort(preferredGameIds, sortCriterion, sortAscending)
			errs := storage.ResumeActions(gamesStorage, concurrency, downloader, gamesMax, sort)
			processErrors(errs)
		},
	}

	storageResumeCmd.Flags().StringVarP(&path, "path", "p", "games", "Path to your games' storage (directory if it is of type fs, json configuration file if it is of type s3)")
	storageResumeCmd.Flags().StringVarP(&storageType, "storage", "k", "fs", "The type of storage you are using. Can be 'fs' (for file system) or 's3' (for s3 store)")
	storageResumeCmd.Flags().IntVarP(&gamesMax, "maximum", "x", -1, "The maximum number of games to upload into storage.")
	storageResumeCmd.Flags().IntVarP(&concurrency, "concurrency", "r", 10, "Number of downloads that should be attempted at the same time")
	storageResumeCmd.Flags().Int64SliceVarP(&preferredGameIds, "preferred-ids", "f", []int64{}, "Ids of games to download first")
	storageResumeCmd.Flags().StringVarP(&sortCriterion, "sort-criterion", "t", "none", "Criteria to sort games download order by. Can be: id, title, size, none")
	storageResumeCmd.Flags().BoolVarP(&sortAscending, "ascending", "n", true, "If set to true, game downloads will be sorted in ascending order given the sort criterion")

	return storageResumeCmd
}

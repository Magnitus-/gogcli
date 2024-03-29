package cmd

import (
	"gogcli/manifest"
	"gogcli/sdk"
	"gogcli/storage"

	"github.com/spf13/cobra"
)

func generateStorageExecuteActionsCmd() *cobra.Command {
	var path string
	var storageType string
	var gamesMax int
	var concurrency int
	var preferredGameIds []int64
	var sortCriterion string
	var sortAscending bool
	var downloadRetries int

	storageExecuteActionsCmd := &cobra.Command{
		Use:   "execute-actions",
		Short: "Runs any pending actions in the storage",
		Run: func(cmd *cobra.Command, args []string) {
			var downloader storage.Downloader
			gamesStorage, _ := getStorage(path, storageType, logSource, "destination")

			source, err := gamesStorage.LoadSource()
			if err != nil {
				processError(err)
			}
			if source.Type == "gog" {
				downloader = sdk.Downloader{sdkPtr}
			} else if source.Type == "fs" {
				fs, sourceErr := storage.GetFileSystemFromSource(*source, logSource, "source")
				if sourceErr != nil {
					processError(sourceErr)
				}
				downloader = storage.FileSystemDownloader{fs}
			} else {
				s3, sourceErr := storage.GetS3StoreFromSource(*source, logSource, "source")
				if sourceErr != nil {
					processError(sourceErr)
				}
				downloader = storage.S3StoreDownloader{s3}
			}

			sort := manifest.NewActionIteratorSort(preferredGameIds, sortCriterion, sortAscending)
			proc := storage.GetActionsProcessor(concurrency, downloadRetries, gamesMax, sort, logSource)
			errs := storage.ExecuteActions(gamesStorage, downloader, proc)
			processErrors(errs)
		},
	}

	storageExecuteActionsCmd.Flags().StringVarP(&path, "path", "p", "games", "Path to your games' storage (directory if it is of type fs, json configuration file if it is of type s3)")
	storageExecuteActionsCmd.Flags().StringVarP(&storageType, "storage", "k", "fs", "The type of storage you are using. Can be 'fs' (for file system) or 's3' (for s3 store)")
	storageExecuteActionsCmd.Flags().IntVarP(&gamesMax, "maximum", "x", -1, "The maximum number of games to upload into storage.")
	storageExecuteActionsCmd.Flags().IntVarP(&concurrency, "concurrency", "r", 4, "Number of downloads that should be attempted at the same time")
	storageExecuteActionsCmd.Flags().Int64SliceVarP(&preferredGameIds, "preferred-ids", "f", []int64{}, "Ids of games to download first")
	storageExecuteActionsCmd.Flags().StringVarP(&sortCriterion, "sort-criterion", "t", "none", "Criteria to sort games download order by. Can be: id, title, size, none")
	storageExecuteActionsCmd.Flags().BoolVarP(&sortAscending, "ascending", "n", true, "If set to true, game downloads will be sorted in ascending order given the sort criterion")
	storageExecuteActionsCmd.Flags().IntVarP(&downloadRetries, "download-retries", "d", 2, "How many times to retry a failed download before giving up")

	return storageExecuteActionsCmd
}

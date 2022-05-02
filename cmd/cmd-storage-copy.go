package cmd

import (
	"gogcli/manifest"
	"gogcli/storage"

	"github.com/spf13/cobra"
)

func generateStorageCopyCmd() *cobra.Command {
	var concurrency int
	var sourcePath string
	var destinationPath string
	var sourceStorage string
	var destinationStorage string
	var gamesMax int
	var preferredGameIds []int64
	var sortCriterion string
	var sortAscending bool
	var downloadRetries int

	storageCopyCmd := &cobra.Command{
		Use:   "copy",
		Short: "Copy the game files from one storage to another",
		Run: func(cmd *cobra.Command, args []string) {
			source, downloader := getStorage(sourcePath, sourceStorage, logSource, "source")
			destination, _ := getStorage(destinationPath, destinationStorage, logSource, "destination")

			sort := manifest.NewActionIteratorSort(preferredGameIds, sortCriterion, sortAscending)
			proc := storage.GetActionsProcessor(concurrency, downloadRetries, gamesMax, sort, logSource)
			errs := storage.Copy(source, destination, downloader, proc)
			processErrors(errs)
		},
	}

	storageCopyCmd.Flags().IntVarP(&concurrency, "concurrency", "r", 4, "Number of downloads that should be attempted at the same time")
	storageCopyCmd.Flags().StringVarP(&sourcePath, "source-path", "s", "games", "Path to the source of your games (directory if it is of type fs, json configuration file if it is of type s3)")
	storageCopyCmd.Flags().StringVarP(&sourceStorage, "source-storage", "t", "fs", "Kind of storage your source is. Can be 'fs' (for file system) or 's3' (for s3 store)")
	storageCopyCmd.Flags().StringVarP(&destinationPath, "destination-path", "n", "games-copy", "Path to the destination of your games (directory if it is of type fs, json configuration file if it is of type s3)")
	storageCopyCmd.Flags().StringVarP(&destinationStorage, "destination-storage", "o", "fs", "Kind of storage your destination is. Can be 'fs' (for file system) or 's3' (for s3 store)")
	storageCopyCmd.Flags().IntVarP(&gamesMax, "maximum", "x", -1, "The maximum number of games to copy into storage.")
	storageCopyCmd.Flags().Int64SliceVarP(&preferredGameIds, "preferred-ids", "f", []int64{}, "Ids of games to download first")
	storageCopyCmd.Flags().StringVarP(&sortCriterion, "sort-criterion", "i", "none", "Criteria to sort games download order by. Can be: id, title, size, none")
	storageCopyCmd.Flags().BoolVarP(&sortAscending, "ascending", "a", true, "If set to true, game downloads will be sorted in ascending order given the sort criterion")
	storageCopyCmd.Flags().IntVarP(&downloadRetries, "download-retries", "d", 2, "How many times to retry a failed download before giving up")

	return storageCopyCmd
}

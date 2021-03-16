package cmd

import (
	"bytes"
	"encoding/json"
	"gogcli/manifest"
	"gogcli/sdk"
	"gogcli/storage"
	"io/ioutil"

	"github.com/spf13/cobra"
)

func generateStorageApplyCmd() *cobra.Command {
	var m manifest.Manifest
	var manifestPath string
	var concurrency int
	var path string
	var storageType string
	var gamesMax int
	var allowEmptyCheckum bool
	var preferredGameIds []int64
	var sortCriterion string
	var sortAscending bool

	storageApplyCmd := &cobra.Command{
		Use:   "apply",
		Short: "Change the files in a given storage to match the content of a manifest, uploading and deleting files as necessary",
		PreRun: func(cmd *cobra.Command, args []string) {
			var err error
			m, err = loadManifestFromFile(manifestPath)
			processError(err)
		},
		Run: func(cmd *cobra.Command, args []string) {
			var buf bytes.Buffer
			gamesStorage, _ := getStorage(path, storageType, logSource, "")
			downloader := sdk.Downloader{sdkPtr}

			err := storage.EnsureInitialization(gamesStorage)
			processError(err)

			sort := manifest.NewActionIteratorSort(preferredGameIds, sortCriterion, sortAscending)
			errs := storage.UploadManifest(&m, gamesStorage, storage.Source{Type: "gog"}, concurrency, downloader, gamesMax, sort, allowEmptyCheckum)
			processErrors(errs)

			output, _ := json.Marshal(&m)
			json.Indent(&buf, output, "", "  ")
			output = buf.Bytes()
			err = ioutil.WriteFile(manifestPath, output, 0644)
			processError(err)				
		},
	}

	storageApplyCmd.Flags().StringVarP(&manifestPath, "manifest", "m", "manifest.json", "Path were the manifest you want to apply is")
	storageApplyCmd.MarkFlagFilename("manifest")
	storageApplyCmd.Flags().IntVarP(&concurrency, "concurrency", "r", 10, "Number of downloads that should be attempted at the same time")
	storageApplyCmd.Flags().StringVarP(&path, "path", "p", "games", "Path to the directory where game files should be stored")
	storageApplyCmd.Flags().StringVarP(&storageType, "storage", "k", "fs", "The type of storage you are using. Can be 'fs' (for file system) or 's3' (for s3 store)")
	storageApplyCmd.Flags().IntVarP(&gamesMax, "maximum", "x", -1, "The maximum number of games to upload into storage.")
	storageApplyCmd.Flags().BoolVarP(&allowEmptyCheckum, "empty-checksum", "s", false, "If set to true, manifest files with empty checksums will count as already uploaded if everything else matches")
	storageApplyCmd.Flags().Int64SliceVarP(&preferredGameIds, "preferred-ids", "f", []int64{}, "Ids of games to download first")
	storageApplyCmd.Flags().StringVarP(&sortCriterion, "sort-criterion", "t", "none", "Criteria to sort games download order by. Can be: id, title, size, none")
	storageApplyCmd.Flags().BoolVarP(&sortAscending, "ascending", "n", true, "If set to true, game downloads will be sorted in ascending order given the sort criterion")
	return storageApplyCmd
}

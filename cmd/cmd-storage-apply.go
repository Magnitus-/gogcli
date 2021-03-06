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

	storageApplyCmd := &cobra.Command{
		Use:   "apply",
		Short: "Change the files in a given storage to match the content of a manifest, uploading and deleting files as necessary",
		PreRun: func(cmd *cobra.Command, args []string) {
			var err error
			m, err = loadManifestFromFile(path)
			processError(err)
		},
		Run: func(cmd *cobra.Command, args []string) {
			var buf bytes.Buffer
			gamesStorage, _ := getStorage(path, storageType, debugMode, "")
			downloader := sdk.Downloader{sdkPtr}

			err := storage.EnsureInitialization(gamesStorage)
			processError(err)

			errs := storage.UploadManifest(&m, gamesStorage, storage.Source{Type: "gog"}, concurrency, downloader, gamesMax, allowEmptyCheckum)
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

	return storageApplyCmd
}

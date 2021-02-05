package cmd

import (
	"encoding/json"
	"fmt"
	"gogcli/manifest"
	"gogcli/storage"
	"io"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

func generateApplyManifestFsCmd(m *manifest.Manifest, concurrency *int, fn GetDownloadHandle) *cobra.Command {
	var path string

	applyManifestFsCmd := &cobra.Command{
		Use:   "fs",
		Short: "Apply the manifest on a file system",
		PreRun: func(cmd *cobra.Command, args []string) {
			_, err := os.Stat(path)
			if err != nil {
				if os.IsNotExist(err) {
					err = os.MkdirAll(path, 0755)
					if err != nil {
						fmt.Println("Failed to create a directory at the specified path: ", err)
						os.Exit(1)
					}
				} else {
					fmt.Println("An error occured while trying to gather info on the specified path: ", err)
					os.Exit(1)
				}
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			fs := storage.GetFileSystem(path, debugMode)
			errs := uploadManifest(m, fs, *concurrency, fn)
			if len(errs) > 0 {
				for _, err := range errs {
					fmt.Println(err)
				}
				os.Exit(1)
			}
		},
	}

	applyManifestFsCmd.Flags().StringVarP(&path, "path", "p", "games", "Path to the directory where game files should be stored")

	return applyManifestFsCmd
}

func generateApplyManifestCmd() *cobra.Command {
	var m manifest.Manifest
	var manifestPath string
	var concurrency int

	download := func(gameId int, add manifest.FileAction) (io.ReadCloser, int, string, error) {
		return sdkPtr.GetDownloadHandle(add.Url, debugMode)
	}

	applyManifestCmd := &cobra.Command{
		Use:   "apply-manifest",
		Short: "Change the files in a given storage to match the content of a manifest, uploading and deleting files as necessary",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			bs, err := ioutil.ReadFile(manifestPath)
			if err != nil {
				fmt.Println("Could not load the manifest: ", err)
				os.Exit(1)
			}

			err = json.Unmarshal(bs, &m)
			if err != nil {
				fmt.Println("Manifest file doesn't appear to contain valid json: ", err)
				os.Exit(1)
			}
		},
	}

	applyManifestCmd.PersistentFlags().StringVarP(&manifestPath, "manifest", "m", "manifest.json", "Path were the manifest you want to apply is")
	applyManifestCmd.MarkPersistentFlagFilename("manifest")
	applyManifestCmd.PersistentFlags().IntVarP(&concurrency, "concurrency", "r", 10, "Number of downloads that should be attempted at the same time")

	applyManifestCmd.AddCommand(generateApplyManifestFsCmd(&m, &concurrency, download))

	return applyManifestCmd
}

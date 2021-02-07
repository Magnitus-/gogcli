package cmd

import (
	"fmt"
	"gogcli/storage"
	"os"

	"github.com/spf13/cobra"
)

func generateCopyStorageFsToFsCmd(concurrency *int) *cobra.Command {
	var sourcePath string
	var destinationPath string

	copyStorageFsToFsCmd := &cobra.Command{
		Use:   "fs-to-fs",
		Short: "Transfer the files stored in one filesystem to another",
		PreRun: func(cmd *cobra.Command, args []string) {
			_, err := os.Stat(sourcePath)
			if err != nil {
				fmt.Println("An error occured while trying to gather info on the source path: ", err)
				os.Exit(1)
			}

			_, err = os.Stat(destinationPath)
			if err != nil {
				if os.IsNotExist(err) {
					err = os.MkdirAll(destinationPath, 0755)
					if err != nil {
						fmt.Println("Failed to create a directory for the destination path: ", err)
						os.Exit(1)
					}
				} else {
					fmt.Println("An error occured while trying to gather info on the destination path: ", err)
					os.Exit(1)
				}
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			sourceFs := storage.GetFileSystem(sourcePath, debugMode, "SOURCE")
			destinationFs := storage.GetFileSystem(destinationPath, debugMode, "DESTINATION")
			m, err := sourceFs.LoadManifest()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			errs := uploadManifest(m, destinationFs, *concurrency, storage.FileSystemDownloader{sourceFs})
			if len(errs) > 0 {
				for _, err := range errs {
					fmt.Println(err)
				}
				os.Exit(1)
			}
		},
	}

	copyStorageFsToFsCmd.Flags().StringVarP(&sourcePath, "source-path", "s", "games", "Path to the directory where game files to copy are stored")
	copyStorageFsToFsCmd.Flags().StringVarP(&destinationPath, "destination-path", "n", "games-copy", "Path to the directory where game files from source are to be copied")

	return copyStorageFsToFsCmd
}

func generateCopyStorageCmd() *cobra.Command {
	var concurrency int

	copyStorageCmd := &cobra.Command{
		Use:   "copy-storage",
		Short: "Copy the game files from one storage to another",
	}

	copyStorageCmd.PersistentFlags().IntVarP(&concurrency, "concurrency", "r", 10, "Number of downloads that should be attempted at the same time")

	copyStorageCmd.AddCommand(generateCopyStorageFsToFsCmd(&concurrency))

	return copyStorageCmd
}

package cmd

import (
	"fmt"
	"gogcli/storage"
	"os"

	"github.com/spf13/cobra"
)

func processErrors(errs []error) {
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
		os.Exit(1)
	}	
}

func generateStorageCopyFsToS3Cmd(concurrency *int) *cobra.Command {
	var sourcePath string
	var destinationConfigsPath string

	storageCopyFsToS3Cmd := &cobra.Command{
		Use:   "fs-to-s3",
		Short: "The source storage is a file system and the destination storage is an s3 store",
		Run: func(cmd *cobra.Command, args []string) {
			source := storage.GetFileSystem(sourcePath, debugMode, "SOURCE")
			destination, err := storage.GetS3StoreFromConfigFile(destinationConfigsPath, debugMode, "DESTINATION")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			errs := storage.Copy(source, destination, storage.FileSystemDownloader{source}, *concurrency)
			processErrors(errs)
		},
	}

	storageCopyFsToS3Cmd.Flags().StringVarP(&sourcePath, "source-path", "s", "games", "Path to the directory where game files to copy are stored")
	storageCopyFsToS3Cmd.Flags().StringVarP(&destinationConfigsPath, "destination-configs-path", "n", "copy-destination.json", "Path to a json configuration file for your S3 store. The following values are expected: Endpoint (string), Region (string), Bucket (string), Tls (boolean), AccessKey (string), SecretKey (string)")

	return storageCopyFsToS3Cmd
}

func generateStorageCopyFsToFsCmd(concurrency *int) *cobra.Command {
	var sourcePath string
	var destinationPath string

	storageCopyFsToFsCmd := &cobra.Command{
		Use:   "fs-to-fs",
		Short: "The source storage is a file system and the destination storage is a file system",
		Run: func(cmd *cobra.Command, args []string) {
			source := storage.GetFileSystem(sourcePath, debugMode, "SOURCE")
			destination := storage.GetFileSystem(destinationPath, debugMode, "DESTINATION")
			errs := storage.Copy(source, destination, storage.FileSystemDownloader{source}, *concurrency)
			processErrors(errs)
		},
	}

	storageCopyFsToFsCmd.Flags().StringVarP(&sourcePath, "source-path", "s", "games", "Path to the directory where game files to copy are stored")
	storageCopyFsToFsCmd.Flags().StringVarP(&destinationPath, "destination-path", "n", "games-copy", "Path to the directory where game files from source are to be copied")

	return storageCopyFsToFsCmd
}

func generateStorageCopyCmd() *cobra.Command {
	var concurrency int

	storageCopyCmd := &cobra.Command{
		Use:   "copy",
		Short: "Copy the game files from one storage to another",
	}

	storageCopyCmd.PersistentFlags().IntVarP(&concurrency, "concurrency", "r", 10, "Number of downloads that should be attempted at the same time")

	storageCopyCmd.AddCommand(generateStorageCopyFsToFsCmd(&concurrency))
	storageCopyCmd.AddCommand(generateStorageCopyFsToS3Cmd(&concurrency))
	return storageCopyCmd
}

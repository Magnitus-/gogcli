package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func generateUrlPathInfoCmd() *cobra.Command {
	var path string

	urlPathInfoCmd := &cobra.Command{
		Use:   "url-path-info",
		Short: "Given a download path, retrieve the filename, size and checksum of the file that would be downloaded. Valid paths can be obtained from the manifest.",
		Run: func(cmd *cobra.Command, args []string) {
			filename, checksum, size, err := sdkPtr.GetDownloadFileInfo(path)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Println("File Name: ", filename)
			fmt.Println("Checksum:", checksum)
			fmt.Println("Size:", size)
		},
	}

	urlPathInfoCmd.Flags().StringVarP(&path, "path", "p", "", "Url path to download")
	urlPathInfoCmd.MarkFlagRequired("path")

	return urlPathInfoCmd
}

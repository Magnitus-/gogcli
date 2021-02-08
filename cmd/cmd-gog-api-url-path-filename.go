package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func generateUrlPathFilenameCmd() *cobra.Command {
	var path string

	urlPathFilenameCmd := &cobra.Command{
		Use:   "url-path-filename",
		Short: "Given a download path, retrieve the filename of the file that would be downloaded. Valid paths can be obtained from the manifest.",
		Run: func(cmd *cobra.Command, args []string) {
			filename, err := sdkPtr.GetDownloadFilename(path)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Println(filename)
		},
	}

	urlPathFilenameCmd.Flags().StringVarP(&path, "path", "p", "", "Url path to download")
	urlPathFilenameCmd.MarkFlagRequired("path")

	return urlPathFilenameCmd
}

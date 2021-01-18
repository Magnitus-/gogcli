package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

func generateDownloadPathCmd() *cobra.Command {
	var path string
	var file string

	downloadPathCmd := &cobra.Command{
		Use:   "download-path",
		Short: "Download a single file with the given path from GOG. Valid paths can be obtained from the manifest.",
		Run: func(cmd *cobra.Command, args []string) {
			body, _, err := sdkPtr.GetDownloadHandle(path, debugMode)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			defer body.Close()

			out, fErr := os.Create(file)
			if fErr != nil {
				fmt.Println("Error creating file: ", fErr)
				os.Exit(1)
			}
			defer out.Close()

			_, wErr := io.Copy(out, body)
			if wErr != nil {
				fmt.Println("Error writing to the file: ", wErr)
				os.Exit(1)
			}
		},
	}

	downloadPathCmd.Flags().StringVarP(&path, "path", "p", "", "Url path to download")
	downloadPathCmd.MarkFlagRequired("path")

	downloadPathCmd.Flags().StringVarP(&file, "file", "f", "", "File to put the download in")
	downloadPathCmd.MarkFlagRequired("file")

	return downloadPathCmd
}

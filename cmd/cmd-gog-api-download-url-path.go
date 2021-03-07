package cmd

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/spf13/cobra"
)

func generateDownloadUrlPathCmd() *cobra.Command {
	var relUrl string
	var dir string

	downloadUrlPathCmd := &cobra.Command{
		Use:   "download-url-path",
		Short: "Download a single file with the given path from GOG. Valid paths can be obtained from the manifest.",
		PreRun: func(cmd *cobra.Command, args []string) {
			var err error
			if dir == "" {
				dir, err = os.Getwd()
				processError(err)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			body, _, file, err := sdkPtr.GetDownloadHandle(relUrl)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			defer body.Close()

			out, fErr := os.Create(path.Join(dir, file))
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

	downloadUrlPathCmd.Flags().StringVarP(&relUrl, "path", "p", "", "Url path to download")
	downloadUrlPathCmd.MarkFlagRequired("path")
	downloadUrlPathCmd.Flags().StringVarP(&dir, "directory", "r", "", "Directory to download the file in. will default to the current directory")

	return downloadUrlPathCmd
}

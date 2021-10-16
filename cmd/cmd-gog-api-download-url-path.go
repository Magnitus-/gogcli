package cmd

import (
	"fmt"
	"io"
	"time"
	"os"
	"path"

	"github.com/spf13/cobra"
)

func generateDownloadUrlPathCmd() *cobra.Command {
	var relUrl string
	var dir string
	var perf bool

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
			startTime := time.Now()
			body, size, file, err := sdkPtr.GetDownloadHandle(relUrl)
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
			endTime := time.Now()
			if perf {
				elapsedTime := endTime.Sub(startTime)
				elapsedTimeInSeconds := float64(elapsedTime) / float64(1000000000)
				sizeInMB := float64(size) / float64(1000 * 1000)
				sizeInMb := float64(size*8) / float64(1000 * 1000)

				fmt.Printf("Download time: %.2f seconds.\n", elapsedTimeInSeconds)
				fmt.Printf("Metrics in Megabytes (typical storage size notation):\n")
				fmt.Printf("    Download Size: %.2f MB\n", sizeInMB)
				fmt.Printf("    Download Speed: %.2f MBps\n", sizeInMB / elapsedTimeInSeconds)
				fmt.Printf("Metrics in Megabits (typical networking speed notation):\n")
				fmt.Printf("    Download Size: %.2f Mb\n", sizeInMb)
				fmt.Printf("    Download Speed: %.2f Mbps\n", sizeInMb / elapsedTimeInSeconds)
			}
		},
	}

	downloadUrlPathCmd.Flags().StringVarP(&relUrl, "path", "p", "", "Url path to download")
	downloadUrlPathCmd.MarkFlagRequired("path")
	downloadUrlPathCmd.Flags().StringVarP(&dir, "directory", "r", "", "Directory to download the file in. will default to the current directory")
    downloadUrlPathCmd.Flags().BoolVarP(&perf, "performance", "f", false, "Get performance characteristics of the download")

	return downloadUrlPathCmd
}

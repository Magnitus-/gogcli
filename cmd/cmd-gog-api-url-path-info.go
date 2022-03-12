package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type FileInfo struct {
	Name     string
	Checksum string
	Size     int64
}

func generateUrlPathInfoCmd() *cobra.Command {
	var path string
	var file string
	var terminalOutput bool
	var jsonOutput bool

	urlPathInfoCmd := &cobra.Command{
		Use:   "url-path-info",
		Short: "Given a download path, retrieve the filename, size and checksum of the file that would be downloaded. Valid paths can be obtained from the manifest.",
		Run: func(cmd *cobra.Command, args []string) {
			filename, checksum, size, err, _, _ := sdkPtr.GetDownloadFileInfo(path)
			if !jsonOutput {
				processError(err)
				fmt.Println("File Name:", filename)
				fmt.Println("Checksum:", checksum)
				fmt.Println("Size:", size)
			} else {
				errs := make([]error, 0)
				if err != nil {
					errs = append(errs, err)
				}
				processSerializableOutput(
					FileInfo{
						filename,
						checksum,
						size,
					},
					errs,
					terminalOutput,
					file,
				)
			}
		},
	}

	urlPathInfoCmd.Flags().StringVarP(&path, "path", "p", "", "Url path to download")
	urlPathInfoCmd.MarkFlagRequired("path")
	urlPathInfoCmd.Flags().StringVarP(&file, "file", "f", "url-path-info.json", "File to output the url path information in if in json format")
	urlPathInfoCmd.Flags().BoolVarP(&terminalOutput, "terminal", "t", true, "If set to true and json format is used, the url path information will be output on the terminal instead of in a file")
	urlPathInfoCmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "If set to true, the output will be in json format either on the terminal or in a file. Otherwise, it will be in human readable format on the terminal.")

	return urlPathInfoCmd
}

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
	var tolerateBadFileMetadata bool

	urlPathInfoCmd := &cobra.Command{
		Use:   "url-path-info",
		Short: "Given a download path, retrieve the filename, size and checksum of the file that would be downloaded. Valid paths can be obtained from the manifest.",
		Run: func(cmd *cobra.Command, args []string) {
			info := sdkPtr.GetFileInfo(path, tolerateBadFileMetadata)
			if !jsonOutput {
				processError(info.Error)
				fmt.Println("File Name:", info.Name)
				fmt.Println("Checksum:", info.Checksum)
				fmt.Println("Size:", info.Size)
			} else {
				errs := make([]error, 0)
				if info.Error != nil {
					errs = append(errs, info.Error)
				}
				processSerializableOutput(
					FileInfo{
						info.Name,
						info.Checksum,
						info.Size,
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
	urlPathInfoCmd.Flags().BoolVarP(&tolerateBadFileMetadata, "tolerate-bad-metadata", "b", true, "Tolerate files for which metadata cannot be retrieved. The checksum will be infered by performing a throwaway file download instead.")

	return urlPathInfoCmd
}

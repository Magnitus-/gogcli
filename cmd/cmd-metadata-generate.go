package cmd

import (
	"github.com/spf13/cobra"
)

func generateMetadataGenerateCmd() *cobra.Command {
	var concurrency int
	var pause int
	var file string
	var terminalOutput bool
	var tolerateDangles bool
	var warningFile string
	var duplicatesFile string

	metadataGenerateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a games metadata file from the GOG Api, which can then be applied to a storage",
		Run: func(cmd *cobra.Command, args []string) {
			m, errs, errs404 := sdkPtr.GetMetadata(concurrency, pause, tolerateDangles)
			//duplicates := m.Finalize()
			processSerializableOutput(m, errs, terminalOutput, file)
			
			/*if len(duplicates) > 0 {
				processSerializableOutput(duplicates, []error{}, false, duplicatesFile)
			}*/
			if len(errs404) > 0 {
				errs404Output := Errors{make([]string, len(errs404))}
				for idx, _ := range errs404 {
					errs404Output.Errors[idx] = errs404[idx].Error()
				}
				processSerializableOutput(errs404Output, []error{}, false, warningFile)
			}
		},
	}


	metadataGenerateCmd.Flags().IntVarP(&concurrency, "concurrency", "r", 10, "Maximum number of concurrent requests that will be made on the GOG api")
	metadataGenerateCmd.Flags().IntVarP(&pause, "pause", "s", 200, "Number of milliseconds to wait between batches of api calls")
	metadataGenerateCmd.Flags().StringVarP(&file, "file", "f", "metadata.json", "File to output the metadata in")
	metadataGenerateCmd.Flags().BoolVarP(&terminalOutput, "terminal", "t", false, "If set to true, the metadata will be output on the terminal instead of in a file")
	metadataGenerateCmd.Flags().BoolVarP(&tolerateDangles, "tolerate-dangles", "d", true, "If set to true, undownloadable dangling files (ie, 404 code on download url) will be tolerated and will not prevent metadata generation")
	metadataGenerateCmd.Flags().StringVarP(&warningFile, "warning-file", "w", "metadata-warnings.json", "Warnings from files whose download url return 404 will be listed in this file. Will only be generated if tolerate-dangles is set to true")
	metadataGenerateCmd.Flags().StringVarP(&duplicatesFile, "duplicates-file", "u", "duplicates.json", "Files that had duplicate filenames within the same game and had to be renamed will be listed in this file")
	return metadataGenerateCmd
}

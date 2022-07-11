package cmd

import (
	"gogcli/metadata"

	"github.com/spf13/cobra"
)

func generateMetadataGenerateCmd() *cobra.Command {
	var concurrency int
	var pause int
	var metadataFile string
	var progressFile string
	var terminalOutput bool
	var tolerateDangles bool
	var warningFile string

	metadataGenerateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a games metadata file from the GOG Api, which can then be applied to a storage",
		PreRun: func(cmd *cobra.Command, args []string) {
			CleanupFile(warningFile)
		},
		Run: func(cmd *cobra.Command, args []string) {
			progressFn := PersistMetadataProgress(progressFile)
			writer := metadata.NewMetadataGamesWriter(
				metadata.NewMetadataGamesWriterState([]int64{}),
				logSource,
			)
			errs := writer.Write( 
				sdkPtr.GenerateMetadataGameGetter(concurrency, pause, tolerateDangles),
				progressFn,
			)
			m, warnings := writer.State.Metadata, writer.State.Warnings
			processErrors(errs)

			if len(warnings) > 0 {
				warningsOutput := Errors{make([]string, len(warnings))}
				for idx, _ := range warnings {
					warningsOutput.Errors[idx] = warnings[idx]
				}
				processSerializableOutput(warningsOutput, []error{}, false, warningFile)
			}

			processSerializableOutput(m, []error{}, terminalOutput, metadataFile)
		
			CleanupFile(progressFile)
		},
	}

	metadataGenerateCmd.Flags().IntVarP(&concurrency, "concurrency", "r", 10, "Maximum number of concurrent requests that will be made on the GOG api")
	metadataGenerateCmd.Flags().IntVarP(&pause, "pause", "s", 200, "Number of milliseconds to wait between batches of api calls")
	metadataGenerateCmd.Flags().StringVarP(&metadataFile, "file", "f", "metadata.json", "File to output the metadata in")
	metadataGenerateCmd.Flags().BoolVarP(&terminalOutput, "terminal", "t", false, "If set to true, the metadata will be output on the terminal instead of in a file")
	metadataGenerateCmd.Flags().BoolVarP(&tolerateDangles, "tolerate-dangles", "d", true, "If set to true, undownloadable dangling files (ie, 404 code on download url) will be tolerated and will not prevent metadata generation")
	metadataGenerateCmd.Flags().StringVarP(&warningFile, "warning-file", "w", "metadata-warnings.json", "Warnings from files whose download url return 404 will be listed in this file. Will only be generated if tolerate-dangles is set to true")
	metadataGenerateCmd.Flags().StringVarP(&progressFile, "progress-file", "z", "metadata-generation-progress.json", "File to save transient progress for the metadata generation in")
	return metadataGenerateCmd
}

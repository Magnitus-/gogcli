package cmd

import (
	"encoding/json"
	"fmt"
	"gogcli/metadata"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

func generateMetadataGenerateResumeCmd() *cobra.Command {
	var s *metadata.MetadataGamesWriterState
	var concurrency int
	var pause int
	var metadataFile string
	var progressFile string
	var terminalOutput bool
	var tolerateDangles bool
	var warningFile string

	metadataGenerateResumeCmd := &cobra.Command{
		Use:   "generate-resume",
		Short: "Resume an interrupted metadata file generation",
		PreRun: func(cmd *cobra.Command, args []string) {
			bs, err := ioutil.ReadFile(progressFile)
			if err != nil {
				fmt.Println("Could not load the progress file: ", err)
				os.Exit(1)
			}

			s = &metadata.MetadataGamesWriterState{}
			err = json.Unmarshal(bs, s)
			if err != nil {
				fmt.Println("Progress file doesn't appear to contain valid json: ", err)
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			progressFn := PersistMetadataProgress(progressFile)
			writer := metadata.NewMetadataGamesWriter(
				*s,
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

	metadataGenerateResumeCmd.Flags().IntVarP(&concurrency, "concurrency", "r", 10, "Maximum number of concurrent requests that will be made on the GOG api")
	metadataGenerateResumeCmd.Flags().IntVarP(&pause, "pause", "s", 200, "Number of milliseconds to wait between batches of api calls")
	metadataGenerateResumeCmd.Flags().StringVarP(&metadataFile, "file", "f", "metadata.json", "File to output the metadata in")
	metadataGenerateResumeCmd.Flags().BoolVarP(&terminalOutput, "terminal", "t", false, "If set to true, the metadata will be output on the terminal instead of in a file")
	metadataGenerateResumeCmd.Flags().BoolVarP(&tolerateDangles, "tolerate-dangles", "d", true, "If set to true, undownloadable dangling files (ie, 404 code on download url) will be tolerated and will not prevent metadata generation")
	metadataGenerateResumeCmd.Flags().StringVarP(&warningFile, "warning-file", "w", "metadata-warnings.json", "Warnings from files whose download url return 404 will be listed in this file. Will only be generated if tolerate-dangles is set to true")
	metadataGenerateResumeCmd.Flags().StringVarP(&progressFile, "progress-file", "z", "metadata-generation-progress.json", "File to save transient progress for the metadata generation in")
	return metadataGenerateResumeCmd
}

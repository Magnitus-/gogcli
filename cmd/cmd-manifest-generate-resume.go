package cmd

import (
	"encoding/json"
	"fmt"
	"gogcli/manifest"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

func generateManifestGenerateResumeCmd() *cobra.Command {
	var s *manifest.ManifestGamesWriterState
	var concurrency int
	var pause int
	var manifestFile string
	var progressFile string
	var terminalOutput bool
	var tolerateDangles bool
	var warningFile string
	var duplicatesFile string
	var tolerateBadFileMetadata bool

	manifestGenerateResumeCmd := &cobra.Command{
		Use:   "generate-resume",
		Short: "Resume an interrupted manifest file generation",
		PreRun: func(cmd *cobra.Command, args []string) {
			bs, err := ioutil.ReadFile(progressFile)
			if err != nil {
				fmt.Println("Could not load the progress file: ", err)
				os.Exit(1)
			}

			s = &manifest.ManifestGamesWriterState{}
			err = json.Unmarshal(bs, s)
			if err != nil {
				fmt.Println("Progress file doesn't appear to contain valid json: ", err)
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			progressFn := PersistProgress(progressFile)
			writer := manifest.NewManifestGamesWriter(
				*s,
				logSource,
			)
			errs := writer.Write( 
				sdkPtr.GenerateManifestGameGetter((*s).Manifest.Filter, concurrency, pause, tolerateDangles, tolerateBadFileMetadata),
				progressFn,
			)
			m, warnings := writer.State.Manifest, writer.State.Warnings
			
			if len(warnings) > 0 {
				warningsOutput := Errors{make([]string, len(warnings))}
				for idx, _ := range warnings {
					warningsOutput.Errors[idx] = warnings[idx]
				}
				processSerializableOutput(warningsOutput, []error{}, false, warningFile)
			}
			processErrors(errs)
			
			duplicates := m.Finalize()
			if len(duplicates) > 0 {
				processSerializableOutput(duplicates, []error{}, false, duplicatesFile)
			}

			processSerializableOutput(m, []error{}, terminalOutput, manifestFile)

			CleanupFile(progressFile)
		},
	}
	
	manifestGenerateResumeCmd.Flags().IntVarP(&concurrency, "concurrency", "r", 4, "Maximum number of concurrent requests that will be made on the GOG api")
	manifestGenerateResumeCmd.Flags().IntVarP(&pause, "pause", "s", 200, "Number of milliseconds to wait between batches of api calls")
	manifestGenerateResumeCmd.Flags().StringVarP(&manifestFile, "manifest-file", "f", "manifest.json", "File to output the manifest in")
	manifestGenerateResumeCmd.Flags().StringVarP(&progressFile, "progress-file", "z", "manifest-generation-progress.json", "File to save transient progress for the manifest generation in")
	manifestGenerateResumeCmd.Flags().BoolVarP(&terminalOutput, "terminal", "t", false, "If set to true, the manifest will be output on the terminal instead of in a file")
	manifestGenerateResumeCmd.Flags().BoolVarP(&tolerateDangles, "tolerate-dangles", "d", true, "If set to true, undownloadable dangling files (ie, 404 code on download url) will be tolerated and will not prevent manifest generation")
	manifestGenerateResumeCmd.Flags().StringVarP(&warningFile, "warning-file", "w", "manifest-generation-warnings.json", "Warnings from files whose download url return 404 will be listed in this file. Will only be generated if tolerate-dangles is set to true")
	manifestGenerateResumeCmd.Flags().StringVarP(&duplicatesFile, "duplicates-file", "u", "manifest-generation-duplicates.json", "Files that had duplicate filenames within the same game and had to be renamed will be listed in this file")
	manifestGenerateResumeCmd.Flags().BoolVarP(&tolerateBadFileMetadata, "tolerate-bad-metadata", "b", true, "Tolerate files for which metadata cannot be retrieved. The checksum will be infered by performing a throwaway file download instead.")
	return manifestGenerateResumeCmd
}
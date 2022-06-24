package cmd

import (
	"encoding/json"
	"fmt"
	"gogcli/manifest"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

func generateManifestUpdateResumeCmd() *cobra.Command {
	var previousWarnings []error
	var s *manifest.ManifestGamesWriterState
	var m manifest.Manifest
	var manifestFile string
	var progressFile string
	var warningFile string
	var duplicatesFile string
	var concurrency int
	var pause int
	var tolerateDangles bool
	var tolerateBadFileMetadata bool

	manifestUpdateResumeCmd := &cobra.Command{
		Use:   "update-resume",
		Short: "Resume an interrupted manifest file update",
		PreRun: func(cmd *cobra.Command, args []string) {
			var err error
			m, err = loadManifestFromFile(manifestFile)
			processError(err)

			var bs []byte
			bs, err = ioutil.ReadFile(progressFile)
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

			previousWarnings = deserializeErrors(warningFile)
		},
		Run: func(cmd *cobra.Command, args []string) {
			progressFn := PersistProgress(progressFile)
			writer := manifest.NewManifestGamesWriter(
				*s,
				logSource,
			)
			result := writer.Write( 
				sdkPtr.GenerateManifestGameGetter(m.Filter, concurrency, pause, tolerateDangles, tolerateBadFileMetadata),
				progressFn,
			)
			uManifest, errs, warnings := writer.State.Manifest, result.Errors, result.Warnings

			warnings = append(warnings, previousWarnings...)
			if len(warnings) > 0 {
				warningsOutput := Errors{make([]string, len(warnings))}
				for idx, _ := range warnings {
					warningsOutput.Errors[idx] = warnings[idx].Error()
				}
				processSerializableOutput(warningsOutput, []error{}, false, warningFile)
			}
			processErrors(errs)

			m.OverwriteGames(uManifest.Games)

			duplicates := m.Finalize()
			if len(duplicates) > 0 {
				processSerializableOutput(duplicates, []error{}, false, duplicatesFile)
			}

			processSerializableOutput(m, []error{}, false, manifestFile)
			
			CleanupFile(progressFile)
		},
	}

	manifestUpdateResumeCmd.Flags().IntVarP(&concurrency, "concurrency", "r", 4, "Maximum number of concurrent requests that will be made on the GOG api")
	manifestUpdateResumeCmd.Flags().IntVarP(&pause, "pause", "s", 200, "Number of milliseconds to wait between batches of api calls")
	manifestUpdateResumeCmd.Flags().StringVarP(&manifestFile, "manifest-file", "f", "manifest.json", "Manifest file to update")
	manifestUpdateResumeCmd.MarkFlagFilename("manifest-file")
	manifestUpdateResumeCmd.Flags().StringVarP(&progressFile, "progress-file", "z", "manifest-update-progress.json", "File to resume update from")
	manifestUpdateResumeCmd.Flags().BoolVarP(&tolerateDangles, "tolerate-dangles", "d", true, "If set to true, undownloadable dangling files (ie, 404 code on download url) will be tolerated and will not prevent manifest generation")
	manifestUpdateResumeCmd.Flags().StringVarP(&warningFile, "warning-file", "w", "manifest-update-warnings.json", "Warnings from files whose download url return 404 will be listed in this file. Will only be generated if tolerate-dangles is set to true")
	manifestUpdateResumeCmd.Flags().StringVarP(&duplicatesFile, "duplicates-file", "l", "manifest-update-duplicates.json", "Files that had duplicate filenames within the same game and had to be renamed will be listed in this file")
	manifestUpdateResumeCmd.Flags().BoolVarP(&tolerateBadFileMetadata, "tolerate-bad-metadata", "b", true, "Tolerate files for which metadata cannot be retrieved. The checksum will be infered by performing a throwaway file download instead.")
	return manifestUpdateResumeCmd
}

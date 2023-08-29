package cmd

import (
	"gogcli/manifest"

	"github.com/spf13/cobra"
)

func generateManifestTrimPatchesCmd() *cobra.Command {
	var m manifest.Manifest
	var manifestFile string
	var gameId int64
	var concurrency int
	var pause int
	var tolerateDangles bool
	var tolerateBadFileMetadata bool
	var warningFile string
	var duplicatesFile string

	manifestTrimPatchesCmd := &cobra.Command{
		Use:   "trim-patches",
		Short: "Command to trim patches for a given game in the manifest. Note that game files will also otherwise be updated if dated.",
		PreRun: func(cmd *cobra.Command, args []string) {
			var err error
			m, err = loadManifestFromFile(manifestFile)
			processError(err)
		},
		Run: func(cmd *cobra.Command, args []string) {
			err := m.AddUrlFilterForGame(gameId, "[a-z]{2,2}[0-9]patch[0-9]")
			processError(err)

			writer := manifest.NewManifestGamesWriter(
				manifest.NewManifestGamesWriterState(m.Filter, []int64{gameId}),
				logSource,
			)
			errs := writer.Write( 
				sdkPtr.GenerateManifestGameGetter(m.Filter, concurrency, pause, tolerateDangles, tolerateBadFileMetadata),
				func(state manifest.ManifestGamesWriterState) error {return nil},
			)
			uManifest, warnings := writer.State.Manifest, writer.State.Warnings

			if len(warnings) > 0 {
				warningsOutput := Errors{make([]string, len(warnings))}
				for idx, _ := range warnings {
					warningsOutput.Errors[idx] = warnings[idx]
				}
				processSerializableOutput(warningsOutput, []error{}, false, warningFile)
			}
			processErrors(errs)

			m.OverwriteGames(uManifest.Games)

			duplicates := m.Finalize()
			if len(duplicates) > 0 {
				processSerializableOutput(duplicates, []error{}, false, duplicatesFile)
			}

			processSerializableOutput(m, []error{}, false, manifestFile)		},
	}

	manifestTrimPatchesCmd.Flags().StringVarP(&manifestFile, "manifest-file", "f", "manifest.json", "Manifest file to update")
	manifestTrimPatchesCmd.MarkFlagFilename("manifest-file")
	manifestTrimPatchesCmd.Flags().Int64VarP(&gameId, "id", "i", -1, "File to output the manifest summary in if in json format")
	manifestTrimPatchesCmd.MarkFlagRequired("id")
	manifestTrimPatchesCmd.Flags().StringVarP(&warningFile, "warning-file", "w", "manifest-update-warnings.json", "Warnings from files whose download url return 404 will be listed in this file. Will only be generated if tolerate-dangles is set to true")
	manifestTrimPatchesCmd.Flags().StringVarP(&duplicatesFile, "duplicates-file", "l", "manifest-update-duplicates.json", "Files that had duplicate filenames within the same game and had to be renamed will be listed in this file")
	manifestTrimPatchesCmd.Flags().IntVarP(&concurrency, "concurrency", "r", 4, "Maximum number of concurrent requests that will be made on the GOG api")
	manifestTrimPatchesCmd.Flags().IntVarP(&pause, "pause", "s", 200, "Number of milliseconds to wait between batches of api calls")
	manifestTrimPatchesCmd.Flags().BoolVarP(&tolerateDangles, "tolerate-dangles", "d", true, "If set to true, undownloadable dangling files (ie, 404 code on download url) will be tolerated and will not prevent manifest generation")
	manifestTrimPatchesCmd.Flags().BoolVarP(&tolerateBadFileMetadata, "tolerate-bad-metadata", "b", true, "Tolerate files for which metadata cannot be retrieved. The checksum will be infered by performing a throwaway file download instead.")

	return manifestTrimPatchesCmd
}

package cmd

import (
	"gogcli/manifest"

	"github.com/spf13/cobra"
)

func generateManifestGenerateCmd() *cobra.Command {
	var oses []string
	var languages []string
	var gameTagFilters []string
	var gameTitleFilters []string
	var downloads bool
	var extras bool
	var extraTypeFilters []string
	var skipUrlFilters []string
	var concurrency int
	var pause int
	var manifestFile string
	var progressFile string
	var terminalOutput bool
	var tolerateDangles bool
	var warningFile string
	var duplicatesFile string
	var tolerateBadFileMetadata bool

	manifestGenerateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a games manifest from the GOG Api, which can then be applied to a storage",
		PreRun: func(cmd *cobra.Command, args []string) {
			CleanupFile(warningFile)
			CleanupFile(duplicatesFile)
		},
		Run: func(cmd *cobra.Command, args []string) {
			f := manifest.NewManifestFilter(
				gameTitleFilters,
				oses,
				languages,
				gameTagFilters,
				downloads,
				extras,
				extraTypeFilters,
				skipUrlFilters,
			)
			progressFn := PersistManifestProgress(progressFile)
			writer := manifest.NewManifestGamesWriter(
				manifest.NewManifestGamesWriterState(f, []int64{}),
				logSource,
			)
			errs := writer.Write( 
				sdkPtr.GenerateManifestGameGetter(f, concurrency, pause, tolerateDangles, tolerateBadFileMetadata),
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

	manifestGenerateCmd.Flags().StringArrayVarP(&oses, "os", "o", []string{}, "If you want to include only specific oses. Valid values: windows, mac, linux")
	manifestGenerateCmd.Flags().StringArrayVarP(&languages, "lang", "l", []string{}, "If you want to include only specific languages. Valid values: english, french, spanish, dutch, portuguese_brazilian, russian, korean, chinese_simplified, japanese, polish, italian, german, czech, hungarian, portuguese, danish, finnish, swedish, turkish, arabic, romanian, unknown")
	manifestGenerateCmd.Flags().StringArrayVarP(&gameTagFilters, "tag", "a", []string{}, "If you want to include only games having specific tags")
	manifestGenerateCmd.Flags().StringArrayVarP(&gameTitleFilters, "title", "i", []string{}, "If you want to include only games with title that contain at least one of the given strings")
	manifestGenerateCmd.Flags().BoolVarP(&downloads, "installers", "n", true, "Whether to incluse installer downloads")
	manifestGenerateCmd.Flags().BoolVarP(&extras, "extras", "e", true, "Whether to incluse extras")
	manifestGenerateCmd.Flags().StringArrayVarP(&extraTypeFilters, "extratype", "x", []string{}, "If you want to include only extras whole type contain one of the given strings. Look at full generated manifest without this flag to figure out valid types")
	manifestGenerateCmd.Flags().StringArrayVarP(&skipUrlFilters, "skip-url", "v", []string{}, "Regex of file urls that should be skipped")
	manifestGenerateCmd.Flags().IntVarP(&concurrency, "concurrency", "r", 4, "Maximum number of concurrent requests that will be made on the GOG api")
	manifestGenerateCmd.Flags().IntVarP(&pause, "pause", "s", 200, "Number of milliseconds to wait between batches of api calls")
	manifestGenerateCmd.Flags().StringVarP(&manifestFile, "manifest-file", "f", "manifest.json", "File to output the manifest in")
	manifestGenerateCmd.Flags().StringVarP(&progressFile, "progress-file", "z", "manifest-generation-progress.json", "File to save transient progress for the manifest generation in")
	manifestGenerateCmd.Flags().BoolVarP(&terminalOutput, "terminal", "t", false, "If set to true, the manifest will be output on the terminal instead of in a file")
	manifestGenerateCmd.Flags().BoolVarP(&tolerateDangles, "tolerate-dangles", "d", true, "If set to true, undownloadable dangling files (ie, 404 code on download url) will be tolerated and will not prevent manifest generation")
	manifestGenerateCmd.Flags().StringVarP(&warningFile, "warning-file", "w", "manifest-generation-warnings.json", "Warnings from files whose download url return 404 will be listed in this file. Will only be generated if tolerate-dangles is set to true")
	manifestGenerateCmd.Flags().StringVarP(&duplicatesFile, "duplicates-file", "u", "manifest-generation-duplicates.json", "Files that had duplicate filenames within the same game and had to be renamed will be listed in this file")
	manifestGenerateCmd.Flags().BoolVarP(&tolerateBadFileMetadata, "tolerate-bad-metadata", "b", true, "Tolerate files for which metadata cannot be retrieved. The checksum will be infered by performing a throwaway file download instead.")
	return manifestGenerateCmd
}

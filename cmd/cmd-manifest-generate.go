package cmd

import (
	"github.com/spf13/cobra"
)

func generateManifestGenerateCmd() *cobra.Command {
	var oses []string
	var languages []string
	var gameTagFilters []string
	var gameTitleFilter string
	var downloads bool
	var extras bool
	var extraTypeFilters []string
	var concurrency int
	var pause int
	var file string
	var terminalOutput bool

	manifestGenerateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a games manifest from the GOG Api, which can then be applied to a storage",
		Run: func(cmd *cobra.Command, args []string) {
			m, errs := sdkPtr.GetManifest(gameTitleFilter, concurrency, pause)
			processSerializableOutput(m, errs, terminalOutput, file)
		},
	}

	manifestGenerateCmd.Flags().StringArrayVarP(&oses, "os", "o", []string{}, "If you want to include only specific oses. Valid values: windows, mac, linux")
	manifestGenerateCmd.Flags().StringArrayVarP(&languages, "lang", "l", []string{}, "If you want to include only specific languages. Valid values: english, french, spanish, dutch, portuguese_brazilian, russian, korean, chinese_simplified, japanese, polish, italian, german, czech, hungarian, portuguese, danish, finnish, swedish, turkish, arabic, romanian, unknown")
	manifestGenerateCmd.Flags().StringArrayVarP(&gameTagFilters, "tag", "a", []string{}, "If you want to include only games having specific tags")
	manifestGenerateCmd.Flags().StringVarP(&gameTitleFilter, "title", "i", "", "If you want to include only games with title that contain the given string")
	manifestGenerateCmd.Flags().BoolVarP(&downloads, "installers", "n", true, "Whether to incluse installer downloads")
	manifestGenerateCmd.Flags().BoolVarP(&extras, "extras", "e", true, "Whether to incluse extras")
	manifestGenerateCmd.Flags().StringArrayVarP(&extraTypeFilters, "extratype", "x", []string{}, "If you want to include only extras whole type contain one of the given strings. Look at full generated manifest without this flag to figure out valid types")
	manifestGenerateCmd.Flags().IntVarP(&concurrency, "concurrency", "r", 10, "Maximum number of concurrent requests that will be made on the GOG api")
	manifestGenerateCmd.Flags().IntVarP(&pause, "pause", "s", 200, "Number of milliseconds to wait between batches of api calls")
	manifestGenerateCmd.Flags().StringVarP(&file, "file", "f", "manifest.json", "File to output the manifest in")
	manifestGenerateCmd.Flags().BoolVarP(&terminalOutput, "terminal", "t", false, "If set to true, the manifest will be output on the terminal instead of in a file")
	return manifestGenerateCmd
}

package cmd

import (
	"gogcli/manifest"

	"github.com/spf13/cobra"
)

func generateManifestSearchCmd() *cobra.Command {
	var m manifest.Manifest
	var manifestFile string
	var oses []string
	var languages []string
	var gameTagFilters []string
	var gameTitleFilters []string
	var downloads bool
	var extras bool
	var extraTypeFilters []string
	var file string
	var terminalOutput bool

	manifestSearchCmd := &cobra.Command{
		Use:   "search",
		Short: "Get a subset of a given manifest, corresponding to search terms",
		PreRun: func(cmd *cobra.Command, args []string) {
			var err error
			m, err = loadManifestFromFile(manifestFile)
			processError(err)
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
			)
			f.Intersect(m.Filter)
			m.Filter = f
			m.Trim()
			m.Finalize()
			processSerializableOutput(m, []error{}, terminalOutput, file)
		},
	}

	manifestSearchCmd.Flags().StringVarP(&manifestFile, "manifest", "m", "manifest.json", "Manifest to search")
	manifestSearchCmd.Flags().StringArrayVarP(&oses, "os", "o", []string{}, "If you want to include only specific oses. Valid values: windows, mac, linux")
	manifestSearchCmd.Flags().StringArrayVarP(&languages, "lang", "l", []string{}, "If you want to include only specific languages. Valid values: english, french, spanish, dutch, portuguese_brazilian, russian, korean, chinese_simplified, japanese, polish, italian, german, czech, hungarian, portuguese, danish, finnish, swedish, turkish, arabic, romanian, unknown")
	manifestSearchCmd.Flags().StringArrayVarP(&gameTagFilters, "tag", "a", []string{}, "If you want to include only games having specific tags")
	manifestSearchCmd.Flags().StringArrayVarP(&gameTitleFilters, "title", "i", []string{}, "If you want to include only games with title that contain at least one of the given strings")
	manifestSearchCmd.Flags().BoolVarP(&downloads, "installers", "n", true, "Whether to incluse installer downloads")
	manifestSearchCmd.Flags().BoolVarP(&extras, "extras", "e", true, "Whether to incluse extras")
	manifestSearchCmd.Flags().StringArrayVarP(&extraTypeFilters, "extratype", "x", []string{}, "If you want to include only extras whole type contain one of the given strings. Look at full generated manifest without this flag to figure out valid types")
	manifestSearchCmd.Flags().StringVarP(&file, "file", "f", "search.json", "File to output the search in")
	manifestSearchCmd.Flags().BoolVarP(&terminalOutput, "terminal", "t", true, "If set to true, the search will be output on the terminal instead of in a file")
	return manifestSearchCmd
}

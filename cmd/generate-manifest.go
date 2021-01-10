package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func generateManifestGenerationCmd() *cobra.Command {
	var oses []string
	var languages []string
	var gameTagFilters []string
	var gameTitleFilter string
	var downloads bool
	var extras bool
	var extraTypeFilters []string

	manifestGenerationCmd := &cobra.Command{
		Use:   "generate-manifest",
		Short: "Generate a download manifest from the GOG Api",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("currently a noop")
		},
	}
	manifestGenerationCmd.Flags().StringArrayVarP(&oses, "os", "o", []string{}, "If you want to include only specific oses. Valid values: windows, mac, linux")
	manifestGenerationCmd.Flags().StringArrayVarP(&languages, "lang", "l", []string{}, "If you want to include only specific languages. Valid values: english, french, spanish, spanish_al, polish, dutch, italian, portuguese, portuguese_brazilian, turkish, romanian")
	manifestGenerationCmd.Flags().StringArrayVarP(&gameTagFilters, "tag", "a", []string{}, "If you want to include only games having specific tags")
	manifestGenerationCmd.Flags().StringVarP(&gameTitleFilter, "title", "i", "", "If you want to include only games with title that contain the given string")
	manifestGenerationCmd.Flags().BoolVarP(&downloads, "installers", "n", true, "Whether to incluse installer downloads")
	manifestGenerationCmd.Flags().BoolVarP(&extras, "extras", "e", true, "Whether to incluse extras")
	manifestGenerationCmd.Flags().StringArrayVarP(&extraTypeFilters, "extratype", "x", []string{}, "If you want to include only extras whole type contain one of the given strings. Look at full generated manifest without this flag to figure out valid types")
	return manifestGenerationCmd
}

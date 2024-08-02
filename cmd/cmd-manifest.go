package cmd

import (
	"github.com/spf13/cobra"
)

func generateManifestCmd() *cobra.Command {
	manifestCmd := &cobra.Command{
		Use:   "manifest",
		Short: "Commands to generate, manipulate and get info from a games manifest",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			callPersistentPreRun(cmd, args)
		},
	}

	manifestCmd.AddCommand(generateManifestGenerateCmd())
	manifestCmd.AddCommand(generateManifestGenerateResumeCmd())
	manifestCmd.AddCommand(generateManifestSummaryCmd())
	manifestCmd.AddCommand(generateManifestSearchCmd())
	manifestCmd.AddCommand(generateManifestUpdateCmd())
	manifestCmd.AddCommand(generateManifestUpdateResumeCmd())
	manifestCmd.AddCommand(generateManifestDiffCmd())
	manifestCmd.AddCommand(generateManifestMigrateCmd())
	manifestCmd.AddCommand(generateManifestTrimLanguagesCmd())
	manifestCmd.AddCommand(generateManifestTrimPatchesCmd())
	manifestCmd.AddCommand(generateManifestMissingGamesCmd())

	return manifestCmd
}

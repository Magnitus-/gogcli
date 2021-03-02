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
	manifestCmd.AddCommand(generateManifestSummaryCmd())

	return manifestCmd
}
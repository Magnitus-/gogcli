package cmd

import (
	"github.com/spf13/cobra"
)

func generateMetadataCmd() *cobra.Command {
	metadataCmd := &cobra.Command{
		Use:   "metadata",
		Short: "Commands to generate, manipulate and get info from games metadata",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			callPersistentPreRun(cmd, args)
		},
	}

	metadataCmd.AddCommand(generateMetadataGenerateCmd())
	metadataCmd.AddCommand(generateMetadataGenerateResumeCmd())
	metadataCmd.AddCommand(generateMetadataUpdateCmd())

	return metadataCmd
}

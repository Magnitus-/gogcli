package cmd

import (
	"github.com/spf13/cobra"
)

func generateUpdateCmd() *cobra.Command {
	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Commands to manage update files based of what is new or got updated in GOG.com",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			callPersistentPreRun(cmd, args)
		},
	}

	updateCmd.AddCommand(generateUpdateGenerateCmd())

	return updateCmd
}
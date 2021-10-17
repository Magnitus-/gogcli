package cmd

import (
	"github.com/spf13/cobra"
)

func generateActionsCmd() *cobra.Command {
	actionsCmd := &cobra.Command{
		Use:   "actions",
		Short: "Command to get Information on an actions file",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			callPersistentPreRun(cmd, args)
		},
	}

	actionsCmd.AddCommand(generateActionsSummaryCmd())

	return actionsCmd
}

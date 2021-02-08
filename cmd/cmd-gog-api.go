package cmd

import (
	"github.com/spf13/cobra"
)

func generateGogApiCmd() *cobra.Command {
	gogApiCmd := &cobra.Command{
		Use:   "gog-api",
		Short: "Command to interact with the gog api. Can be used to troubleshoot the sdk or build other tools on top of this client.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			callPersistentPreRun(cmd, args)
		},
	}

	gogApiCmd.AddCommand(generateOwnedGamesCmd())
	gogApiCmd.AddCommand(generateGameDetailsCmd())
	gogApiCmd.AddCommand(generateUserInfoCmd())
	gogApiCmd.AddCommand(generateDownloadUrlPathCmd())
	gogApiCmd.AddCommand(generateUrlPathFilenameCmd())

	return gogApiCmd
}

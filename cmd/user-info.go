package cmd

import (
	"github.com/spf13/cobra"
)

func generateUserInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "user-info",
		Short: "Command to retrieve your GOG user summary",
		Run: func(cmd *cobra.Command, args []string) {
			sdkInst.GetUser(debugMode).Print()
		},
	}
}

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func generateUserInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "user-info",
		Short: "Command to retrieve your GOG user summary",
		Run: func(cmd *cobra.Command, args []string) {
			user, err := sdkPtr.GetUser(debugMode)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			user.Print()
		},
	}
}

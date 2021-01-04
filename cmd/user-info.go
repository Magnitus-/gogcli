package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var userInfoCmd = &cobra.Command{
	Use:   "userinfo",
	Short: "Command to retrieve your GOG user summary",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(sdkInst.GetUser())
	},
}

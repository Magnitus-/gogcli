package cmd

import (
	"fmt"
	"gogcli/login"

	"github.com/spf13/cobra"
)

func generateLoginCmd() *cobra.Command {
	loginCmd := &cobra.Command{
		Use:   "login",
		Short: "Login to gog and generate a cookie file",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			callPersistentPreRun(cmd, args)
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Type the following url in your browser: http://localhost:8080")
			login.Serve()
		},
	}

	return loginCmd
}

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func generateVersionCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Outputs the current version of the tool",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("0.9.0")
		},
	}

	return versionCmd
}
package cmd

import (
	"github.com/spf13/cobra"
)

func generateUpdateGenerateCmd() *cobra.Command {
	var concurrency int
	var pause int
	var file string
	var terminalOutput bool

	updateGenerateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate update file based on what is new or got updated in GOG.com",
		Run: func(cmd *cobra.Command, args []string) {
			u, errs := sdkPtr.GetUpdates(concurrency, pause)
			processSerializableOutput(u, errs, terminalOutput, file)
		},
	}

	updateGenerateCmd.Flags().IntVarP(&concurrency, "concurrency", "r", 10, "Maximum number of concurrent requests that will be made on the GOG api")
	updateGenerateCmd.Flags().IntVarP(&pause, "pause", "s", 200, "Number of milliseconds to wait between batches of api calls")
	updateGenerateCmd.Flags().StringVarP(&file, "file", "f", "updates.json", "File to output the updates in")
	updateGenerateCmd.Flags().BoolVarP(&terminalOutput, "terminal", "t", false, "If set to true, the manifest will be output on the terminal instead of in a file")

	return updateGenerateCmd
}

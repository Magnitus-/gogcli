package cmd

import (
	"gogcli/manifest"

	"github.com/spf13/cobra"
)

func generateActionsSummaryCmd() *cobra.Command {
	var a manifest.GameActions
	var actionsPath string
	var file string
	var terminalOutput bool

	actionsSummaryCmd := &cobra.Command{
		Use:   "summary",
		Short: "Command to retrieve the summary of an action file",
		PreRun: func(cmd *cobra.Command, args []string) {
			var err error
			a, err = loadActionsFromFile(actionsPath)
			processError(err)
		},
		Run: func(cmd *cobra.Command, args []string) {
			summary := a.GetSummary()
			processSerializableOutput(summary, []error{}, terminalOutput, file)
		},
	}

	actionsSummaryCmd.Flags().StringVarP(&actionsPath, "actions", "a", "actions.json", "Actions file to get the summary about")
	actionsSummaryCmd.Flags().StringVarP(&file, "file", "f", "actions-info.json", "File to output the actions summary in if in json format")
	actionsSummaryCmd.Flags().BoolVarP(&terminalOutput, "terminal", "t", true, "If set to true and json format is used, the actions summary will be output on the terminal instead of in a file")

	return actionsSummaryCmd
}

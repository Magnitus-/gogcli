package cmd

import (
	"gogcli/manifest"

	"github.com/spf13/cobra"
)

func generateManifestSummaryCmd() *cobra.Command {
	var m manifest.Manifest
	var manifestPath string
	var summaryFile string
	var terminalOutput bool

	manifestSummaryCmd := &cobra.Command{
		Use:   "summary",
		Short: "Command to retrieve the summary of a manifest",
		PreRun: func(cmd *cobra.Command, args []string) {
			var err error
			m, err = loadManifestFromFile(manifestPath)
			processError(err)
		},
		Run: func(cmd *cobra.Command, args []string) {
			summary := m.GetSummary()
			processSerializableOutput(summary, []error{}, terminalOutput, summaryFile)
		},
	}

	manifestSummaryCmd.Flags().StringVarP(&manifestPath, "manifest", "m", "manifest.json", "Manifest file to get the summary about")
	manifestSummaryCmd.Flags().StringVarP(&summaryFile, "summary-file", "f", "manifest-info.json", "File to output the manifest summary in if in json format")
	manifestSummaryCmd.Flags().BoolVarP(&terminalOutput, "terminal", "t", true, "If set to true and json format is used, the manifest summary will be output on the terminal instead of in a file")

	return manifestSummaryCmd
}

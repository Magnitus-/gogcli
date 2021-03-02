package cmd

import (
	"encoding/json"
	"fmt"
	"gogcli/manifest"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

func generateManifestSummaryCmd() *cobra.Command {
	var m manifest.Manifest
	var manifestPath string
	var file string
	var terminalOutput bool

	manifestSummaryCmd := &cobra.Command{
		Use:   "summary",
		Short: "Command to retrieve the summary of a manifest",
		PreRun: func(cmd *cobra.Command, args []string) {
			bs, err := ioutil.ReadFile(manifestPath)
			if err != nil {
				fmt.Println("Could not load the manifest: ", err)
				os.Exit(1)
			}

			err = json.Unmarshal(bs, &m)
			if err != nil {
				fmt.Println("Manifest file doesn't appear to contain valid json: ", err)
				os.Exit(1)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			summary := m.GetSummary()
			processSerializableOutput(summary, []error{}, terminalOutput, file)
		},
	}

	manifestSummaryCmd.Flags().StringVarP(&manifestPath, "manifest", "m", "manifest.json", "Manifest file to get the summary about")
	manifestSummaryCmd.Flags().StringVarP(&file, "file", "f", "manifest-info.json", "File to output the manifest summary in if in json format")
	manifestSummaryCmd.Flags().BoolVarP(&terminalOutput, "terminal", "t", true, "If set to true and json format is used, the manifest summary will be output on the terminal instead of in a file")

	return manifestSummaryCmd
}

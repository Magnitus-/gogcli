package cmd

import (
	"gogcli/gameupdates"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"github.com/spf13/cobra"
)

func generateUpdateAccumulateCmd() *cobra.Command {
	var prev *gameupdates.Updates
	var concurrency int
	var pause int
	var updateFile string
	var terminalOutput bool

	updateAccumulateCmd := &cobra.Command{
		Use:   "accumulate",
		Short: "Add to the content of an existing update file based on what is new or got updated in GOG.com",
		PreRun: func(cmd *cobra.Command, args []string) {
			bs, err := ioutil.ReadFile(updateFile)
			if err != nil {
				fmt.Println("Could not load the update file: ", err)
				os.Exit(1)
			}

			prev = &gameupdates.Updates{}
			err = json.Unmarshal(bs, prev)
			if err != nil {
				fmt.Println("Update file doesn't appear to contain valid json: ", err)
				os.Exit(1)
			}
		},
        Run: func(cmd *cobra.Command, args []string) {
			u, errs := sdkPtr.GetUpdates(concurrency, pause)
			processErrors(errs)
            u.Merge(prev)
			processSerializableOutput(u, []error{}, terminalOutput, updateFile)
		},
	}

	updateAccumulateCmd.Flags().IntVarP(&concurrency, "concurrency", "r", 4, "Maximum number of concurrent requests that will be made on the GOG api")
	updateAccumulateCmd.Flags().IntVarP(&pause, "pause", "s", 200, "Number of milliseconds to wait between batches of api calls")
	updateAccumulateCmd.Flags().StringVarP(&updateFile, "update-file", "f", "updates.json", "File to add updates to")
	updateAccumulateCmd.Flags().BoolVarP(&terminalOutput, "terminal", "t", false, "If set to true, the updates will be output on the terminal instead of the file")

	return updateAccumulateCmd
}

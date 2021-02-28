package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	
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
			hasErr := false
			var buf bytes.Buffer
			var output []byte
			var e Errors
			u, errs := sdkPtr.GetUpdates(concurrency, pause)

			if len(errs) > 0 {
				for _, err := range errs {
					e.Errors = append(e.Errors, err.Error())
				}
				output, _ = json.Marshal(e)
				hasErr = true
			} else {
				output, _ = json.Marshal(u)
			}

			json.Indent(&buf, output, "", "  ")
			output = buf.Bytes()

			if terminalOutput {
				fmt.Println(string(output))
			} else {
				err := ioutil.WriteFile(file, output, 0644)
				if err != nil {
					fmt.Println(err)
					hasErr = true
				}
			}

			if hasErr {
				os.Exit(1)
			}
		},
	}

	updateGenerateCmd.Flags().IntVarP(&concurrency, "concurrency", "r", 10, "Maximum number of concurrent requests that will be made on the GOG api")
	updateGenerateCmd.Flags().IntVarP(&pause, "pause", "s", 200, "Number of milliseconds to wait between batches of api calls")
	updateGenerateCmd.Flags().StringVarP(&file, "file", "f", "updates.json", "File to output the manifest in")
	updateGenerateCmd.Flags().BoolVarP(&terminalOutput, "terminal", "t", false, "If set to true, the manifest will be output on the terminal instead of in a file")

	return updateGenerateCmd
}
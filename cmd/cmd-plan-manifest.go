package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gogcli/manifest"
	"gogcli/storage"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

type showActions func(*manifest.GameActions)

func generatePlanManifestFsCmd(m *manifest.Manifest, fn showActions) *cobra.Command {
	var path string

	planManifestFsCmd := &cobra.Command{
		Use:   "fs",
		Short: "Plan the manifest on a file system store",
		PreRun: func(cmd *cobra.Command, args []string) {
			_, err := os.Stat(path)
			if err != nil {
				if os.IsNotExist(err) {
					err = os.MkdirAll(path, 0755)
					if err != nil {
						fmt.Println("Failed to create a directory at the specified path: ", err)
						os.Exit(1)
					}
				} else {
					fmt.Println("An error occured while trying to gather info on the specified path: ", err)
					os.Exit(1)
				}
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			actionsPtr, err := storage.PlanManifest(m, storage.GetFileSystem(path, debugMode))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fn(actionsPtr)
		},
	}

	planManifestFsCmd.Flags().StringVarP(&path, "path", "p", "games", "Path to the directory where game files should be stored")
	return planManifestFsCmd
}

func generatePlanManifestCmd() *cobra.Command {
	var m manifest.Manifest
	var manifestPath string
	var file string
	var terminalOutput bool

	show := func(a *manifest.GameActions) {
		var buf bytes.Buffer
		var output []byte
		
		output, _ = json.Marshal((*a))
		json.Indent(&buf, output, "", "  ")
		output = buf.Bytes()

		if terminalOutput {
			fmt.Println(string(output))
		} else {
			err := ioutil.WriteFile(file, output, 0644)
			if err != nil {
				fmt.Println("Could not write output to file:", err)
				os.Exit(1)
			}
		}
	}

	planManifestCmd := &cobra.Command{
		Use:   "plan-manifest",
		Short: "Generate a plan of the actions that would be executed if the manifest was applied",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
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
	}

	planManifestCmd.PersistentFlags().StringVarP(&manifestPath, "manifest", "m", "manifest.json", "Path were the manifest you want to apply is")
	planManifestCmd.MarkPersistentFlagFilename("manifest")
	planManifestCmd.PersistentFlags().StringVarP(&file, "file", "f", "actions.json", "File to output the plan in")
	planManifestCmd.PersistentFlags().BoolVarP(&terminalOutput, "terminal", "t", false, "If set to true, the plan will be output on the terminal instead of in a file")

	planManifestCmd.AddCommand(generatePlanManifestFsCmd(&m, show))

	return planManifestCmd
}

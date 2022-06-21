package cmd

import (
	"gogcli/manifest"

	"github.com/spf13/cobra"
)

func generateManifestDiffCmd() *cobra.Command {
	var curr manifest.Manifest
	var next manifest.Manifest
	var currManifestFile string
	var nextManifestFile string
	var diffFile string
	var terminalOutput bool
	var allowEmptyCheckum bool

	manifestDiffCmd := &cobra.Command{
		Use:   "diff",
		Short: "Get a plan between two manifests",
		PreRun: func(cmd *cobra.Command, args []string) {
			var err error
			curr, err = loadManifestFromFile(currManifestFile)
			processError(err)
			next, err = loadManifestFromFile(nextManifestFile)
			processError(err)
		},
		Run: func(cmd *cobra.Command, args []string) {
			checksumValidation := manifest.ChecksumValidation
			if allowEmptyCheckum {
				checksumValidation = manifest.ChecksumValidationIfPresent
			}

			a := curr.Plan(&next, checksumValidation, false)
			processSerializableOutput(a, []error{}, terminalOutput, diffFile)
		},
	}

	manifestDiffCmd.Flags().StringVarP(&currManifestFile, "current", "u", "manifest.json", "Current manifest that represents the state of your storage")
	manifestDiffCmd.MarkFlagFilename("current")
	manifestDiffCmd.Flags().StringVarP(&nextManifestFile, "next", "n", "next-manifest.json", "Next manifest that represents the desired state of your storage")
	manifestDiffCmd.MarkFlagFilename("next")
	manifestDiffCmd.Flags().StringVarP(&diffFile, "diff-file", "f", "diff-actions.json", "File to output the actions representing the difference")
	manifestDiffCmd.Flags().BoolVarP(&terminalOutput, "terminal", "t", true, "If set to true, the actions will be output on the terminal instead of in a file")
	manifestDiffCmd.Flags().BoolVarP(&allowEmptyCheckum, "empty-checksum", "s", false, "If set to true, files in the desired manifest with empty checksums will count as already uploaded if everything else matches")

	return manifestDiffCmd
}

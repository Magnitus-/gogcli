package cmd

import (
	"encoding/json"
	"gogcli/migration"
	"io/ioutil"

	"github.com/spf13/cobra"
)

func generateManifestMigrateCmd() *cobra.Command {
	var m migration.ManifestV0_9
	var manifestFile string

	manifestMigrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate a manifest file from version v0.9.x to the current format",
		PreRun: func(cmd *cobra.Command, args []string) {
			bs, err := ioutil.ReadFile(manifestFile)
			processError(err)

			err = json.Unmarshal(bs, &m)
			processError(err)
		},
		Run: func(cmd *cobra.Command, args []string) {
			next := m.Migrate()
			processSerializableOutput(next, []error{}, false, manifestFile)
		},
	}

	manifestMigrateCmd.Flags().StringVarP(&manifestFile, "manifest", "m", "manifest.json", "Manifest file to migrate")
	manifestMigrateCmd.MarkFlagFilename("manifest")
	return manifestMigrateCmd
}

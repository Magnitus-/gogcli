package cmd

import (
	"github.com/spf13/cobra"
)

func generateStorageApplyCmd() *cobra.Command {
	storageApplyCmd := &cobra.Command{
		Use:   "apply",
		Short: "Commands to apply state files (currently, manifests only) to the storage",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			callPersistentPreRun(cmd, args)
		},
	}

	storageApplyCmd.AddCommand(generateStorageApplyManifestCmd())

	return storageApplyCmd
}
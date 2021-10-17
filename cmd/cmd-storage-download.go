package cmd

import (
	"github.com/spf13/cobra"
)

func generateStorageDownloadCmd() *cobra.Command {
	storageDownloadCmd := &cobra.Command{
		Use:   "download",
		Short: "Commands to download things from the storage",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			callPersistentPreRun(cmd, args)
		},
	}

	storageDownloadCmd.AddCommand(generateStorageDownloadManifestCmd())
	storageDownloadCmd.AddCommand(generateStorageDownloadActionsCmd())

	return storageDownloadCmd
}

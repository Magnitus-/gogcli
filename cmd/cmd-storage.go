package cmd

import (
	"github.com/spf13/cobra"
)

func generateStorageCmd() *cobra.Command {
	storageCmd := &cobra.Command{
		Use:   "storage",
		Short: "Commands to upload to, download from, copy and verify storage",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			callPersistentPreRun(cmd, args)
		},
	}

	storageCmd.AddCommand(generateStoragePlanCmd())
	storageCmd.AddCommand(generateStorageApplyCmd())
	storageCmd.AddCommand(generateStorageCopyCmd())
	storageCmd.AddCommand(generateStorageValidateCmd())
	storageCmd.AddCommand(generateStorageResumeCmd())
	storageCmd.AddCommand(generateStorageDownloadCmd())
	storageCmd.AddCommand(generateStorageUpdateActionsCmd())
	storageCmd.AddCommand(generateStorageRepairCmd())

	return storageCmd
}
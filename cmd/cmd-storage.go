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
	storageCmd.AddCommand(generateStorageExecuteActionsCmd())
	storageCmd.AddCommand(generateStorageDownloadCmd())
	storageCmd.AddCommand(generateStorageRepairCmd())
	storageCmd.AddCommand(generateStorageRepairResumeCmd())
	storageCmd.AddCommand(generateStorageAddFileProtectionCmd())
	storageCmd.AddCommand(generateStorageRemoveFileProtectionCmd())

	return storageCmd
}

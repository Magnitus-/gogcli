package cmd

import (
	"errors"
	"gogcli/manifest"

	"github.com/spf13/cobra"
)

func generateStorageRemoveFileProtectionCmd() *cobra.Command {
	var gameId int64
	var FileType string
	var FileName string
	var path string
	var storageType string

	storageRemoveFileProtectionCmd := &cobra.Command{
		Use:   "remove-file-protection",
		Short: "Remove protection a file in your manifest has against deletion",
		Run: func(cmd *cobra.Command, args []string) {
			gamesStorage, _ := getStorage(path, storageType, logSource, "")

			exists, err := gamesStorage.Exists()
			processError(err)
			if !exists {
				processError(errors.New("Specified storage doesn't exist"))
			}

			has, hasErr := gamesStorage.HasManifest()
			processError(hasErr)
			if !has {
				processError(errors.New("Specified storage doesn't have a manifest"))
			}

			m, mErr := gamesStorage.LoadManifest()
			processError(mErr)

			m.ProtectedFiles.RemoveGameFile(manifest.FileInfo{
				Game: manifest.GameInfo{
					Id: gameId,
				},
				Kind: FileType,
				Name: FileName,
			})

			sErr := gamesStorage.StoreManifest(m)
			processError(sErr)
		},
	}

	storageRemoveFileProtectionCmd.Flags().StringVarP(&path, "path", "p", "games", "Path to the directory where game files should be stored")
	storageRemoveFileProtectionCmd.Flags().StringVarP(&storageType, "storage", "k", "fs", "The type of storage you are using. Can be 'fs' (for file system) or 's3' (for s3 store)")
	storageRemoveFileProtectionCmd.Flags().StringVarP(&FileType, "file-type", "t", "installer", "Type of the file to remove protection from. Can be 'installer' or 'extra'")
	storageRemoveFileProtectionCmd.Flags().StringVarP(&FileName, "file-name", "n", "", "Name of the file to remove protection from")
	storageRemoveFileProtectionCmd.MarkFlagRequired("file-name")
	storageRemoveFileProtectionCmd.Flags().Int64VarP(&gameId, "game-id", "i", -1, "Id of the game that contains the file you want to remove protection from")
	storageRemoveFileProtectionCmd.MarkFlagRequired("game-id")
	return storageRemoveFileProtectionCmd
}
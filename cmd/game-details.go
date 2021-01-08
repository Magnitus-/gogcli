package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func generateGameDetailsCmd() *cobra.Command {
	var gameId int

	gameDetailsCmd := &cobra.Command{
		Use:   "game-details",
		Short: "Retrieve details about a given game including link to download files",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(sdkInst.GetGameDetails(gameId))
		},
	}

	gameDetailsCmd.Flags().IntVarP(&gameId, "id", "i", 0, "Id of the game to get details from")
	gameDetailsCmd.MarkFlagRequired("id")

	return gameDetailsCmd
}

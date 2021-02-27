package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func generateGameDetailsCmd() *cobra.Command {
	var gameId int64

	gameDetailsCmd := &cobra.Command{
		Use:   "game-details",
		Short: "Retrieve details about a given game including link to download files",
		Run: func(cmd *cobra.Command, args []string) {
			g, err := sdkPtr.GetGameDetails(gameId)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			g.Print()
		},
	}

	gameDetailsCmd.Flags().Int64VarP(&gameId, "id", "i", 0, "Id of the game to get details from")
	gameDetailsCmd.MarkFlagRequired("id")

	return gameDetailsCmd
}

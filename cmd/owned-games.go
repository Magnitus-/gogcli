package cmd

import (
	"github.com/spf13/cobra"
)

func generateOwnedGamesCmd() *cobra.Command {
	var page int

	ownedGamesCmd := &cobra.Command{
		Use:   "owned-games",
		Short: "Command to retrieve a list of games you own",
		Run: func(cmd *cobra.Command, args []string) {
			sdkInst.GetOwnedGames(page).Print()
		},
	}

	ownedGamesCmd.Flags().IntVarP(&page, "page", "p", 1, "Page number to look at from the results pages")

	return ownedGamesCmd
}

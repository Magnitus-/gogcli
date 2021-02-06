package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func generateOwnedGamesCmd() *cobra.Command {
	var page int
	var search string

	ownedGamesCmd := &cobra.Command{
		Use:   "owned-games",
		Short: "Command to retrieve a list of games you own",
		Run: func(cmd *cobra.Command, args []string) {
			o, err := sdkPtr.GetOwnedGames(page, search)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			o.Print()
		},
	}

	ownedGamesCmd.Flags().IntVarP(&page, "page", "p", 1, "Page to fetch if the result spans multiple pages")
	ownedGamesCmd.Flags().StringVarP(&search, "search", "s", "", "Return only games whose title contain the term")

	return ownedGamesCmd
}

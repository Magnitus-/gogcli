package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var getOwnedGamesCmd = &cobra.Command{
	Use:   "get-owned-games",
	Short: "Command to retrieve a list of games you own",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%+v\n", sdkInst.GetOwnedGames(1))
	},
}

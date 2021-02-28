package cmd

import (
	"github.com/spf13/cobra"
)

func generateOwnedGamesCmd() *cobra.Command {
	var page int
	var search string
	var file string
	var terminalOutput bool
	var jsonOutput bool

	ownedGamesCmd := &cobra.Command{
		Use:   "owned-games",
		Short: "Command to retrieve a list of games you own",
		Run: func(cmd *cobra.Command, args []string) {
			o, err := sdkPtr.GetOwnedGames(page, search)
			if !jsonOutput {
				processError(err) 
				o.Print()
			} else {
				errs := make([]error, 0)
				if err != nil {
					errs = append(errs, err)
				}
				processSerializableOutput(o, errs, terminalOutput, file)
			}
		},
	}

	ownedGamesCmd.Flags().IntVarP(&page, "page", "p", 1, "Page to fetch if the result spans multiple pages")
	ownedGamesCmd.Flags().StringVarP(&search, "search", "s", "", "Return only games whose title contain the term")
	ownedGamesCmd.Flags().StringVarP(&file, "file", "f", "owned-games.json", "File to output the owned games information in if in json format")
	ownedGamesCmd.Flags().BoolVarP(&terminalOutput, "terminal", "t", true, "If set to true and json format is used, the owned games information will be output on the terminal instead of in a file")
	ownedGamesCmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "If set to true, the output will be in json format either on the terminal or in a file. Otherwise, it will be in human readable format on the terminal.")

	return ownedGamesCmd
}

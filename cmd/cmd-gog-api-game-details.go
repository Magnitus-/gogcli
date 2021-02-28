package cmd

import (
	"github.com/spf13/cobra"
)

func generateGameDetailsCmd() *cobra.Command {
	var gameId int64
	var file string
	var terminalOutput bool
	var jsonOutput bool

	gameDetailsCmd := &cobra.Command{
		Use:   "game-details",
		Short: "Retrieve details about a given game including link to download files",
		Run: func(cmd *cobra.Command, args []string) {
			g, err := sdkPtr.GetGameDetails(gameId)
			if !jsonOutput {
				processError(err) 
				g.Print()
			} else {
				errs := make([]error, 0)
				if err != nil {
					errs = append(errs, err)
				}
				processSerializableOutput(g, errs, terminalOutput, file)
			}
		},
	}

	gameDetailsCmd.Flags().Int64VarP(&gameId, "id", "i", 0, "Id of the game to get details from")
	gameDetailsCmd.MarkFlagRequired("id")
	gameDetailsCmd.Flags().StringVarP(&file, "file", "f", "game-details.json", "File to output the game details information in if in json format")
	gameDetailsCmd.Flags().BoolVarP(&terminalOutput, "terminal", "t", true, "If set to true and json format is used, the game details information will be output on the terminal instead of in a file")
	gameDetailsCmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "If set to true, the output will be in json format either on the terminal or in a file. Otherwise, it will be in human readable format on the terminal.")

	return gameDetailsCmd
}

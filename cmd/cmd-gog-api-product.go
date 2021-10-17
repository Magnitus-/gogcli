package cmd

import (
	"github.com/spf13/cobra"
)

func generateProductCmd() *cobra.Command {
	var gameId int64
	var file string
	var terminalOutput bool

	productCmd := &cobra.Command{
		Use:   "product",
		Short: "Command to retrieve product information from a game you own",
		Run: func(cmd *cobra.Command, args []string) {
			o, _, err := sdkPtr.GetProduct(gameId)
			errs := make([]error, 0)
			if err != nil {
				errs = append(errs, err)
			}
			processSerializableOutput(o, errs, terminalOutput, file)
		},
	}

	productCmd.Flags().Int64VarP(&gameId, "id", "i", 0, "Id of the game to get product info from")
	productCmd.MarkFlagRequired("id")
	productCmd.Flags().StringVarP(&file, "file", "f", "product.json", "File to output the product information in if in json format")
	productCmd.Flags().BoolVarP(&terminalOutput, "terminal", "t", true, "If set to true and json format is used, the product information will be output on the terminal instead of in a file")

	return productCmd
}

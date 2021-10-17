package cmd

import (
	"github.com/spf13/cobra"
)

func generateUserInfoCmd() *cobra.Command {
	var file string
	var terminalOutput bool
	var jsonOutput bool

	userInfoCmd := &cobra.Command{
		Use:   "user-info",
		Short: "Command to retrieve your GOG user summary",
		Run: func(cmd *cobra.Command, args []string) {
			user, err := sdkPtr.GetUser()
			if !jsonOutput {
				processError(err)
				user.Print()
			} else {
				errs := make([]error, 0)
				if err != nil {
					errs = append(errs, err)
				}
				processSerializableOutput(user, errs, terminalOutput, file)
			}
		},
	}

	userInfoCmd.Flags().StringVarP(&file, "file", "f", "user.json", "File to output the user information in if in json format")
	userInfoCmd.Flags().BoolVarP(&terminalOutput, "terminal", "t", true, "If set to true and json format is used, the user information will be output on the terminal instead of in a file")
	userInfoCmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "If set to true, the output will be in json format either on the terminal or in a file. Otherwise, it will be in human readable format on the terminal.")

	return userInfoCmd
}

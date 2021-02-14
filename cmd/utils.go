package cmd

import (
    "fmt"
    "os"

	"github.com/spf13/cobra"
)

func processErrors(errs []error) {
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
		os.Exit(1)
	}	
}

func processError(err error) {
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

//https://github.com/spf13/cobra/issues/216#issuecomment-703846787
func callPersistentPreRun(cmd *cobra.Command, args []string) { 
	parent := cmd.Parent()
	if parent != nil { 
        if parent.PersistentPreRun != nil { 
            parent.PersistentPreRun(parent, args) 
        } 
    } 
} 
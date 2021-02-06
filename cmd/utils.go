package cmd

import (
	"github.com/spf13/cobra"
)

//https://github.com/spf13/cobra/issues/216#issuecomment-703846787
func callPersistentPreRun(cmd *cobra.Command, args []string) { 
	parent := cmd.Parent()
	if parent != nil { 
        if parent.PersistentPreRun != nil { 
            parent.PersistentPreRun(parent, args) 
        } 
    } 
} 
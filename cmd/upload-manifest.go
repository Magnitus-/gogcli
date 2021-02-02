package cmd

import (
	"fmt"
	"gogcli/manifest"
	"gogcli/storage"
	"os"
)

func uploadManifest(m *manifest.Manifest, s storage.Storage, concurrency int, pause int) {
	actionsPtr, err := storage.PlanManifest(m, s)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len((*actionsPtr)) == 0 {
		fmt.Println("No action to be applied. Aborting the process.")
		os.Exit(1)
	}

	err = s.StoreActions(actionsPtr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)	
	}

	err = s.StoreManifest(m)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)	
	}


}
package cmd

import (
	"fmt"
	"gogcli/manifest"
	"gogcli/storage"
	"os"
)

func uploadManifest(m *manifest.Manifest, s storage.Storage, concurrency int, pause int) {
	var hasActions bool
	var actions *manifest.GameActions
	var err error

	hasActions, err = s.HasActions()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if hasActions {
		fmt.Println("An unfinished manifest apply is already in progress. Aborting.")
		os.Exit(1)	
	}

	actions, err = storage.PlanManifest(m, s)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = s.StoreActions(actions)
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
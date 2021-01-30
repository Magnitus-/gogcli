package cmd

import (
	"fmt"
	"gogcli/manifest"
	"gogcli/storage"
	"os"
)

func uploadManifest(m *manifest.Manifest, s storage.Storage, concurrency int) {
	//TO FINISH
	/*actionsPtr*/ _, err := storage.PlanManifest(m, s)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
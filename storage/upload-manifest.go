package storage

import (
	"errors"
	"gogcli/manifest"
)

func UploadManifest(m *manifest.Manifest, s Storage, concurrency int, d Downloader, gamesMax int) []error {
	var hasActions bool
	var actions *manifest.GameActions
	var err error

	hasActions, err = s.HasActions()
	if err != nil {
		return []error{err}
	}
	
	if hasActions {
		return []error{errors.New("An unfinished manifest apply is already in progress. Aborting.")}
	}

	actions, err = PlanManifest(m, s)
	if err != nil {
		return []error{err}
	}

	err = s.StoreActions(actions)
	if err != nil {
		return []error{err}
	}

	err = s.StoreManifest(m)
	if err != nil {
		return []error{err}	
	}

	return processGameActions(m, actions, s, concurrency, d, gamesMax)
}
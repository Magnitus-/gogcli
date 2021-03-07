package storage

import (
	"errors"
	"gogcli/manifest"
)

func UploadManifest(m *manifest.Manifest, s Storage, src Source, concurrency int, d Downloader, gamesMax int, emptyChecksumOk bool) []error {
	var hasSource bool
	var actions *manifest.GameActions
	var err error

	hasSource, err = s.HasSource()
	if err != nil {
		return []error{err}
	}
	
	if hasSource {
		return []error{errors.New("Unfinished actions are pending in the storage. Aborting.")}
	}

	actions, err = PlanManifest(m, s, emptyChecksumOk)
	if err != nil {
		return []error{err}
	}

	if emptyChecksumOk {
		prevManifest, err := s.LoadManifest()
		if err != nil {
			return []error{err}
		}
		m.ImprintMissingChecksums(prevManifest)
	}

	err = s.StoreSource(&src)
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
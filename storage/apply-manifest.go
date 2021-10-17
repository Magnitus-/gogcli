package storage

import (
	"gogcli/manifest"
)

func ApplyManifest(m *manifest.Manifest, s Storage, src Source, emptyChecksumOk bool) error {
	var hasSource bool
	var actions *manifest.GameActions
	var err error

	hasSource, err = s.HasSource()
	if err != nil {
		return err
	}
	
	if !hasSource {
		err = s.StoreSource(&src)
		if err != nil {
			return err
		}
	}

	actions, err = PlanManifest(m, s, emptyChecksumOk)
	if err != nil {
		return err
	}

	hasManifest, hasManifestErr := s.HasManifest()
	if hasManifestErr != nil {
		return hasManifestErr
	}

	if emptyChecksumOk && hasManifest {
		prevManifest, err := s.LoadManifest()
		if err != nil {
			return err
		}
		m.ImprintMissingChecksums(prevManifest)
	}

	err = s.StoreActions(actions)
	if err != nil {
		return err
	}

	err = s.StoreManifest(m)
	if err != nil {
		return err
	}

	return nil
}
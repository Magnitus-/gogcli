package storage

import (
	"gogcli/manifest"
)

func ImprintFilter(m *manifest.Manifest, s Storage) error {
	hasManifest, hasManifestErr := s.HasManifest()
	if hasManifestErr != nil {
		return hasManifestErr
	}

	if !hasManifest {
		return nil
	}

	storedManifest, loadManifestErr := s.LoadManifest()
	if loadManifestErr != nil {
		return loadManifestErr
	}

	m.ImprintFilter(storedManifest)
	return nil
}
package storage

import (
	"errors"
	"fmt"
	"gogcli/manifest"
)

func PlanManifest(m *manifest.Manifest, s Storage, emptyChecksumOk bool) (*manifest.GameActions, error) {
	var storedManifest *manifest.Manifest
	hasManifest, err := s.HasManifest()

	if err != nil {
		msg := fmt.Sprintf("PlanManifest(...) -> Error checking manifest existance: %s", err.Error())
		return nil, errors.New(msg)
	}

	if hasManifest {
		storedManifest, err = s.LoadManifest()
		if err != nil {
			msg := fmt.Sprintf("PlanManifest(...) -> Error loading manifest: %s", err.Error())
			return nil, errors.New(msg)
		}
	} else {
	    storedManifest = manifest.NewEmptyManifest((*m).Filter)
	}

	return storedManifest.Plan(m, emptyChecksumOk, false), nil
}

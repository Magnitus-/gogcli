package storage

import (
	"gogcli/manifest"
)

func PlanManifest(m *manifest.Manifest, s Storage, emptyChecksumOk bool) (*manifest.GameActions, error) {
	var storedManifest *manifest.Manifest
	var loadManifestErr error

	hasManifest, hasManifestErr := s.HasManifest()
	if hasManifestErr != nil {
		return nil, hasManifestErr
	}

	if hasManifest {
		storedManifest, loadManifestErr = s.LoadManifest()
		if loadManifestErr != nil {
			return nil, loadManifestErr
		}
	} else {
	    storedManifest = manifest.NewEmptyManifest((*m).Filter)
	}

	hasActions, hasActionsErr := s.HasActions()
	if hasActionsErr != nil {
		return nil, hasActionsErr
	}

	if hasActions {
		storedActions, loadActionsErr := s.LoadActions()
		if loadActionsErr != nil {
			return nil, loadActionsErr
		}
		
		actionsUpdate := storedManifest.Plan(m, emptyChecksumOk, false)
		updateErr := storedActions.Update(actionsUpdate)
		if updateErr != nil {
			return nil, updateErr
		}

		return storedActions, nil
	}

	return storedManifest.Plan(m, emptyChecksumOk, false), nil
}

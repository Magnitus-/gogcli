package storage

import (
	"errors"
	"gogcli/manifest"
)

func UpdateActions(nextMan *manifest.Manifest, s Storage, emptyChecksumOk bool) error {
	hasSource, err := s.HasSource()
	if err != nil {
		return err
	}
	
	if !hasSource {
		return errors.New("No actions are pending in the storage. Aborting.")
	}

	storedMan, manErr := s.LoadManifest()
	if manErr != nil {
		return manErr
	}

	storedActions, actionsErr := s.LoadActions()
	if actionsErr != nil {
		return actionsErr
	}

	actionsUpdate := storedMan.Plan(nextMan, emptyChecksumOk, false)
	updateErr := storedActions.Update(actionsUpdate)
	if updateErr != nil {
		return updateErr
	}

	if emptyChecksumOk {
		err = nextMan.ImprintMissingChecksums(storedMan)
		if err != nil {
			return err
		}
	}

	err = s.StoreManifest(nextMan)
	if err != nil {
		return err
	}

	err = s.StoreActions(storedActions)
	if err != nil {
		return err
	}

	return nil
}
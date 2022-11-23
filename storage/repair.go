package storage

import (
	"gogcli/manifest"
)

func Repair(authMan *manifest.Manifest, storeMan *manifest.Manifest, s Storage, src Source, verifyChecksum bool) error {
	hasActions, actErr := s.HasActions()
	if actErr != nil {
		return actErr
	}
	if hasActions {
		err := s.RemoveActions()
		if err != nil {
			return err
		}
	}

	hasSource, srcErr := s.HasSource()
	if srcErr != nil {
		return srcErr
	}
	if hasSource {
		err := s.RemoveSource()
		if err != nil {
			return err
		}
	}

	manErr := s.StoreManifest(authMan)
	if manErr != nil {
		return manErr
	}

	checksumValidation := manifest.ChecksumNoValidation
	if verifyChecksum {
		checksumValidation = manifest.ChecksumValidation
	}

	actions := storeMan.Plan(authMan, checksumValidation, true)
	if len(*actions) > 0 {
		srcErr = s.StoreSource(&src)
		if srcErr != nil {
			return srcErr
		}

		actErr = s.StoreActions(actions)
		if actErr != nil {
			return actErr
		}
	}
	return nil
}

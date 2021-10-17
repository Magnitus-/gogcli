package storage

import "gogcli/manifest"

func Repair(m *manifest.Manifest, s Storage, src Source, concurrency int) error {
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

	manErr := s.StoreManifest(m)
	if manErr != nil {
		return manErr
	}

	l, lErr := s.GetListing()
	if lErr != nil {
		return lErr
	}

	filesManifest, err := l.GetManifest(concurrency)
	if err != nil {
		return err
	}

	actions := filesManifest.Plan(m, false, true)
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

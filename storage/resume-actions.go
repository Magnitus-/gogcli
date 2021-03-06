package storage

import (
	"errors"
	"fmt"
	"gogcli/manifest"
)

func ResumeActions(s Storage, concurrency int, d Downloader, gamesMax int, gamesSort manifest.ActionsIteratorSort) []error {
	exists, err := s.Exists()
	if err != nil {
		return []error{err}
	}

	if !exists {
		msg := fmt.Sprintf("Storage %s does not exist", s.GetPrintableSummary())
		return []error{errors.New(msg)}
	}

	hasManifest, hasManifestErr := s.HasManifest()
	if hasManifestErr != nil {
		return []error{hasManifestErr}
	}

	if !hasManifest {
		msg := fmt.Sprintf("Storage %s does not have a manifest", s.GetPrintableSummary())
		return []error{errors.New(msg)}
	}

	manifest, manifestErr := s.LoadManifest()
	if manifestErr != nil {
		return []error{manifestErr}
	}

	hasActions, hasActionsErr := s.HasActions()
	if hasActionsErr != nil {
		return []error{hasActionsErr}
	}

	if !hasActions {
		msg := fmt.Sprintf("Storage %s does not actions to execute", s.GetPrintableSummary())
		return []error{errors.New(msg)}
	}

	actions, actionsErr := s.LoadActions()
	if err != nil {
		return []error{actionsErr}
	}

	return processGameActions(manifest, actions, s, concurrency, d, gamesMax, gamesSort)
}

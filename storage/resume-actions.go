package storage

import (
	"errors"
	"fmt"
)

func ResumeActions(s Storage, d Downloader, a ActionsProcessor) []error {
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

	return a.ProcessGameActions(manifest, actions, s, d)
}

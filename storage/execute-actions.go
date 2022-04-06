package storage

import (
	"errors"
	"fmt"
)

func ExecuteActions(s Storage, d Downloader, a ActionsProcessor) []error {
	exists, err := s.Exists()
	if err != nil {
		return []error{err}
	}

	summary, summaryErr := s.GetPrintableSummary()

	if !exists {
		if summaryErr != nil {
			return []error{errors.New("Storage does not exist")}
		}
		msg := fmt.Sprintf("Storage %s does not exist", summary)
		return []error{errors.New(msg)}
	}

	hasManifest, hasManifestErr := s.HasManifest()
	if hasManifestErr != nil {
		return []error{hasManifestErr}
	}

	if !hasManifest {
		if summaryErr != nil {
			return []error{errors.New("Storage does not have a manifest")}
		}
		msg := fmt.Sprintf("Storage %s does not have a manifest", summary)
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
		if summaryErr != nil {
			return []error{errors.New("Storage does not have actions to execute")}
		}
		msg := fmt.Sprintf("Storage %s does not have actions to execute", summary)
		return []error{errors.New(msg)}
	}

	actions, actionsErr := s.LoadActions()
	if err != nil {
		return []error{actionsErr}
	}

	return a.ProcessGameActions(manifest, actions, s, d)
}

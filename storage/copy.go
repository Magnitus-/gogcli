package storage

import (
	"errors"
	"fmt"
)

func Copy(source Storage, destination Storage, sourceDownloader Downloader, a ActionsProcessor) []error {
	exists, err := source.Exists()
	if err != nil {
		return []error{err}
	}

	if !exists {
		msg := fmt.Sprintf("Source storage %s does not exist", source.GetPrintableSummary())
		return []error{errors.New(msg)}
	}

	hasSource, hasSourceErr := source.HasSource()
	if hasSourceErr != nil {
		return []error{hasSourceErr}
	}

	if hasSource {
		return []error{errors.New("Unfinished actions are pending in the source storage. Aborting.")}
	}

	err = EnsureInitialization(destination)
	if err != nil {
		return []error{err}
	}

	hasSource, hasSourceErr = destination.HasSource()
	if hasSourceErr != nil {
		return []error{hasSourceErr}
	}

	if hasSource {
		return []error{errors.New("Unfinished actions are pending in the desination storage. Aborting.")}
	}

	m, loadErr := source.LoadManifest()
	if loadErr != nil {
		return []error{loadErr}
	}

	err = ApplyManifest(m, destination, *source.GenerateSource(), false)
	if err != nil {
		return []error{err}
	}

	actions, loadActionsErr := destination.LoadActions()
	if loadActionsErr != nil {
		return []error{loadActionsErr}
	}

	return a.ProcessGameActions(m, actions, destination, sourceDownloader)
}

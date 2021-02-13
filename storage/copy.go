package storage

import (
	"errors"
	"fmt"
)

func Copy(source Storage, destination Storage, sourceDownloader Downloader, concurrency int) []error {
	errs := make([]error, 0)

	exists, err := source.Exists()
	if err != nil {
		errs := append(errs, err)
		return errs
	}

	if !exists {
		msg := fmt.Sprintf("Source storage %s does not exist", source.GetPrintableSummary())
		errs := append(errs, errors.New(msg))
		return errs
	}

	err = EnsureInitialization(destination)
	if err != nil {
		errs := append(errs, err)
		return errs
	}

	m, loadErr := source.LoadManifest()
	if loadErr != nil {
		errs := append(errs, loadErr)
		return errs
	}

	errs = UploadManifest(m, source, concurrency, sourceDownloader)
	return errs
}
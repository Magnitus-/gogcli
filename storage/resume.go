package storage

import (
	"errors"
	"fmt"
)

func ResumeUpload(s Storage, concurrency int, d Downloader, gamesMax int) []error {
	errs := make([]error, 0)
	
	exists, err := s.Exists()
	if err != nil {
		errs := append(errs, err)
		return errs
	}

	if !exists {
		msg := fmt.Sprintf("Storage %s does not exist", s.GetPrintableSummary())
		errs := append(errs, errors.New(msg))
		return errs
	}

	hasManifest, hasManifestErr := s.HasManifest()
	if hasManifestErr != nil {
		errs := append(errs, hasManifestErr)
		return errs
	}

	if !hasManifest {
		msg := fmt.Sprintf("Storage %s does not have a manifest", s.GetPrintableSummary())
		errs := append(errs, errors.New(msg))
		return errs
	}

	manifest, manifestErr := s.LoadManifest()
	if manifestErr != nil {
		errs := append(errs, manifestErr)
		return errs
	}

	return UploadManifest(manifest, s, concurrency, d, gamesMax)
}

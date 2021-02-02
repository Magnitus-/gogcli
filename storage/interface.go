package storage

import (
	"errors"
	"fmt"
	"gogcli/manifest"
	"io"
)

type Storage interface {
	HasManifest() (bool, error)
	HasActions() (bool, error)
	StoreManifest(m *manifest.Manifest) error
	StoreActions(a *manifest.GameActions) error
	LoadManifest() (*manifest.Manifest, error)
	LoadActions() (*manifest.GameActions, error)
	AddGame(gameId int) error
	RemoveGame(gameId int) error
	UploadFile(source io.ReadCloser, gameId int, kind string, name string) ([]byte, error)
	RemoveFile(gameId int, kind string, name string) error
}

func PlanManifest(m *manifest.Manifest, s Storage) (*manifest.GameActions, error) {
	var storedManifest *manifest.Manifest
	hasManifest, err := s.HasManifest()

	if err != nil {
		msg := fmt.Sprintf("PlanManifest(...) -> Error checking manifest existance: %s", err.Error())
		return nil, errors.New(msg)
	}

	if hasManifest {
		storedManifest, err = s.LoadManifest()
		if err != nil {
			msg := fmt.Sprintf("PlanManifest(...) -> Error loading manifest: %s", err.Error())
			return nil, errors.New(msg)
		}
	} else {
	    storedManifest = manifest.NewEmptyManifest()
	}

	return storedManifest.Plan(m), nil
}

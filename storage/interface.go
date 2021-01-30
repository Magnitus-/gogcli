package storage

import (
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
	storedManifestPtr := manifest.NewEmptyManifest()
	hasManifest, err := s.HasManifest()

	if err != nil {
		return nil, err
	}

	if hasManifest {
		err = s.StoreManifest(storedManifestPtr)
		if err != nil {
			return nil, err
		}
	}

	return storedManifestPtr.Plan(m), nil
}

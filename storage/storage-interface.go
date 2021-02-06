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
	RemoveActions() error
	AddGame(gameId int) error
	RemoveGame(gameId int) error
	UploadFile(source io.ReadCloser, gameId int, kind string, name string) (string, error)
	RemoveFile(gameId int, kind string, name string) error
}
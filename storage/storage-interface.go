package storage

import (
	"gogcli/manifest"
	"io"
)

type Storage interface {
	GetPrintableSummary() string
	Exists() (bool, error)
	Initialize() error
	HasManifest() (bool, error)
	HasActions() (bool, error)
	StoreManifest(m *manifest.Manifest) error
	StoreActions(a *manifest.GameActions) error
	LoadManifest() (*manifest.Manifest, error)
	LoadActions() (*manifest.GameActions, error)
	RemoveActions() error
	AddGame(gameId int) error
	RemoveGame(gameId int) error
	UploadFile(source io.ReadCloser, gameId int, kind string, name string, expectedSize int64) (string, error)
	RemoveFile(gameId int, kind string, name string) error
	DownloadFile(gameId int, kind string, name string) (io.ReadCloser, int64, error)
}
package storage

import (
	"gogcli/manifest"
	"io"
)

type Storage interface {
	GetListing() (*StorageListing, error)
	SupportsReaderAt() bool
	GenerateSource() *Source
	GetPrintableSummary() string
	Exists() (bool, error)
	Initialize() error
	HasManifest() (bool, error)
	HasActions() (bool, error)
	HasSource() (bool, error)
	StoreManifest(m *manifest.Manifest) error
	StoreActions(a *manifest.GameActions) error
	StoreSource(s *Source) error
	LoadManifest() (*manifest.Manifest, error)
	LoadActions() (*manifest.GameActions, error)
	LoadSource() (*Source, error)
	RemoveActions() error
	RemoveSource() error
	AddGame(gameId int64) error
	RemoveGame(gameId int64) error
	UploadFile(source io.ReadCloser, gameId int64, kind string, name string, expectedSize int64) (string, error)
	RemoveFile(gameId int64, kind string, name string) error
	DownloadFile(gameId int64, kind string, name string) (io.ReadCloser, int64, error)
}
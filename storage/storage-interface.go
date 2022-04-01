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
	AddGame(game manifest.GameInfo) error
	RemoveGame(game manifest.GameInfo) error
	UploadFile(source io.ReadCloser, file manifest.FileInfo) (string, error)
	RemoveFile(file manifest.FileInfo) error
	DownloadFile(file manifest.FileInfo) (io.ReadCloser, int64, error)
}

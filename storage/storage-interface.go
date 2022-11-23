package storage

import (
	"gogcli/manifest"
	"gogcli/metadata"
	"io"
)

type Storage interface {
	GetGameIds() ([]int64, error)
	GetGameFiles(GameId int64) ([]manifest.FileInfo, error)
	SupportsReaderAt() bool
	IsSelfValidating() (bool, error)
	GenerateSource() *Source
	GetPrintableSummary() (string, error)
	Exists() (bool, error)
	Initialize() error
	HasManifest() (bool, error)
	HasMetadata() (bool, error)
	HasActions() (bool, error)
	HasSource() (bool, error)
	StoreManifest(m *manifest.Manifest) error
	StoreMetadata(m *metadata.Metadata) error
	StoreActions(a *manifest.GameActions) error
	StoreSource(s *Source) error
	LoadManifest() (*manifest.Manifest, error)
	LoadMetadata() (*metadata.Metadata, error)
	LoadActions() (*manifest.GameActions, error)
	LoadSource() (*Source, error)
	RemoveActions() error
	RemoveSource() error
	AddGame(game manifest.GameInfo) error
	RemoveGame(game manifest.GameInfo) error
	UploadFile(source io.ReadCloser, file manifest.FileInfo) (string, error)
	RemoveFile(file manifest.FileInfo) error
	DownloadFile(file manifest.FileInfo) (io.ReadCloser, int64, error)
	UploadImage(source io.ReadCloser, image metadata.GameMetadataImage) (string, error)
	RemoveImage(image metadata.GameMetadataImage) error
	DownloadImage(image metadata.GameMetadataImage) (io.ReadCloser, int64, error)
}

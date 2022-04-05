package storage

import (
	"gogcli/manifest"
	_ "gogcli/storagegrpc"
	"io"
)

type GrpcStore struct {
	Endpoint  string
}

func (g GrpcStore) GetListing() (*StorageListing, error) {
	return nil, nil
}

func (g GrpcStore) SupportsReaderAt() bool {
	return true
}

func (g GrpcStore) IsSelfVerifying() bool {
	return false
}

func (g GrpcStore) GenerateSource() *Source {
	return nil
}

func (g GrpcStore) GetPrintableSummary() string {
	return ""
}

func (g GrpcStore) Exists() (bool, error) {
	return false, nil
}

func (g GrpcStore) Initialize() error {
	return nil
}

func (g GrpcStore) HasManifest() (bool, error) {
	return false, nil
}

func (g GrpcStore) HasActions() (bool, error) {
	return false, nil
}

func (g GrpcStore) HasSource() (bool, error) {
	return false, nil
}

func (g GrpcStore) StoreManifest(m *manifest.Manifest) error {
	return nil
}

func (g GrpcStore) StoreActions(a *manifest.GameActions) error {
	return nil
}

func (g GrpcStore) StoreSource(s *Source) error {
	return nil
}

func (g GrpcStore) LoadManifest() (*manifest.Manifest, error) {
	return nil, nil
}

func (g GrpcStore) LoadActions() (*manifest.GameActions, error) {
	return nil, nil
}

func (g GrpcStore) LoadSource() (*Source, error) {
	return nil, nil
}

func (g GrpcStore) RemoveActions() error {
	return nil
}

func (g GrpcStore) RemoveSource() error {
	return nil
}

func (g GrpcStore) AddGame(game manifest.GameInfo) error {
	return nil
}

func (g GrpcStore) RemoveGame(game manifest.GameInfo) error {
	return nil
}

func (g GrpcStore) UploadFile(source io.ReadCloser, file manifest.FileInfo) (string, error) {
	return "", nil
}

func (g GrpcStore) RemoveFile(file manifest.FileInfo) error {
	return nil
}

func (g GrpcStore) DownloadFile(file manifest.FileInfo) (io.ReadCloser, int64, error) {
	return nil, 0, nil
}

package storage

import (
	"gogcli/manifest"
	_ "gogcli/storageplugin"
	"io"
)

type PluginStore struct {
	Endpoint  string
}

func (p PluginStore) GetListing() (*StorageListing, error) {
	return nil, nil
}

func (p PluginStore) SupportsReaderAt() bool {
	return true
}

func (p PluginStore) IsSelfVerifying() bool {
	return false
}

func (p PluginStore) GenerateSource() *Source {
	return nil
}

func (p PluginStore) GetPrintableSummary() string {
	return ""
}

func (p PluginStore) Exists() (bool, error) {
	return false, nil
}

func (p PluginStore) Initialize() error {
	return nil
}

func (p PluginStore) HasManifest() (bool, error) {
	return false, nil
}

func (p PluginStore) HasActions() (bool, error) {
	return false, nil
}

func (p PluginStore) HasSource() (bool, error) {
	return false, nil
}

func (p PluginStore) StoreManifest(m *manifest.Manifest) error {
	return nil
}

func (p PluginStore) StoreActions(a *manifest.GameActions) error {
	return nil
}

func (p PluginStore) StoreSource(s *Source) error {
	return nil
}

func (p PluginStore) LoadManifest() (*manifest.Manifest, error) {
	return nil, nil
}

func (p PluginStore) LoadActions() (*manifest.GameActions, error) {
	return nil, nil
}

func (p PluginStore) LoadSource() (*Source, error) {
	return nil, nil
}

func (p PluginStore) RemoveActions() error {
	return nil
}

func (p PluginStore) RemoveSource() error {
	return nil
}

func (p PluginStore) AddGame(game manifest.GameInfo) error {
	return nil
}

func (p PluginStore) RemoveGame(game manifest.GameInfo) error {
	return nil
}

func (p PluginStore) UploadFile(source io.ReadCloser, file manifest.FileInfo) (string, error) {
	return "", nil
}

func (p PluginStore) RemoveFile(file manifest.FileInfo) error {
	return nil
}

func (p PluginStore) DownloadFile(file manifest.FileInfo) (io.ReadCloser, int64, error) {
	return nil, 0, nil
}

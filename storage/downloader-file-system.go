package storage

import (
	"gogcli/manifest"
	"io"
)

//Implementation of the Downloader interface
type FileSystemDownloader struct {
	Fs FileSystem
}

func (d FileSystemDownloader) Download(gameId int, add manifest.FileAction) (io.ReadCloser, int, string, error) {
	handle, size, err := d.Fs.DownloadFile(gameId, add.Kind, add.Name)
	return handle, size, add.Name, err
}
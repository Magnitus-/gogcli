package storage

import (
	"gogcli/manifest"
	"io"
)

//Implementation of the Downloader interface
type S3StoreDownloader struct {
	S3 S3Store
}

func (d S3StoreDownloader) Download(gameId int64, add manifest.FileAction) (io.ReadCloser, int64, string, error) {
	handle, size, err := d.S3.DownloadFile(gameId, add.Kind, add.Name)
	return handle, size, add.Name, err
}

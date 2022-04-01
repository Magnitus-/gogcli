package storage

import (
	"gogcli/manifest"
	"io"
)

//Implementation of the Downloader interface
type S3StoreDownloader struct {
	S3 S3Store
}

func (d S3StoreDownloader) Download(file manifest.FileInfo) (io.ReadCloser, int64, string, error) {
	handle, size, err := d.S3.DownloadFile(file)
	return handle, size, file.Name, err
}

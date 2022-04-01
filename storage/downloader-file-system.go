package storage

import (
	"gogcli/manifest"
	"io"
)

//Implementation of the Downloader interface
type FileSystemDownloader struct {
	Fs FileSystem
}

func (d FileSystemDownloader) Download(file manifest.FileInfo) (io.ReadCloser, int64, string, error) {
	handle, size, err := d.Fs.DownloadFile(file)
	return handle, size, file.Name, err
}

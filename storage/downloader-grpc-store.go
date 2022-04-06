package storage

import (
	"gogcli/manifest"
	"io"
)

//Implementation of the Downloader interface
type GrpcStoreDownloader struct {
	Grpc GrpcStore
}

func (d GrpcStoreDownloader) Download(file manifest.FileInfo) (io.ReadCloser, int64, string, error) {
	handle, size, err := d.Grpc.DownloadFile(file)
	return handle, size, file.Name, err
}

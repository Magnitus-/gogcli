package sdk

import (
	"gogcli/manifest"
	"io"
)

//Implementation of the Downloader interface
type Downloader struct {
	SdkPtrPtr *Sdk
}

func (d Downloader) Download(file manifest.FileInfo) (io.ReadCloser, int64, string, error) {
	return d.SdkPtrPtr.GetDownloadHandle(file.Url)
}

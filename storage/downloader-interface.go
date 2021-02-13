package storage

import (
	"gogcli/manifest"
	"io"
)

type Downloader interface {
	Download(int, manifest.FileAction) (io.ReadCloser, int64, string, error)
}
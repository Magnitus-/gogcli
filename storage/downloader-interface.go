package storage

import (
	"gogcli/manifest"
	"io"
)

type Downloader interface {
	Download(int, manifest.FileAction) (io.ReadCloser, int64, string, error)
}

type Source struct {
	Type string
	S3Params S3Configs
	FsPath string
}
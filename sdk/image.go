package sdk

import (
	"fmt"
	"path"
)

type ImageInfo struct {
	Name string
	Size int64
	Checksum string
}

func (s *Sdk) GetImageInfo(url string) (ImageInfo, error) {
	fn := fmt.Sprintf("GetImageInfo(url=%s)", url)

	reply, err := s.getUrlBodyChecksum(url, fn, (*s).maxRetries)
	if err != nil {
		return ImageInfo{}, err
	}

	return ImageInfo{
		Name: path.Base(url),
		Checksum: reply.BodyChecksum,
		Size: reply.BodyLength,
	}, nil
}
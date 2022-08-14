package sdk

import (
	"fmt"
	"path"
	"strings"

	"gogcli/metadata"
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

func (s *Sdk) GetImageInfoWithFallback(image *metadata.GameMetadataImage) (ImageInfo, error) {
	info, err := s.GetImageInfo((*image).Url)
	if err != nil {
		if strings.Contains((*image).Url, "ggvgl_2x") {
			adjustedUrl := strings.Replace((*image).Url, "ggvgl_2x", "ggvgl", -1)
			(*s).logger.Warning(fmt.Sprintf("Unretrievable image %s: Will retry with the following version of the image that has worse resolution: %s.", (*image).Url, adjustedUrl))
			(*image).Url = adjustedUrl
			
			return s.GetImageInfo((*image).Url)
		}

		return info, err
	}

	return info, nil
}
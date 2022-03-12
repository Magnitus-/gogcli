package sdk

import (
	"fmt"
	"time"
)

type DownloadFileInfo struct {
	url      string
	name     string
	checksum string
	size     int64
}

type DownloadFileInfoReturn struct {
	url         string
	name        string
	checksum    string
	size        int64
	err         error
	dangling    bool
	badMetadata bool
}

func (s *Sdk) GetDownloadFilenInfoAsync(downloadPath string, tolerateBadFileMetadata bool, returnVal chan DownloadFileInfoReturn) {
	name, checksum, size, err, dangling, badMetadata := s.GetDownloadFileInfo(downloadPath)
	if badMetadata && tolerateBadFileMetadata {
		var workaroundErr error
		name, checksum, size, workaroundErr = s.GetDownloadFileInfoWorkaroundWay(downloadPath)
		if workaroundErr != nil {
			returnVal <- DownloadFileInfoReturn{url: downloadPath, name: name, checksum: checksum, size: size, err: workaroundErr, dangling: false, badMetadata: false}
			return
		}
		returnVal <- DownloadFileInfoReturn{url: downloadPath, name: name, checksum: checksum, size: size, err: err, dangling: false, badMetadata: true}
		return
	}
	returnVal <- DownloadFileInfoReturn{url: downloadPath, name: name, checksum: checksum, size: size, err: err, dangling: dangling, badMetadata: false}
}

func (s *Sdk) GetManyDownloadFileInfo(downloadPaths []string, concurrency int, pause int, tolerateDangles bool, tolerateBadFileMetadata bool) ([]DownloadFileInfo, []error, []error) {
	var errs []error
	var warnings []error
	var downloadFileInfos []DownloadFileInfo
	c := make(chan DownloadFileInfoReturn)

	i := 0
	for i < len(downloadPaths) {
		beginning := i
		target := min(len(downloadPaths), i+concurrency)
		for i < target {
			go s.GetDownloadFilenInfoAsync(downloadPaths[i], tolerateBadFileMetadata, c)
			i++
		}

		y := beginning
		for y < target {
			returnVal := <-c
			if returnVal.err != nil {
				if returnVal.badMetadata && tolerateBadFileMetadata {
					(*s).logger.Warning(fmt.Sprintf("Bad metadata for %s: File metadata was still fetched using much longer workaround method.", returnVal.url))
					warnings = append(warnings, returnVal.err)
					downloadFileInfos = append(downloadFileInfos, DownloadFileInfo{url: returnVal.url, name: returnVal.name, checksum: returnVal.checksum, size: returnVal.size})
				} else if returnVal.dangling && tolerateDangles {
					(*s).logger.Warning(fmt.Sprintf("Bad download link for %s: File was not added to manifest.", returnVal.url))
					warnings = append(warnings, returnVal.err)
				} else {
					errs = append(errs, returnVal.err)
				}
			} else {
				downloadFileInfos = append(downloadFileInfos, DownloadFileInfo{url: returnVal.url, name: returnVal.name, checksum: returnVal.checksum, size: returnVal.size})
			}
			y++
		}

		if len(errs) > 0 {
			return downloadFileInfos, errs, warnings
		}

		if i < len(downloadPaths) {
			time.Sleep(time.Duration(pause) * time.Millisecond)
		}
	}

	return downloadFileInfos, nil, warnings
}

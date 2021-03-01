package sdk

import "time"

type DownloadFileInfo struct {
	url      string
	name     string
	checksum string
	size     int64
}

type DownloadFileInfoReturn struct {
	url      string
	name     string
	checksum string
	size     int64
	err      error
	dangling bool
}

func (s *Sdk) GetDownloadFilenInfoAsync(downloadPath string, returnVal chan DownloadFileInfoReturn) {
	name, checksum, size, err, dangling := s.GetDownloadFileInfo(downloadPath)
	returnVal <- DownloadFileInfoReturn{url: downloadPath, name: name, checksum: checksum, size: size, err: err, dangling: dangling}
}

func (s *Sdk) GetManyDownloadFileInfo(downloadPaths []string, concurrency int, pause int, tolerateDangles bool) ([]DownloadFileInfo, []error, []error) {
	var errs []error
	var danglingErrs []error
	var downloadFileInfos []DownloadFileInfo
	c := make(chan DownloadFileInfoReturn)

	i := 0
	for i < len(downloadPaths) {
		beginning := i
		target := min(len(downloadPaths), i+concurrency)
		for i < target {
			go s.GetDownloadFilenInfoAsync(downloadPaths[i], c)
			i++
		}

		y := beginning
		for y < target {
			returnVal := <-c
			if returnVal.err != nil {
				if returnVal.dangling && tolerateDangles {
					danglingErrs = append(danglingErrs, returnVal.err)
				} else {
					errs = append(errs, returnVal.err)
				}
			} else {
				downloadFileInfos = append(downloadFileInfos, DownloadFileInfo{url: returnVal.url, name: returnVal.name, checksum: returnVal.checksum, size: returnVal.size})
			}
			y++
		}

		if len(errs) > 0 {
			return downloadFileInfos, errs, danglingErrs
		}

		if i < len(downloadPaths) {
			time.Sleep(time.Duration(pause) * time.Millisecond)
		}
	}

	return downloadFileInfos, nil, danglingErrs
}
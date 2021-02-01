package sdk

import "time"

type DownloadFilename struct {
	url      string
	filename string
}

type DownloadFilenameReturn struct {
	url      string
	filename string
	err      error
}

func (s *Sdk) GetDownloadFilenameAsync(downloadPath string, debug bool, returnVal chan DownloadFilenameReturn) {
	filename, err := s.GetDownloadFilename(downloadPath, debug)
	returnVal <- DownloadFilenameReturn{url: downloadPath, filename: filename, err: err}
}

func (s *Sdk) GetManyDownloadFilename(downloadPaths []string, concurrency int, pause int, debug bool) ([]DownloadFilename, []error) {
	var errs []error
	var downloadFilenames []DownloadFilename
	c := make(chan DownloadFilenameReturn)

	i := 0
	for i < len(downloadPaths) {
		beginning := i
		target := min(len(downloadPaths), i+concurrency)
		for i < target {
			go s.GetDownloadFilenameAsync(downloadPaths[i], debug, c)
			i++
		}

		y := beginning
		for y < target {
			returnVal := <-c
			if returnVal.err != nil {
				errs = append(errs, returnVal.err)
			} else {
				downloadFilenames = append(downloadFilenames, DownloadFilename{url: returnVal.url, filename: returnVal.filename})
			}
			y++
		}

		if len(errs) > 0 {
			return downloadFilenames, errs
		}

		if i < len(downloadPaths) {
			time.Sleep(time.Duration(pause) * time.Millisecond)
		}
	}

	return downloadFilenames, nil
}
package storage

import (
	"archive/zip"
	"errors"
	"fmt"
	"gogcli/manifest"
	"io"
	"io/ioutil"
)

func ValidateZipArchive(s Storage, file manifest.FileInfo) error {
	fn := fmt.Sprintf("ValidateZipArchive(..., file={game={Id=%d, ...}, Kind=%s, Name=%s, ...})", file.Game.Id, file.Kind, file.Name)
	if !s.SupportsReaderAt() {
		msg := fmt.Sprintf("%s -> Provided storage doesn't support downloading fixed length subset of file from a given offset", fn)
		return errors.New(msg)
	}

	download, size, err := s.DownloadFile(file)
	if err != nil {
		msg := fmt.Sprintf("%s -> Error occurred getting download from store: %s", fn, err.Error())
		return errors.New(msg)
	}

	downloadReaderAt, ok := download.(io.ReaderAt)
	if !ok {
		msg := fmt.Sprintf("%s -> Provided download doesn't support downloading fixed length subset of file from a given offset", fn)
		return errors.New(msg)
	}
	defer download.Close()

	zipReader, zipErr := zip.NewReader(downloadReaderAt, size)
	if zipErr != nil {
		msg := fmt.Sprintf("%s -> Error occurred opening the zip archive: %s", fn, zipErr.Error())
		return errors.New(msg)
	}

	for _, zipFile := range zipReader.File {
		zipFileReader, openErr := zipFile.Open()
		if openErr != nil {
			msg := fmt.Sprintf("%s -> Error occurred accessing file %s in zip archive: %s", fn, zipFile.Name, openErr.Error())
			return errors.New(msg)
		}
		defer zipFileReader.Close()

		_, err = io.Copy(ioutil.Discard, zipFileReader)
		if err != nil {
			msg := fmt.Sprintf("%s -> Error occurred reading content of file %s in zip archive: %s", fn, zipFile.Name, err.Error())
			return errors.New(msg)
		}
	}

	return nil
}

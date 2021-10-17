package storage

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"gogcli/manifest"
	"io"
)

func validateFile(info manifest.FileInfo, s Storage, errChan chan error) {
	downloadHandle, size, err := s.DownloadFile(info.GameId, info.Kind, info.Name)
	if err != nil {
		msg := fmt.Sprintf("validateFile(FileInfo{GameId: %d, Kind: %s, Name: %s}, ...) -> Error occured while getting the file's download handle: %s", info.GameId, info.Kind, info.Name, err.Error())
		errChan <- errors.New(msg)
		return
	}

	h := md5.New()
	io.Copy(h, downloadHandle)
	checksum := hex.EncodeToString(h.Sum(nil))

	if size != info.Size {
		msg := fmt.Sprintf("validateFile(FileInfo{GameId: %d, Kind: %s, Name: %s}, ...) -> Actual file size of %d did not match the expected size of %d", info.GameId, info.Kind, info.Name, size, info.Size)
		errChan <- errors.New(msg)
		return
	}

	if size != info.Size {
		msg := fmt.Sprintf("validateFile(FileInfo{GameId: %d, Kind: %s, Name: %s}, ...) -> Actual file size of %d did not match the expected size of %d", info.GameId, info.Kind, info.Name, size, info.Size)
		errChan <- errors.New(msg)
		return
	}

	if checksum != info.Checksum {
		msg := fmt.Sprintf("validateFile(FileInfo{GameId: %d, Kind: %s, Name: %s}, ...) -> Actual file checksum of %s did not match the expected checksum of %s", info.GameId, info.Kind, info.Name, checksum, info.Checksum)
		errChan <- errors.New(msg)
		return
	}

	errChan <- nil
}

func ValidateManifest(s Storage, concurrency int) []error {
	jobsRunning := 0
	errChan := make(chan error)

	errs := make([]error, 0)
	has, err := s.HasManifest()
	if err != nil {
		msg := fmt.Sprintf("ValidateManifest(...) -> Error checking manifest existance: %s", err.Error())
		errs = append(errs, errors.New(msg))
		return errs
	} else if !has {
		msg := fmt.Sprintf("ValidateManifest(...) -> Manifest not found")
		errs = append(errs, errors.New(msg))
		return errs
	}

	m, loadErr := s.LoadManifest()
	if err != nil {
		msg := fmt.Sprintf("ValidateManifest(...) -> Error occured while loading the manifest: %s", loadErr.Error())
		errs = append(errs, errors.New(msg))
		return errs
	}

	iterator := manifest.NewManifestFileInterator(m)
	for true {
		if jobsRunning > 0 && ((!iterator.HasMore()) || concurrency <= 0) {
			err := <-errChan
			if err != nil {
				errs = append(errs, err)
			}
			jobsRunning--
		} else if !iterator.HasMore() {
			break
		}
		if iterator.HasMore() {
			file, err := iterator.Next()
			if err != nil {
				errs = append(errs, err)
			} else {
				go validateFile(file, s, errChan)
				jobsRunning++
				concurrency--
			}
		}
	}

	return errs
}

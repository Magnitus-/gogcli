package storage

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"encoding/hex"
	"errors"
	"fmt"
	"gogcli/manifest"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
)

type FileSystem struct {
	Path string
	debug bool
	logger *log.Logger
}

func GetFileSystemFromSource(s Source, debug bool, tag string) (FileSystem, error) {
	if s.Type != "fs" {
		msg := fmt.Sprintf("Cannot load file system from source of type %s", s.Type)
		return FileSystem{"", false, nil}, errors.New(msg)
	}
	return GetFileSystem(s.FsPath, debug, tag), nil
}

func GetFileSystem(path string, debug bool, tag string) FileSystem {
	var logPrefix string
	if tag == "" {
		logPrefix = "FS: "
	} else {
		logPrefix = fmt.Sprintf("FS-%s: ", tag)
	}
	return FileSystem{path, debug, log.New(os.Stdout, logPrefix, log.Lshortfile)}
}

func (f FileSystem) GenerateSource() *Source {
	src := Source{Type: "fs", FsPath: f.Path}
	return &src
}

func (f FileSystem) GetPrintableSummary() string {
	return fmt.Sprintf("FileSystem{Path: %s}", f.Path)
}

func (f FileSystem) Exists() (bool, error) {
	_, err := os.Stat(f.Path)
	if err != nil {
		if os.IsNotExist(err) {
			if f.debug {
				f.logger.Println("Exists() -> File system path not found")
			}
			return false, nil
		} else {
			msg := fmt.Sprintf("Exists() -> The following error occured while ascertaining existance of path %s: %s", f.Path, err.Error())
			return true, errors.New(msg)
		}
	}

	if f.debug {
		f.logger.Println("Exists() -> File system path found")
	}
	return true, nil
}

func (f FileSystem) Initialize() error {
	err := os.MkdirAll(f.Path, 0755)
	if err != nil {
		msg := fmt.Sprintf("Initialize() -> Failed to create a directory at the specified path: %s", err.Error())
		return errors.New(msg)
	}

	if f.debug {
		msg := fmt.Sprintf("Initialize() -> File system path %s created", f.Path)
		f.logger.Println(msg)
	}
	return nil
}

func (f FileSystem) HasManifest() (bool, error) {
	_, err := os.Stat(path.Join(f.Path, "manifest.json"))
	if err != nil {
		if os.IsNotExist(err) {
			if f.debug {
				f.logger.Println("HasManifest() -> Manifest not found")
			}
			return false, nil
		}

		msg := fmt.Sprintf("HasManifest() -> The following error occured while ascertaining manifest's existance: %s", err.Error())
		return true, errors.New(msg)
	}
	if f.debug {
		f.logger.Println("HasManifest() -> Manifest found")
	}
	return true, nil
}

func (f FileSystem) HasActions() (bool, error) {
	_, err := os.Stat(path.Join(f.Path, "actions.json"))
	if err != nil {
		if os.IsNotExist(err) {
			if f.debug {
				f.logger.Println("HasActions() -> Actions not found")
			}
			return false, nil
		}

		msg := fmt.Sprintf("HasActions() -> The following error occured while ascertaining actions' existance: %s", err.Error())
		return true, errors.New(msg)
	}
	if f.debug {
		f.logger.Println("HasActions() -> Actions found")
	}
	return true, nil
}

func (f FileSystem) HasSource() (bool, error) {
	_, err := os.Stat(path.Join(f.Path, "source.json"))
	if err != nil {
		if os.IsNotExist(err) {
			if f.debug {
				f.logger.Println("HasSource() -> Source not found")
			}
			return false, nil
		}

		msg := fmt.Sprintf("HasSource() -> The following error occured while ascertaining source's existance: %s", err.Error())
		return true, errors.New(msg)
	}
	if f.debug {
		f.logger.Println("HasSource() -> Source found")
	}
	return true, nil
}

func (f FileSystem) StoreManifest(m *manifest.Manifest) error {
	var err error
	var buf bytes.Buffer
	var output []byte

	output, err = json.Marshal(*m)

	if err != nil {
		return err
	}

	json.Indent(&buf, output, "", "  ")
	output = buf.Bytes()

	err = ioutil.WriteFile(path.Join(f.Path, "manifest.json"), output, 0644)
	if err == nil && f.debug {
		f.logger.Println(fmt.Sprintf("StoreManifest(...) -> Stored manifest with %d games", len((*m).Games)))
	}
	return err
}

func (f FileSystem) StoreActions(a *manifest.GameActions) error {
	var err error
	var buf bytes.Buffer
	var output []byte

	output, err = json.Marshal(*a)

	if err != nil {
		return err
	}

	json.Indent(&buf, output, "", "  ")
	output = buf.Bytes()

	err = ioutil.WriteFile(path.Join(f.Path, "actions.json"), output, 0644)
	if err == nil && f.debug {
		f.logger.Println(fmt.Sprintf("StoreActions(...) -> Stored actions on %d games", len(*a)))
	}
	return err
}

func (f FileSystem) StoreSource(s *Source) error {
	var err error
	var buf bytes.Buffer
	var output []byte

	output, err = json.Marshal(*s)

	if err != nil {
		return err
	}

	json.Indent(&buf, output, "", "  ")
	output = buf.Bytes()

	err = ioutil.WriteFile(path.Join(f.Path, "source.json"), output, 0644)
	if err == nil && f.debug {
		f.logger.Println(fmt.Sprintf("StoreSource(...) -> Stored source of type %s", s.Type))
	}
	return err
}

func (f FileSystem) LoadManifest() (*manifest.Manifest, error) {
	var m manifest.Manifest

	bs, err := ioutil.ReadFile(path.Join(f.Path, "manifest.json"))
	if err != nil {
		return &m, err
	}

	err = json.Unmarshal(bs, &m)
	if err != nil {
		return &m, err
	}

	if f.debug {
		f.logger.Println(fmt.Sprintf("LoadManifest() -> Loaded manifest with %d games", len(m.Games)))
	}
	return &m, nil
}

func (f FileSystem) LoadActions() (*manifest.GameActions, error) {
	var a *manifest.GameActions

	bs, err := ioutil.ReadFile(path.Join(f.Path, "actions.json"))
	if err != nil {
		return a, err
	}

	err = json.Unmarshal(bs, &a)
	if err != nil {
		return a, err
	}

	if f.debug {
		f.logger.Println(fmt.Sprintf("LoadActions() -> Loaded actions on %d games", len(*a)))
	}
	return a, nil
}

func (f FileSystem) LoadSource() (*Source, error) {
	var s *Source

	bs, err := ioutil.ReadFile(path.Join(f.Path, "source.json"))
	if err != nil {
		return s, err
	}

	err = json.Unmarshal(bs, &s)
	if err != nil {
		return s, err
	}

	if f.debug {
		f.logger.Println(fmt.Sprintf("LoadSource() -> Loaded source of type %s", (*s).Type))
	}
	return s, nil
}

func (f FileSystem) RemoveActions() error {
	has, err := f.HasActions()
	if err != nil {
		return err
	}

	if has {
	    err = os.Remove(path.Join(f.Path, "actions.json"))
	}
	if err == nil && f.debug {
		f.logger.Println("RemoveActions(...) -> Removed actions file")
	}
	return err
}

func (f FileSystem) RemoveSource() error {
	has, err := f.HasSource()
	if err != nil {
		return err
	}

	if has {
	    err = os.Remove(path.Join(f.Path, "source.json"))
	}
	if err == nil && f.debug {
		f.logger.Println("RemoveSource(...) -> Removed source file")
	}
	return err
}

func (f FileSystem) AddGame(gameId int) error {
	gameDir := path.Join(f.Path, strconv.Itoa(gameId))
	instDir := path.Join(gameDir, "installers")
	extrDir := path.Join(gameDir, "extras")

	_, err := os.Stat(instDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(instDir, 0755)
			if err != nil {
				msg := fmt.Sprintf("AddGame(gameId=%d) -> Error occured while creating installers directory exists: %s", gameId, err.Error())
				return errors.New(msg)
			}
		} else {
			msg := fmt.Sprintf("AddGame(gameId=%d) -> Error occured while checking if installers directory exists: %s", gameId, err.Error())
			return errors.New(msg)
		}
	}

	_, err = os.Stat(extrDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(extrDir, 0755)
			if err != nil {
				msg := fmt.Sprintf("AddGame(gameId=%d) -> Error occured while creating extras directory exists: %s", gameId, err.Error())
				return errors.New(msg)
			}
		} else {
			msg := fmt.Sprintf("AddGame(gameId=%d) -> Error occured while checking if extras directory exists: %s", gameId, err.Error())
			return errors.New(msg)
		}
	}

	
	if f.debug {
		f.logger.Println(fmt.Sprintf("AddGame(gameId=%d) -> Created game directory", gameId))
	}
	return nil
}

func (f FileSystem) RemoveGame(gameId int) error {
	gameDir := path.Join(f.Path, strconv.Itoa(gameId))

	_, err := os.Stat(gameDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	}

	err = os.RemoveAll(gameDir)

	if err == nil && f.debug {
		f.logger.Println(fmt.Sprintf("RemoveGame(gameId=%d) -> Removed game directory", gameId))
	}
	return err
}

func (f FileSystem) UploadFile(source io.ReadCloser, gameId int, kind string, name string, expectedSize int64) (string, error) {
	var fPath string
	if kind == "installer" {
		fPath = path.Join(f.Path, strconv.Itoa(gameId), "installers", name)
	} else if kind == "extra" {
		fPath = path.Join(f.Path, strconv.Itoa(gameId), "extras", name)
	} else {
		return "", errors.New("Unknown kind of file")
	}

	h := md5.New()

	dest, err := os.Create(fPath)
	if err != nil {
		return "", err
	}

	w := io.MultiWriter(dest, h)
	io.Copy(w, source)
	dest.Close()

	info, infoErr := os.Stat(fPath)
	if infoErr != nil {
		return "", infoErr
	} else if info.Size() != expectedSize {
		msg := fmt.Sprintf("Created file at %d has size %d which doesn't match expected size %s", fPath, info.Size(), expectedSize)
		return "", errors.New(msg)
	}

	if f.debug {
		f.logger.Println(fmt.Sprintf("UploadFile(source=..., gameId=%d, kind=%s, name=%s) -> Uploaded file", gameId, kind, name))
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func (f FileSystem) RemoveFile(gameId int, kind string, name string) error {
	var fPath string
	if kind == "installer" {
		fPath = path.Join(f.Path, strconv.Itoa(gameId), "installers", name)
	} else if kind == "extra" {
		fPath = path.Join(f.Path, strconv.Itoa(gameId), "extras", name)
	} else {
		return errors.New("Unknown kind of file")
	}

	err := os.Remove(fPath)

	if err == nil && f.debug {
		f.logger.Println(fmt.Sprintf("RemoveFile(gameId=%d, kind=%s, name=%s) -> Removed file", gameId, kind, name))
	}
	return err
}

func (f FileSystem) DownloadFile(gameId int, kind string, name string) (io.ReadCloser, int64, error) {
	var fPath string
	if kind == "installer" {
		fPath = path.Join(f.Path, strconv.Itoa(gameId), "installers", name)
	} else if kind == "extra" {
		fPath = path.Join(f.Path, strconv.Itoa(gameId), "extras", name)
	} else {
		msg := fmt.Sprintf("DownloadFile(gameId=%d, kind=%s, name=%s) -> Unknown kind of file", gameId, kind, name)
		return nil, 0, errors.New(msg)
	}

	fi, err := os.Stat(fPath)
	if err != nil {
		msg := fmt.Sprintf("DownloadFile(gameId=%d, kind=%s, name=%s) -> Error occured while retrieving file size: %s", gameId, kind, name, err.Error())
		return nil, 0, errors.New(msg)
	}
	size := fi.Size()

	downloadHandle, openErr := os.Open(fPath)
	if openErr != nil {
		msg := fmt.Sprintf("DownloadFile(gameId=%d, kind=%s, name=%s) -> Error occured while opening file for download: %s", gameId, kind, name, openErr.Error())
		return nil, 0, errors.New(msg)
	}

	if f.debug {
		f.logger.Println(fmt.Sprintf("DownloadFile(gameId=%d, kind=%s, name=%s) -> Fetched file download handle", gameId, kind, name))
	}
	return downloadHandle, size, nil
}

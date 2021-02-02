package storage

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
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

func GetFileSystem(path string, debug bool) FileSystem {
	return FileSystem{path, debug, log.New(os.Stdout, "FS STORE: ", log.Lshortfile)}
}

func (f FileSystem) HasManifest() (bool, error) {
	_, err := os.Stat(path.Join(f.Path, "manifest.json"))
	if err != nil {
		if os.IsNotExist(err) {
			if f.debug {
				f.logger.Println("HasManifest() -> Manifest not found")
			}
			return false, nil
		} else {
			return true, err
		}
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
		} else {
			return true, err
		}
	}
	if f.debug {
		f.logger.Println("HasActions() -> Actions found")
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
	if err == nil {
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
	if err == nil {
		f.logger.Println(fmt.Sprintf("StoreActions(...) -> Stored actions on %d games", len(*a)))
	}
	return err
}

func (f FileSystem) LoadManifest() (*manifest.Manifest, error) {
	var m *manifest.Manifest

	bs, err := ioutil.ReadFile(path.Join(f.Path, "manifest.json"))
	if err != nil {
		return m, err
	}

	err = json.Unmarshal(bs, m)
	if err != nil {
		return m, err
	}

	return m, nil
}

func (f FileSystem) LoadActions() (*manifest.GameActions, error) {
	var a *manifest.GameActions

	bs, err := ioutil.ReadFile(path.Join(f.Path, "actions.json"))
	if err != nil {
		return a, err
	}

	err = json.Unmarshal(bs, a)
	if err != nil {
		return a, err
	}

	return a, nil
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
				return err
			}
		} else {
			return err
		}
	}

	_, err = os.Stat(extrDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(instDir, 0755)
			if err != nil {
				return err
			}
		} else {
			return err
		}
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
	if err != nil {
		return err
	}

	return nil
}

func (f FileSystem) UploadFile(source io.ReadCloser, gameId int, kind string, name string) ([]byte, error) {
	var fPath string
	if kind == "installer" {
		fPath = path.Join(f.Path, strconv.Itoa(gameId), "installers", name)
	} else if kind == "extra" {
		fPath = path.Join(f.Path, strconv.Itoa(gameId), "extras", name)
	} else {
		return nil, errors.New("Unknown kind of file")
	}

	h := md5.New()

	dest, err := os.Create(fPath)
	if err != nil {
		return nil, err
	}
	defer dest.Close()

	w := io.MultiWriter(dest, h)
	io.Copy(w, source)

	return h.Sum(nil), nil
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
	return err
}

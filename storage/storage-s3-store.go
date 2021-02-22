package storage

import (
	"bytes"
	"context"
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

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Configs struct {
	Endpoint  string
	Region    string
	Bucket    string
	Tls       bool
	AccessKey string
	SecretKey string
}

type S3Store struct {
	client *minio.Client
	configs *S3Configs
	debug bool
	logger *log.Logger
}

func GetS3StoreFromConfigFile(path string, debug bool, tag string) (S3Store, error) {
	var configs S3Configs

	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return S3Store{nil, nil, false, nil}, err
	}

	err = json.Unmarshal(bs, &configs)
	if err != nil {
		return S3Store{nil, nil, false, nil}, err
	}

	return getS3Store(&configs, debug, tag)
}

func GetS3StoreFromSource(s Source, debug bool, tag string) (S3Store, error) {
	if s.Type != "s3" {
		msg := fmt.Sprintf("Cannot load S3 store from source of type %s", s.Type)
		return S3Store{nil, nil, false, nil}, errors.New(msg)
	}
	return getS3Store(&(s.S3Params), debug, tag)
}

func getS3Store(configs *S3Configs, debug bool, tag string) (S3Store, error) {
	var logPrefix string
	if tag == "" {
		logPrefix = "FS: "
	} else {
		logPrefix = fmt.Sprintf("FS-%s: ", tag)
	}

	client, err := minio.New((*configs).Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4((*configs).AccessKey, (*configs).SecretKey, ""),
		Secure: (*configs).Tls,
		Region: (*configs).Region,
	})

	if err != nil {
		msg := fmt.Sprintf("GetS3Store(endpoint=%s, ...) -> Error connecting to the s3 store: %s", (*configs).Endpoint, err.Error())
		return S3Store{nil, nil, false, nil}, errors.New(msg)
	}

	return S3Store{
		client: client,
		configs: configs,
		debug: debug,
		logger: log.New(os.Stdout, logPrefix, log.Lshortfile),
	}, nil
}

func (s S3Store) GenerateSource() *Source {
	src := Source{
		Type: "s3",
		S3Params: (*s.configs),
	}
	return &src
}

func (s S3Store) GetPrintableSummary() string {
	configs := *s.configs
	return fmt.Sprintf("S3Store{endpoint: %s, region: %s, bucket: %s}", configs.Endpoint, configs.Region, configs.Bucket)
}

func (s S3Store) Exists() (bool, error) {
	configs := *s.configs
	found, existsErr := s.client.BucketExists(context.Background(), configs.Bucket)
	if existsErr != nil {
		msg := fmt.Sprintf("Exists() -> Error occured while trying to ascertain the bucket's existance: %s", existsErr.Error())
		return true, errors.New(msg)
	}

	if s.debug {
		if found {
			s.logger.Println("Exists() -> Bucket found")	
		} else {
			s.logger.Println("Exists() -> Bucket not found")			
		}
	}

	return found, nil
}

func (s S3Store) Initialize() error {
	configs := *s.configs
	makeErr := s.client.MakeBucket(context.Background(), configs.Bucket, minio.MakeBucketOptions{Region: configs.Region})
	if makeErr != nil {
		msg := fmt.Sprintf("Initialize() -> Error occured while trying to create bucket %s: %s", configs.Bucket, makeErr.Error())
		return errors.New(msg)
	}

	if s.debug {
		msg := fmt.Sprintf("Initialize() -> Bucket %s created", configs.Bucket)
		s.logger.Println(msg)
	}
	return nil
}

func (s S3Store) HasManifest() (bool, error) {
	configs := *s.configs
	_, err := s.client.StatObject(context.Background(), configs.Bucket, "manifest.json", minio.StatObjectOptions{})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			if s.debug {
				s.logger.Println("HasManifest() -> Manifest not found")
			}
			return false, nil
		}

		msg := fmt.Sprintf("HasManifest() -> The following error occured while ascertaining manifest's existance: %s", err.Error())
		return true, errors.New(msg)
	}

	if s.debug {
		s.logger.Println("HasManifest() -> Manifest found")
	}

	return true, nil
}

func (s S3Store) HasActions() (bool, error) {
	configs := *s.configs
	_, err := s.client.StatObject(context.Background(), configs.Bucket, "actions.json", minio.StatObjectOptions{})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			if s.debug {
				s.logger.Println("HasActions() -> Actions not found")
			}
			return false, nil
		}
		
		msg := fmt.Sprintf("HasActions() -> The following error occured while ascertaining actions' existance: %s", err.Error())
		return true, errors.New(msg)
	}

	if s.debug {
		s.logger.Println("HasActions() -> Actions found")
	}

	return true, nil
}

func (s S3Store) HasSource() (bool, error) {
	configs := *s.configs
	_, err := s.client.StatObject(context.Background(), configs.Bucket, "source.json", minio.StatObjectOptions{})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			if s.debug {
				s.logger.Println("HasSource() -> Source not found")
			}
			return false, nil
		}
		
		msg := fmt.Sprintf("HasSource() -> The following error occured while ascertaining source's existance: %s", err.Error())
		return true, errors.New(msg)
	}

	if s.debug {
		s.logger.Println("HasSource() -> Source found")
	}

	return true, nil
}

func (s S3Store) StoreManifest(m *manifest.Manifest) error {
	var err error
	var buf bytes.Buffer
	var output []byte
	configs := *s.configs
	
	output, err = json.Marshal(*m)

	if err != nil {
		return err
	}

	json.Indent(&buf, output, "", "  ")
	output = buf.Bytes()

	_, err = s.client.PutObject(context.Background(), configs.Bucket, "manifest.json", bytes.NewReader(output), int64(len(output)), minio.PutObjectOptions{ContentType:"application/json"})
	if err == nil && s.debug {
		s.logger.Println(fmt.Sprintf("StoreManifest(...) -> Stored manifest with %d games", len((*m).Games)))
	}
	return err
}

func (s S3Store) StoreActions(a *manifest.GameActions) error {
	var err error
	var buf bytes.Buffer
	var output []byte
	configs := *s.configs

	output, err = json.Marshal(*a)

	if err != nil {
		return err
	}

	json.Indent(&buf, output, "", "  ")
	output = buf.Bytes()

	_, err = s.client.PutObject(context.Background(), configs.Bucket, "actions.json", bytes.NewReader(output), int64(len(output)), minio.PutObjectOptions{ContentType:"application/json"})
	if err == nil && s.debug {
		s.logger.Println(fmt.Sprintf("StoreActions(...) -> Stored actions on %d games", len(*a)))
	}
	return err
}

func (s S3Store) StoreSource(o *Source) error {
	var err error
	var buf bytes.Buffer
	var output []byte
	configs := *s.configs

	output, err = json.Marshal(*o)

	if err != nil {
		return err
	}

	json.Indent(&buf, output, "", "  ")
	output = buf.Bytes()

	_, err = s.client.PutObject(context.Background(), configs.Bucket, "source.json", bytes.NewReader(output), int64(len(output)), minio.PutObjectOptions{ContentType:"application/json"})
	if err == nil && s.debug {
		s.logger.Println(fmt.Sprintf("StoreSource(...) -> Stored source of type %s", o.Type))
	}
	return err
}

func (s S3Store) LoadManifest() (*manifest.Manifest, error) {
	var m manifest.Manifest
	configs := *s.configs

	objPtr, err := s.client.GetObject(context.Background(), configs.Bucket, "manifest.json", minio.GetObjectOptions{})
	if err != nil {
		return &m, err
	}

	bs, bErr := ioutil.ReadAll(objPtr)
	if bErr != nil {
		return &m, bErr
	}

	err = json.Unmarshal(bs, &m)
	if err != nil {
		return &m, err
	}

	if s.debug {
		s.logger.Println(fmt.Sprintf("LoadManifest() -> Loaded manifest with %d games", len(m.Games)))
	}
	return &m, nil
}

func (s S3Store) LoadActions() (*manifest.GameActions, error) {
	var a *manifest.GameActions
	configs := *s.configs

	objPtr, err := s.client.GetObject(context.Background(), configs.Bucket, "actions.json", minio.GetObjectOptions{})
	if err != nil {
		return a, err
	}

	bs, bErr := ioutil.ReadAll(objPtr)
	if bErr != nil {
		return a, bErr
	}

	err = json.Unmarshal(bs, &a)
	if err != nil {
		return a, err
	}

	if s.debug {
		s.logger.Println(fmt.Sprintf("LoadActions() -> Loaded actions on %d games", len(*a)))
	}
	return a, nil
}

func (s S3Store) LoadSource() (*Source, error) {
	var o *Source
	configs := *s.configs

	objPtr, err := s.client.GetObject(context.Background(), configs.Bucket, "source.json", minio.GetObjectOptions{})
	if err != nil {
		return o, err
	}

	bs, bErr := ioutil.ReadAll(objPtr)
	if bErr != nil {
		return o, bErr
	}

	err = json.Unmarshal(bs, &o)
	if err != nil {
		return o, err
	}

	if s.debug {
		s.logger.Println(fmt.Sprintf("LoadSource() -> Loaded source of type %s", (*o).Type))
	}
	return o, nil
}

func (s S3Store) RemoveActions() error {
	configs := *s.configs

	has, err := s.HasActions()
	if err != nil {
		return err
	}

	if has {
		err = s.client.RemoveObject(context.Background(), configs.Bucket, "actions.json", minio.RemoveObjectOptions{})
	}
	if err == nil && s.debug {
		s.logger.Println("RemoveActions(...) -> Removed actions file")
	}
	return err
}

func (s S3Store) RemoveSource() error {
	configs := *s.configs

	has, err := s.HasSource()
	if err != nil {
		return err
	}

	if has {
		err = s.client.RemoveObject(context.Background(), configs.Bucket, "source.json", minio.RemoveObjectOptions{})
	}
	if err == nil && s.debug {
		s.logger.Println("RemoveSource(...) -> Removed source file")
	}
	return err
}

func (s S3Store) AddGame(gameId int) error {
	if s.debug {
		s.logger.Println(fmt.Sprintf("AddGame(gameId=%d) -> No-op as s3 store doesn't have a real directory structure", gameId))
	}

	return nil
}

func (s S3Store) RemoveGame(gameId int) error {
	if s.debug {
		s.logger.Println(fmt.Sprintf("RemoveGame(gameId=%d) -> No-op as s3 store doesn't have a real directory structure", gameId))
	}

	return nil
}

func (s S3Store) UploadFile(source io.ReadCloser, gameId int, kind string, name string, expectedSize int64) (string, error) {
	configs := *s.configs

	var fPath string
	if kind == "installer" {
		fPath = path.Join(strconv.Itoa(gameId), "installers", name)
	} else if kind == "extra" {
		fPath = path.Join(strconv.Itoa(gameId), "extras", name)
	} else {
		return "", errors.New("Unknown kind of file")
	}

	_, err := s.client.PutObject(context.Background(), configs.Bucket, fPath, source, expectedSize, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}

	info, statErr := s.client.StatObject(context.Background(), configs.Bucket, fPath, minio.StatObjectOptions{})
	if statErr != nil {
		return "", statErr
	}

	if info.Size != expectedSize {
		msg := fmt.Sprintf("Object %s has a size of %d which doesn't match expected size of %d", fPath, info.Size, expectedSize)
		return "", errors.New(msg)
	}

	if s.debug {
		s.logger.Println(fmt.Sprintf("UploadFile(source=..., gameId=%d, kind=%s, name=%s) -> Uploaded file", gameId, kind, name))
	}
	return info.ETag, nil
}

func (s S3Store) RemoveFile(gameId int, kind string, name string) error {
	configs := *s.configs
	
	var oPath string
	if kind == "installer" {
		oPath = path.Join(strconv.Itoa(gameId), "installers", name)
	} else if kind == "extra" {
		oPath = path.Join(strconv.Itoa(gameId), "extras", name)
	} else {
		return errors.New("Unknown kind of file")
	}

	err := s.client.RemoveObject(context.Background(), configs.Bucket, oPath, minio.RemoveObjectOptions{})

	if err == nil && s.debug {
		s.logger.Println(fmt.Sprintf("RemoveFile(gameId=%d, kind=%s, name=%s) -> Removed file", gameId, kind, name))
	}
	return err
}

func (s S3Store) DownloadFile(gameId int, kind string, name string) (io.ReadCloser, int64, error) {
	configs := *s.configs

	var fPath string
	if kind == "installer" {
		fPath = path.Join(strconv.Itoa(gameId), "installers", name)
	} else if kind == "extra" {
		fPath = path.Join(strconv.Itoa(gameId), "extras", name)
	} else {
		msg := fmt.Sprintf("DownloadFile(gameId=%d, kind=%s, name=%s) -> Unknown kind of file", gameId, kind, name)
		return nil, 0, errors.New(msg)
	}

	fi, err := s.client.StatObject(context.Background(), configs.Bucket, fPath, minio.StatObjectOptions{})
	if err != nil {
		msg := fmt.Sprintf("DownloadFile(gameId=%d, kind=%s, name=%s) -> Error occured while retrieving file size: %s", gameId, kind, name, err.Error())
		return nil, 0, errors.New(msg)
	}
	size := fi.Size

	downloadHandle, openErr := s.client.GetObject(context.Background(), configs.Bucket, fPath, minio.GetObjectOptions{})
	if openErr != nil {
		msg := fmt.Sprintf("DownloadFile(gameId=%d, kind=%s, name=%s) -> Error occured while opening file for download: %s", gameId, kind, name, openErr.Error())
		return nil, 0, errors.New(msg)
	}

	if s.debug {
		s.logger.Println(fmt.Sprintf("DownloadFile(gameId=%d, kind=%s, name=%s) -> Fetched file download handle", gameId, kind, name))
	}
	return downloadHandle, size, nil
}

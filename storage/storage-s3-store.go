package storage

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"gogcli/logging"
	"gogcli/manifest"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

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
	client  *minio.Client
	configs *S3Configs
	logger  *logging.Logger
}

func GetS3StoreFromConfigFile(path string, logSource *logging.Source, tag string) (S3Store, error) {
	var configs S3Configs

	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return S3Store{nil, nil, nil}, err
	}

	err = json.Unmarshal(bs, &configs)
	if err != nil {
		return S3Store{nil, nil, nil}, err
	}

	return getS3Store(&configs, logSource, tag)
}

func GetS3StoreFromSource(s Source, logSource *logging.Source, tag string) (S3Store, error) {
	if s.Type != "s3" {
		msg := fmt.Sprintf("Cannot load S3 store from source of type %s", s.Type)
		return S3Store{nil, nil, nil}, errors.New(msg)
	}
	return getS3Store(&(s.S3Params), logSource, tag)
}

func getS3Store(configs *S3Configs, logSource *logging.Source, tag string) (S3Store, error) {
	var logPrefix string
	if tag == "" {
		logPrefix = "[s3] "
	} else {
		logPrefix = fmt.Sprintf("[s3-%s] ", tag)
	}

	client, err := minio.New((*configs).Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4((*configs).AccessKey, (*configs).SecretKey, ""),
		Secure: (*configs).Tls,
		Region: (*configs).Region,
	})

	if err != nil {
		msg := fmt.Sprintf("GetS3Store(endpoint=%s, ...) -> Error connecting to the s3 store: %s", (*configs).Endpoint, err.Error())
		return S3Store{nil, nil, nil}, errors.New(msg)
	}

	return S3Store{
		client:  client,
		configs: configs,
		logger:  logSource.CreateLogger(os.Stdout, logPrefix, log.Lmsgprefix),
	}, nil
}

func (s S3Store) GetListing() (*StorageListing, error) {
	listing := NewEmptyStorageListing(S3StoreDownloader{s})
	configs := *s.configs
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	gameFileRegex := regexp.MustCompile(`^(?P<id>\d+)/(?P<kind>(?:installers)|(?:extras))/(?P<file>.+)$`)

	objChan := s.client.ListObjects(ctx, configs.Bucket, minio.ListObjectsOptions{
		Recursive: true,
	})
	for obj := range objChan {
		if obj.Err != nil {
			return nil, obj.Err
		}
		if gameFileRegex.MatchString(obj.Key) {
			match := gameFileRegex.FindStringSubmatch(obj.Key)
			gameId, _ := strconv.ParseInt(match[1], 10, 64)
			gameListing, ok := listing.Games[gameId]
			if !ok {
				gameListing = StorageListingGame{
					Game:       manifest.GameInfo{Id: gameId},
					Installers: make([]manifest.FileInfo, 0),
					Extras:     make([]manifest.FileInfo, 0),
				}
			}
			if match[2] == "installers" {
				fileInfo := manifest.FileInfo{Game: listing.Games[gameId].Game, Name: match[3], Kind: "installer"}
				gameListing.Installers = append(gameListing.Installers, fileInfo)
			} else {
				fileInfo := manifest.FileInfo{Game: listing.Games[gameId].Game, Name: match[3], Kind: "extra"}
				gameListing.Extras = append(gameListing.Extras, fileInfo)
			}
			listing.Games[gameId] = gameListing
		}
	}
	return &listing, nil
}

func (s S3Store) SupportsReaderAt() bool {
	return true
}

func (s S3Store) GenerateSource() *Source {
	src := Source{
		Type:     "s3",
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

	if found {
		s.logger.Debug("Exists() -> Bucket found")
	} else {
		s.logger.Debug("Exists() -> Bucket not found")
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

	msg := fmt.Sprintf("Initialize() -> Bucket %s created", configs.Bucket)
	s.logger.Debug(msg)
	return nil
}

func (s S3Store) HasManifest() (bool, error) {
	configs := *s.configs
	_, err := s.client.StatObject(context.Background(), configs.Bucket, "manifest.json", minio.StatObjectOptions{})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			s.logger.Debug("HasManifest() -> Manifest not found")
			return false, nil
		}

		msg := fmt.Sprintf("HasManifest() -> The following error occured while ascertaining manifest's existance: %s", err.Error())
		return true, errors.New(msg)
	}

	s.logger.Debug("HasManifest() -> Manifest found")
	return true, nil
}

func (s S3Store) HasActions() (bool, error) {
	configs := *s.configs
	_, err := s.client.StatObject(context.Background(), configs.Bucket, "actions.json", minio.StatObjectOptions{})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			s.logger.Debug("HasActions() -> Actions not found")
			return false, nil
		}

		msg := fmt.Sprintf("HasActions() -> The following error occured while ascertaining actions' existance: %s", err.Error())
		return true, errors.New(msg)
	}

	s.logger.Debug("HasActions() -> Actions found")
	return true, nil
}

func (s S3Store) HasSource() (bool, error) {
	configs := *s.configs
	_, err := s.client.StatObject(context.Background(), configs.Bucket, "source.json", minio.StatObjectOptions{})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			s.logger.Debug("HasSource() -> Source not found")
			return false, nil
		}

		msg := fmt.Sprintf("HasSource() -> The following error occured while ascertaining source's existance: %s", err.Error())
		return true, errors.New(msg)
	}

	s.logger.Debug("HasSource() -> Source found")
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

	_, err = s.client.PutObject(context.Background(), configs.Bucket, "manifest.json", bytes.NewReader(output), int64(len(output)), minio.PutObjectOptions{ContentType: "application/json"})
	if err == nil {
		s.logger.Debug(fmt.Sprintf("StoreManifest(...) -> Stored manifest with %d games", len((*m).Games)))
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

	_, err = s.client.PutObject(context.Background(), configs.Bucket, "actions.json", bytes.NewReader(output), int64(len(output)), minio.PutObjectOptions{ContentType: "application/json"})
	if err == nil {
		s.logger.Debug(fmt.Sprintf("StoreActions(...) -> Stored actions on %d games", len(*a)))
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

	_, err = s.client.PutObject(context.Background(), configs.Bucket, "source.json", bytes.NewReader(output), int64(len(output)), minio.PutObjectOptions{ContentType: "application/json"})
	if err == nil {
		s.logger.Debug(fmt.Sprintf("StoreSource(...) -> Stored source of type %s", o.Type))
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

	s.logger.Debug(fmt.Sprintf("LoadManifest() -> Loaded manifest with %d games", len(m.Games)))
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

	s.logger.Debug(fmt.Sprintf("LoadActions() -> Loaded actions on %d games", len(*a)))
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

	s.logger.Debug(fmt.Sprintf("LoadSource() -> Loaded source of type %s", (*o).Type))
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
	if err == nil {
		s.logger.Debug("RemoveActions(...) -> Removed actions file")
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
	if err == nil {
		s.logger.Debug("RemoveSource(...) -> Removed source file")
	}
	return err
}

func (s S3Store) AddGame(game manifest.GameInfo) error {
	s.logger.Debug(fmt.Sprintf("AddGame(game={Id=%d, ...}) -> No-op as s3 store doesn't have a real directory structure", game.Id))
	return nil
}

func (s S3Store) RemoveGame(game manifest.GameInfo) error {
	s.logger.Debug(fmt.Sprintf("RemoveGame(game={Id=%d, ...}) -> No-op as s3 store doesn't have a real directory structure", game.Id))
	return nil
}

func (s S3Store) UploadFile(source io.ReadCloser, file manifest.FileInfo) (string, error) {
	configs := *s.configs

	var fPath string
	if file.Kind == "installer" {
		arr := []string{strconv.FormatInt(file.Game.Id, 10), "installers", file.Name}
		fPath = strings.Join(arr, "/")
	} else if file.Kind == "extra" {
		arr := []string{strconv.FormatInt(file.Game.Id, 10), "extras", file.Name}
		fPath = strings.Join(arr, "/")
	} else {
		return "", errors.New("Unknown kind of file")
	}

	_, err := s.client.PutObject(context.Background(), configs.Bucket, fPath, source, file.Size, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}

	downloadHandle, size, downErr := s.DownloadFile(file)
	if downErr != nil {
		return "", downErr
	}
	defer downloadHandle.Close()
	h := md5.New()
	io.Copy(h, downloadHandle)
	checksum := hex.EncodeToString(h.Sum(nil))

	if size != file.Size {
		msg := fmt.Sprintf("Object %s has a size of %d which doesn't match expected size of %d", fPath, size, file.Size)
		return "", errors.New(msg)
	}

	s.logger.Debug(fmt.Sprintf("UploadFile(source=..., gameId=%d, kind=%s, name=%s) -> Uploaded file", file.Game.Id, file.Kind, file.Name))
	return checksum, nil
}

func (s S3Store) RemoveFile(file manifest.FileInfo) error {
	configs := *s.configs

	var oPath string
	if file.Kind == "installer" {
		arr := []string{strconv.FormatInt(file.Game.Id, 10), "installers", file.Name}
		oPath = strings.Join(arr, "/")
	} else if file.Kind == "extra" {
		arr := []string{strconv.FormatInt(file.Game.Id, 10), "extras", file.Name}
		oPath = strings.Join(arr, "/")
	} else {
		return errors.New("Unknown kind of file")
	}

	_, err := s.client.StatObject(context.Background(), configs.Bucket, oPath, minio.StatObjectOptions{})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code != "NoSuchKey" {
			return err
		}
	} else {
		err := s.client.RemoveObject(context.Background(), configs.Bucket, oPath, minio.RemoveObjectOptions{})
		if err != nil {
			return err
		}
	}

	s.logger.Debug(fmt.Sprintf("RemoveFile(gameId=%d, kind=%s, name=%s) -> Removed file", file.Game.Id, file.Kind, file.Name))
	return nil
}

func (s S3Store) DownloadFile(file manifest.FileInfo) (io.ReadCloser, int64, error) {
	configs := *s.configs

	var fPath string
	if file.Kind == "installer" {
		arr := []string{strconv.FormatInt(file.Game.Id, 10), "installers", file.Name}
		fPath = strings.Join(arr, "/")
	} else if file.Kind == "extra" {
		arr := []string{strconv.FormatInt(file.Game.Id, 10), "extras", file.Name}
		fPath = strings.Join(arr, "/")
	} else {
		msg := fmt.Sprintf("DownloadFile(gameId=%d, kind=%s, name=%s) -> Unknown kind of file", file.Game.Id, file.Kind, file.Name)
		return nil, 0, errors.New(msg)
	}

	fi, err := s.client.StatObject(context.Background(), configs.Bucket, fPath, minio.StatObjectOptions{})
	if err != nil {
		msg := fmt.Sprintf("DownloadFile(gameId=%d, kind=%s, name=%s) -> Error occured while retrieving file size: %s", file.Game.Id, file.Kind, file.Name, err.Error())
		return nil, 0, errors.New(msg)
	}
	size := fi.Size

	downloadHandle, openErr := s.client.GetObject(context.Background(), configs.Bucket, fPath, minio.GetObjectOptions{})
	if openErr != nil {
		msg := fmt.Sprintf("DownloadFile(gameId=%d, kind=%s, name=%s) -> Error occured while opening file for download: %s", file.Game.Id, file.Kind, file.Name, openErr.Error())
		return nil, 0, errors.New(msg)
	}

	s.logger.Debug(fmt.Sprintf("DownloadFile(gameId=%d, kind=%s, name=%s) -> Fetched file download handle", file.Game.Id, file.Kind, file.Name))
	return downloadHandle, size, nil
}

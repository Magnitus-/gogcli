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

type S3Store struct {
	client *minio.Client
	endpoint string
	region string
	bucket string
	debug bool
	logger *log.Logger
}

type S3Configs struct {
	Endpoint  string
	Region    string
	Bucket    string
	Tls       bool
	AccessKey string
	SecretKey string
}

func GetS3StoreFromConfigFile(path string, debug bool, tag string) (S3Store, error) {
	var configs S3Configs

	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return S3Store{nil, "", "", "", false, nil}, err
	}

	err = json.Unmarshal(bs, &configs)
	if err != nil {
		return S3Store{nil, "", "", "", false, nil}, err
	}

	return GetS3Store(configs.Endpoint, configs.Region, configs.Bucket, configs.AccessKey, configs.SecretKey, configs.Tls, debug, tag)
}

func GetS3Store(endpoint string, region string, bucket string, accessKey string, secretKey string, tls bool, debug bool, tag string) (S3Store, error) {
	var logPrefix string
	if tag == "" {
		logPrefix = "FS: "
	} else {
		logPrefix = fmt.Sprintf("FS-%s: ", tag)
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: tls,
		Region: region,
	})

	if err != nil {
		msg := fmt.Sprintf("GetS3Store(endpoint=%s, ...) -> Error connecting to the s3 store: %s", endpoint, err.Error())
		return S3Store{nil, "", "", "", false, nil}, errors.New(msg)
	}

	return S3Store{
		client: client,
		endpoint: endpoint,
		region: region,
		bucket: bucket,
		debug: debug,
		logger: log.New(os.Stdout, logPrefix, log.Lshortfile),
	}, nil
}

func (s S3Store) GetPrintableSummary() string {
	return fmt.Sprintf("S3Store{endpoint: %s, region: %s, bucket: %s}", s.endpoint, s.region, s.bucket)
}

func (s S3Store) Exists() (bool, error) {
	found, existsErr := s.client.BucketExists(context.Background(), s.bucket)
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
	makeErr := s.client.MakeBucket(context.Background(), s.bucket, minio.MakeBucketOptions{Region: s.region})
	if makeErr != nil {
		msg := fmt.Sprintf("Initialize() -> Error occured while trying to create bucket %s: %s", s.bucket, makeErr.Error())
		return errors.New(msg)
	}

	if s.debug {
		msg := fmt.Sprintf("Initialize() -> Bucket %s created", s.bucket)
		s.logger.Println(msg)
	}
	return nil
}

func (s S3Store) HasManifest() (bool, error) {
	_, err := s.client.StatObject(context.Background(), s.bucket, "manifest.json", minio.StatObjectOptions{})
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
	_, err := s.client.StatObject(context.Background(), s.bucket, "actions.json", minio.StatObjectOptions{})
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

func (s S3Store) StoreManifest(m *manifest.Manifest) error {
	var err error
	var buf bytes.Buffer
	var output []byte

	output, err = json.Marshal(*m)

	if err != nil {
		return err
	}

	json.Indent(&buf, output, "", "  ")
	output = buf.Bytes()

	_, err = s.client.PutObject(context.Background(), s.bucket, "manifest.json", bytes.NewReader(output), int64(len(output)), minio.PutObjectOptions{ContentType:"application/json"})
	if err == nil && s.debug {
		s.logger.Println(fmt.Sprintf("StoreManifest(...) -> Stored manifest with %d games", len((*m).Games)))
	}
	return err
}

func (s S3Store) StoreActions(a *manifest.GameActions) error {
	var err error
	var buf bytes.Buffer
	var output []byte

	output, err = json.Marshal(*a)

	if err != nil {
		return err
	}

	json.Indent(&buf, output, "", "  ")
	output = buf.Bytes()

	_, err = s.client.PutObject(context.Background(), s.bucket, "actions.json", bytes.NewReader(output), int64(len(output)), minio.PutObjectOptions{ContentType:"application/json"})
	if err == nil && s.debug {
		s.logger.Println(fmt.Sprintf("StoreActions(...) -> Stored actions on %d games", len(*a)))
	}
	return err
}

func (s S3Store) LoadManifest() (*manifest.Manifest, error) {
	var m manifest.Manifest

	objPtr, err := s.client.GetObject(context.Background(), s.bucket, "manifest.json", minio.GetObjectOptions{})
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

	objPtr, err := s.client.GetObject(context.Background(), s.bucket, "actions.json", minio.GetObjectOptions{})
	if err != nil {
		return a, err
	}

	bs, bErr := ioutil.ReadAll(objPtr)
	if bErr != nil {
		return a, bErr
	}

	err = json.Unmarshal(bs, a)
	if err != nil {
		return a, err
	}

	if s.debug {
		s.logger.Println(fmt.Sprintf("LoadActions() -> Loaded actions on %d games", len(*a)))
	}
	return a, nil
}

func (s S3Store) RemoveActions() error {
	has, err := s.HasActions()
	if err != nil {
		return err
	}

	if has {
		err = s.client.RemoveObject(context.Background(), s.bucket, "actions.json", minio.RemoveObjectOptions{})
	}
	if err == nil && s.debug {
		s.logger.Println("RemoveActions(...) -> Removed actions file")
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
	var fPath string
	if kind == "installer" {
		fPath = path.Join(strconv.Itoa(gameId), "installers", name)
	} else if kind == "extra" {
		fPath = path.Join(strconv.Itoa(gameId), "extras", name)
	} else {
		return "", errors.New("Unknown kind of file")
	}

	_, err := s.client.PutObject(context.Background(), s.bucket, fPath, source, expectedSize, minio.PutObjectOptions{})
	if err != nil {
		return "", err
	}

	info, statErr := s.client.StatObject(context.Background(), s.bucket, fPath, minio.StatObjectOptions{})
	if statErr != nil {
		return "", statErr
	}

	if s.debug {
		s.logger.Println(fmt.Sprintf("UploadFile(source=..., gameId=%d, kind=%s, name=%s) -> Uploaded file", gameId, kind, name))
	}
	return info.ETag, nil
}

func (s S3Store) RemoveFile(gameId int, kind string, name string) error {
	var oPath string
	if kind == "installer" {
		oPath = path.Join(strconv.Itoa(gameId), "installers", name)
	} else if kind == "extra" {
		oPath = path.Join(strconv.Itoa(gameId), "extras", name)
	} else {
		return errors.New("Unknown kind of file")
	}

	err := s.client.RemoveObject(context.Background(), s.bucket, oPath, minio.RemoveObjectOptions{})

	if err == nil && s.debug {
		s.logger.Println(fmt.Sprintf("RemoveFile(gameId=%d, kind=%s, name=%s) -> Removed file", gameId, kind, name))
	}
	return err
}

func (s S3Store) DownloadFile(gameId int, kind string, name string) (io.ReadCloser, int64, error) {
	var fPath string
	if kind == "installer" {
		fPath = path.Join(strconv.Itoa(gameId), "installers", name)
	} else if kind == "extra" {
		fPath = path.Join(strconv.Itoa(gameId), "extras", name)
	} else {
		msg := fmt.Sprintf("DownloadFile(gameId=%d, kind=%s, name=%s) -> Unknown kind of file", gameId, kind, name)
		return nil, 0, errors.New(msg)
	}

	fi, err := s.client.StatObject(context.Background(), s.bucket, fPath, minio.StatObjectOptions{})
	if err != nil {
		msg := fmt.Sprintf("DownloadFile(gameId=%d, kind=%s, name=%s) -> Error occured while retrieving file size: %s", gameId, kind, name, err.Error())
		return nil, 0, errors.New(msg)
	}
	size := fi.Size

	downloadHandle, openErr := s.client.GetObject(context.Background(), s.bucket, fPath, minio.GetObjectOptions{})
	if openErr != nil {
		msg := fmt.Sprintf("DownloadFile(gameId=%d, kind=%s, name=%s) -> Error occured while opening file for download: %s", gameId, kind, name, openErr.Error())
		return nil, 0, errors.New(msg)
	}

	if s.debug {
		s.logger.Println(fmt.Sprintf("DownloadFile(gameId=%d, kind=%s, name=%s) -> Fetched file download handle", gameId, kind, name))
	}
	return downloadHandle, size, nil
}

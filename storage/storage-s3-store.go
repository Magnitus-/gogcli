package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Store struct {
	client *minio.Client
	bucket string
	debug bool
	logger *log.Logger
}

func GetS3Store(endpoint string, region string, bucket string, accessKey string, secretKey string, tls bool, debug bool, tag string) S3Store, error {
	var logPrefix string
	if tag == "" {
		logPrefix = "FS: "
	} else {
		logPrefix = fmt.Sprintf("FS-%s: ", tag)
	}

	client, err := minio.NewWithRegion(
		endpoint,
		accessKey,
		secretKey,
		tls,
		region,
	)

	if err != nil {
		msg := fmt.Sprintf("GetS3Store(endpoint=%s, ...) -> Error connecting to the s3 store: %s", err.Error())
		return nil, errors.New(msg)
	}

	found, existsErr := client.BucketExists(context.Background(), bucket)
	if existsErr != nil {
		msg := fmt.Sprintf("GetS3Store(endpoint=%s, ...) -> Error occured while trying to ascertain the bucket's existance: %s", existsErr.Error())
		return nil, errors.New(msg)
	}

	if !found {
		makeErr = client.MakeBucket(context.Background(), bucket, minio.MakeBucketOptions{Region: region})
		if makeErr != nil {
			msg := fmt.Sprintf("GetS3Store(endpoint=%s, ...) -> Error occured while trying to create missing bucket: %s", makeErr.Error())
			return nil, errors.New(msg)
		}
	}

	return S3Store{
		client: client,
		bucket: bucket,
		debug: debug,
		logger: log.New(os.Stdout, logPrefix, log.Lshortfile),
	}, nil
}

func (s S3Store) HasManifest() (bool, error) {
	_, err := s.client.StatObject(context.Background(), s.bucket, "manifest.json", minio.StatObjectOptions{})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchBucket" {
			if f.debug {
				f.logger.Println("HasManifest() -> Manifest not found")
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
		if errResponse.Code == "NoSuchBucket" {
			if f.debug {
				f.logger.Println("HasActions() -> Actions not found")
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

	_, err := s.client.PutObject(context.Background(), s.bucket, "manifest.json", bytes.NewReader(output), len(output), minio.PutObjectOptions{ContentType:"application/json"})
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

	_, err := s.client.PutObject(context.Background(), s.bucket, "actions.json", bytes.NewReader(output), len(output), minio.PutObjectOptions{ContentType:"application/json"})
	if err == nil && s.debug {
		s.logger.Println(fmt.Sprintf("StoreActions(...) -> Stored actions on %d games", len(*a)))
	}
	return err
}

func (s S3Store) LoadManifest() (*manifest.Manifest, error) {
	var m manifest.Manifest

	bs, err := s.client.GetObject(context.Background(), s.bucket, "manifest.json", minio.GetObjectOptions{})
	if err != nil {
		return a, err
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

	bs, err := s.client.GetObject(context.Background(), s.bucket, "actions.json", minio.GetObjectOptions{})
	if err != nil {
		return a, err
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
		err = minioClient.RemoveObject(context.Background(), s.bucket, "actions.json", minio.RemoveObjectOptions{})
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

func (s S3Store) UploadFile(source io.ReadCloser, gameId int, kind string, name string) (string, error) {
	//TODO
}

func (s S3Store) RemoveFile(gameId int, kind string, name string) error {
	var oPath string
	if kind == "installer" {
		oPath = path.Join(strconv.Itoa(gameId), "installers", name)
	} else if kind == "extra" {
		oPath = path.Join(strconv.Itoa(gameId), "extras", name)
	} else {
		return "", errors.New("Unknown kind of file")
	}

	err = minioClient.RemoveObject(context.Background(), s.bucket, oPath, minio.RemoveObjectOptions{})

	if err != nil && s.debug {
		s.logger.Println(fmt.Sprintf("RemoveFile(gameId=%d, kind=%s, name=%s) -> Removed file", gameId, kind, name))
	}
	return err
}

func (s S3Store) DownloadFile(gameId int, kind string, name string) (io.ReadCloser, int, error) {
	//TODO
}


package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gogcli/logging"
	"gogcli/manifest"
	"gogcli/metadata"
	"gogcli/storage"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

func loadActionsFromFile(path string) (manifest.GameActions, error) {
	var a manifest.GameActions
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return manifest.GameActions{}, err
	}

	err = json.Unmarshal(bs, &a)
	if err != nil {
		return manifest.GameActions{}, err
	}

	return a, nil
}

func loadManifestFromFile(path string) (manifest.Manifest, error) {
	var m manifest.Manifest
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return manifest.Manifest{}, err
	}

	err = json.Unmarshal(bs, &m)
	if err != nil {
		return manifest.Manifest{}, err
	}

	return m, nil
}

func loadMetadataFromFile(path string) (metadata.Metadata, error) {
	var m metadata.Metadata
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return metadata.Metadata{}, err
	}

	err = json.Unmarshal(bs, &m)
	if err != nil {
		return metadata.Metadata{}, err
	}

	return m, nil
}

func getStorage(path string, storageType string, logSource *logging.Source, loggerTag string) (storage.Storage, storage.Downloader) {
	if storageType != "fs" && storageType != "s3" {
		msg := fmt.Sprintf("Source storage type %s is invalid", storageType)
		fmt.Println(msg)
		os.Exit(1)
	}

	if storageType == "fs" {
		gameStorage := storage.GetFileSystem(path, logSource, loggerTag)
		downloader := storage.FileSystemDownloader{gameStorage}
		return gameStorage, downloader
	} else {
		gameStorage, err := storage.GetS3StoreFromConfigFile(path, logSource, loggerTag)
		processError(err)
		downloader := storage.S3StoreDownloader{gameStorage}
		return gameStorage, downloader
	}
}

func processErrors(errs []error) {
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
		os.Exit(1)
	}
}

func processError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type Errors struct {
	Errors []string
}

func deserializeErrors(file string) []error {
	var strErrs Errors
	errs := []error{}

	bs, err := ioutil.ReadFile(file)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Println("Could not load the error/warning file: ", err)
			os.Exit(1)
		}
		return errs
	}

	err = json.Unmarshal(bs, &strErrs)
	if err != nil {
		fmt.Println("Error/warning file doesn't appear to contain valid json: ", err)
		os.Exit(1)
	}

	for _, strErr := range strErrs.Errors {
		errs = append(errs, errors.New(strErr))
	}

	return errs
}

func processSerializableOutput(serializable interface{}, errs []error, terminal bool, file string) {
	hasErr := false
	var output []byte
	var e Errors

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if len(errs) > 0 {
		for _, err := range errs {
			e.Errors = append(e.Errors, err.Error())
		}
		_ = enc.Encode(e)
		hasErr = true
	} else {
		_ = enc.Encode(serializable)
	}

	output = buf.Bytes()

	if terminal {
		fmt.Println(string(output))
	} else {
		err := ioutil.WriteFile(file, output, 0644)
		if err != nil {
			fmt.Println(err)
			hasErr = true
		}
	}

	if hasErr {
		os.Exit(1)
	}
}

func PersistManifestProgress(file string) manifest.ManifestWriterStatePersister {
	return func(state manifest.ManifestGamesWriterState) error {
		var output []byte
		buf := new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "  ")
		_ = enc.Encode(state)
		output = buf.Bytes()
		err := ioutil.WriteFile(file, output, 0644)
		return err
	}		
}

func PersistMetadataProgress(file string) metadata.MetadataWriterStatePersister {
	return func(state metadata.MetadataGamesWriterState) error {
		var output []byte
		buf := new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "  ")
		_ = enc.Encode(state)
		output = buf.Bytes()
		err := ioutil.WriteFile(file, output, 0644)
		return err
	}		
}

func CleanupFile(file string) error {
	err := os.Remove(file)
	if err != nil && (!os.IsNotExist(err)){
		return err
	}

	return nil
}

//https://github.com/spf13/cobra/issues/216#issuecomment-703846787
func callPersistentPreRun(cmd *cobra.Command, args []string) {
	parent := cmd.Parent()
	if parent != nil {
		if parent.PersistentPreRun != nil {
			parent.PersistentPreRun(parent, args)
		}
	}
}

func getLanguageMap() map[string]string {
	langMap := make(map[string]string)
	langMap["english"] = "en"
	langMap["french"] = "fr"
	langMap["dutch"] = "nl"
	langMap["spanish"] = "es"
	langMap["portuguese_brazilian"] = "pt-br"
	langMap["russian"] = "ru"
	langMap["korean"] = "ko"
	langMap["chinese_simplified"] = "zh"
	langMap["japanese"] = "jp"
	langMap["polish"] = "pl"
	langMap["italian"] = "it"
	langMap["german"] = "de"
	langMap["czech"] = "cs"
	langMap["hungarian"] = "hu"
	langMap["portuguese"] = "pt"
	langMap["danish"] = "da"
	langMap["finnish"] = "fi"
	langMap["swedish"] = "sv"
	langMap["turkish"] = "tr"
	langMap["arabic"] = "ar"
	langMap["romanian"] = "ro"
	return langMap
}
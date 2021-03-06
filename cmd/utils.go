package cmd

import (
	"bytes"
	"encoding/json"
    "fmt"
    "gogcli/manifest"
    "gogcli/storage"
    "io/ioutil"
    "os"
    
	"github.com/spf13/cobra"
)

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

func getStorage(path string, storageType string, debugMode bool, loggerTag string) (storage.Storage, storage.Downloader) {
    if storageType != "fs" && storageType != "s3" {
        msg := fmt.Sprintf("Source storage type %s is invalid", storageType)
        fmt.Println(msg)
        os.Exit(1)
    }
    
    if storageType == "fs" {
        gameStorage := storage.GetFileSystem(path, debugMode, loggerTag)
        downloader := storage.FileSystemDownloader{gameStorage}
        return gameStorage, downloader
    } else {
        gameStorage, err := storage.GetS3StoreFromConfigFile(path, debugMode, loggerTag)
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

func processSerializableOutput(serializable interface{}, errs []error, terminal bool, file string) {
    hasErr := false
    var buf bytes.Buffer
    var output []byte
    var e Errors

    if len(errs) > 0 {
        for _, err := range errs {
            e.Errors = append(e.Errors, err.Error())
        }
        output, _ = json.Marshal(e)
        hasErr = true
    } else {
        output, _ = json.Marshal(serializable)
    }

    json.Indent(&buf, output, "", "  ")
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

//https://github.com/spf13/cobra/issues/216#issuecomment-703846787
func callPersistentPreRun(cmd *cobra.Command, args []string) { 
	parent := cmd.Parent()
	if parent != nil { 
        if parent.PersistentPreRun != nil { 
            parent.PersistentPreRun(parent, args) 
        } 
    } 
} 
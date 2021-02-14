package cmd

import (
    "fmt"
    "gogcli/storage"
    "os"
    
	"github.com/spf13/cobra"
)

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

//https://github.com/spf13/cobra/issues/216#issuecomment-703846787
func callPersistentPreRun(cmd *cobra.Command, args []string) { 
	parent := cmd.Parent()
	if parent != nil { 
        if parent.PersistentPreRun != nil { 
            parent.PersistentPreRun(parent, args) 
        } 
    } 
} 
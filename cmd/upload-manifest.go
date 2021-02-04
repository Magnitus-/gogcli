package cmd

import (
	"errors"
	"fmt"
	"gogcli/manifest"
	"gogcli/storage"
	"io"
	"os"
)

type ActionResult {
	gameId int
	fileKind string
	action manifest.FileAction
	fileName string
	fileSize int
	fileChecksum string
	err error
}

type GetDownloadHandle func(int, manifest.FileAction) (io.ReadCloser, int, string, error)

func addFileAction(
	gameId int, 
	fileKind string, 
	action manifest.FileAction, 
	result chan ActionResult, 
	actionErr chan error,
	s storage.Storage,
	fn GetDownloadHandle
) {
	result := ActionResult{
		gameId: gameId,
		fileKind: fileKind,
		action: action,
		fileName: action.Name,
	}
	handle, fSize, fname, err := GetDownloadHandle(gameId, action)
	if err != nil {
		result.err = err
		ActionResult <- result
		actionErr <- err
	} else {
		result.fileSize = fSize
		fChecksum, uploadErr := s.UploadFile(handle, gameId, fileKind, action.Name) ([]byte, error)
		if err != nil {
			result.err = uploadErr
			actionErr <- uploadErr
		} else {
			result.fileChecksum = fChecksum
			actionErr <- nil
		}
	}
}

func launchActions(a *manifest.GameActions, s storage.Storage, concurrency int, fn GetDownloadHandle, result chan ActionResult, actionErrsChan chan []error) {
	errs := make([]error)
	var actionErr chan error
	jobsRunning  := 0
	gameIds := make([]int, len(*a))
	
	idx := 0
	for id, _ := range (*a) {
		gameIds[idx] = id
		idx++
	}
	
	for true {
		if len(gameIds) > 0 {
			if (*a)[gameIds[0]].Action == "add" && len(errs) == 0 {
				err := s.AddGame(gameIds[0])
				if err != nil {
					errs := append(errs, err)
				}
				(*a)[gameIds[0]].Action = "added"
			}

			if (!(*a)[gameIds[0]].HasFileActions()) && len(errs) == 0 {
				if (*a)[gameIds[0]].Action == "remove" {
					err := s.AddRemove(gameIds[0])
					if err != nil {
						errs := append(errs, err)
					}
					(*a)[gameIds[0]].Action = "removed"
				}

				id := gameIds[0]
				if len(gameIds) > 1 {
					gameIds = gameIds[1:]
				}
				delete((*a), id)
			} else if len(errs) == 0 {
				fileAction, fileKind, _ := (*a)[gameIds[0]].ExtractFileAction()
				if fileAction.Action == "add" {
					go addFileAction(
						gameIds[0],
						fileKind,
						fileAction,
						result,
						actionErr,
						s,
						fn
					)
					concurrency--
					jobsRunning++
				} else if fileAction.Action == "remove" {
					err := s.RemoveFile(gameIds[0], fileKind, fileAction.Name)
					if err != nil {
						errs := append(errs, err)
					}
				} else {
					msg := fmt.Sprintf("launchActions(...) -> %s is not a valid type for file actions", fileAction.Action)
					errs := append(errs, errors.New(msg))
				}
			}
		}		

		allDone = (len(gameIds) == 0 && jobsRunning == 0) || len(errs) > 0
		waitOnLingeringJobs = (len(gameIds) == 0 || len(errs) > 0) && jobsRunning > 0
		if allDone {
			break;
		}
		else if concurrency <= 0 || waitOnLingeringJobs {
			err := <- actionErr
			if err != nil {
				errs = append(errs, err)
			}
			jobsRunning--
		}
	}

	actionErrsChan <- errs
}

func keepManifestUpdated(m *manifest.Manifest, s storage.Storage, result chan ActionResult, manifestErrsChan chan []error) {
	//TODO
}

func processGameActions(m *manifest.Manifest, a *manifest.GameActions, s storage.Storage, concurrency int, fn GetDownloadHandle) {
	actionErrsChan chan []error
	manifestErrsChan chan []error
	result chan ActionResult
	go launchActions(a.DeepCopy(), s, concurrency, fn, result, actionErrsChan)
	go keepManifestUpdated(m, s, result, manifestErrsChan)
	//TODO: Some sync to figure out between action and manifest processing if either returns an error
	actionErrs := <- actionErrsChan
	manifestErrs := <- manifestErrsChan
}

func uploadManifest(m *manifest.Manifest, s storage.Storage, concurrency int, fn GetDownloadHandle) {
	var hasActions bool
	var actions *manifest.GameActions
	var err error

	hasActions, err = s.HasActions()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if hasActions {
		fmt.Println("An unfinished manifest apply is already in progress. Aborting.")
		os.Exit(1)	
	}

	actions, err = storage.PlanManifest(m, s)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = s.StoreActions(actions)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)	
	}

	err = s.StoreManifest(m)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)	
	}

	processGameActions(m, actions, s, concurrency, fn)
}
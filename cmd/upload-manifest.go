package cmd

import (
	"errors"
	"fmt"
	"gogcli/manifest"
	"gogcli/storage"
	"io"
	"os"
)

type ActionResult struct {
	gameId int
	fileKind string
	action manifest.FileAction
	fileName string
	fileSize int
	fileChecksum string
	err error
	end bool
}

type GetDownloadHandle func(int, manifest.FileAction) (io.ReadCloser, int, string, error)

func addFileAction(
	gameId int, 
	fileKind string, 
	action manifest.FileAction, 
	result chan ActionResult, 
	actionErr chan error,
	s storage.Storage,
	fn GetDownloadHandle,
) {
	r := ActionResult{
		gameId: gameId,
		fileKind: fileKind,
		action: action,
		fileName: action.Name,
		end: false,
	}
	handle, fSize, _, err := fn(gameId, action)
	if err != nil {
		r.err = err
		actionErr <- err
	} else {
		r.fileSize = fSize
		fChecksum, uploadErr := s.UploadFile(handle, gameId, fileKind, action.Name)
		if err != nil {
			r.err = uploadErr
			actionErr <- uploadErr
		} else {
			r.fileChecksum = string(fChecksum)
			result <- r
			actionErr <- nil
		}
	}
}

func launchActions(a *manifest.GameActions, s storage.Storage, concurrency int, fn GetDownloadHandle, result chan ActionResult, actionErrsChan chan []error) {
	errs := make([]error, 0)
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
			ga := (*a)[gameIds[0]]
			if (*a)[gameIds[0]].Action == "add" && len(errs) == 0 {
				err := s.AddGame(gameIds[0])
				if err != nil {
					errs = append(errs, err)
				}
				ga.Action = "added"
				(*a)[gameIds[0]] = ga
			}

			if (!ga.HasFileActions()) && len(errs) == 0 {
				if (*a)[gameIds[0]].Action == "remove" {
					err := s.RemoveGame(gameIds[0])
					if err != nil {
						errs = append(errs, err)
					}
					ga.Action = "removed"
					(*a)[gameIds[0]] = ga
				}

				id := gameIds[0]
				if len(gameIds) > 1 {
					gameIds = gameIds[1:]
				}
				delete((*a), id)
			} else if len(errs) == 0 {
				fileAction, fileKind, _ := ga.ExtractFileAction()
				(*a)[gameIds[0]] = ga
				if fileAction.Action == "add" {
					go addFileAction(
						gameIds[0],
						fileKind,
						fileAction,
						result,
						actionErr,
						s,
						fn,
					)
					concurrency--
					jobsRunning++
				} else if fileAction.Action == "remove" {
					err := s.RemoveFile(gameIds[0], fileKind, fileAction.Name)
					if err != nil {
						errs = append(errs, err)
					}
				} else {
					msg := fmt.Sprintf("launchActions(...) -> %s is not a valid type for file actions", fileAction.Action)
					errs = append(errs, errors.New(msg))
				}
			}
		}		

		allDone := (len(gameIds) == 0 && jobsRunning == 0) || len(errs) > 0
		waitOnLingeringJobs := (len(gameIds) == 0 || len(errs) > 0) && jobsRunning > 0
		if allDone {
			break
		} else if concurrency <= 0 || waitOnLingeringJobs {
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
	errs := make([]error, 0)
	for true {
		r := <- result
		if r.end {
			break
		}
		err := m.FillMissingFileInfo(r.gameId, r.fileKind, r.fileName, r.fileSize, r.fileChecksum)
		if err != nil {
			errs = append(errs, err)
		} else {
			err = s.StoreManifest(m)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	manifestErrsChan <- errs
}

func processGameActions(m *manifest.Manifest, a *manifest.GameActions, s storage.Storage, concurrency int, fn GetDownloadHandle) []error {
	//Missing: keep game actions updated in storage
	actionErrsChan := make(chan []error)
	manifestErrsChan := make(chan []error)
	result := make(chan ActionResult)
	go launchActions(a.DeepCopy(), s, concurrency, fn, result, actionErrsChan)
	go keepManifestUpdated(m, s, result, manifestErrsChan)
	
	actionErrs := <- actionErrsChan
	result <- ActionResult{end: true}
	manifestErrs := <- manifestErrsChan

	errs := make([]error, len(actionErrs) + len(manifestErrs))
	for idx, err := range actionErrs {
		errs[idx] = err
	}
	for idx, err := range manifestErrs {
		errs[idx + len(actionErrs)] = err
	}
	return errs
}

func uploadManifest(m *manifest.Manifest, s storage.Storage, concurrency int, fn GetDownloadHandle) []error {
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

	return processGameActions(m, actions, s, concurrency, fn)
	//TODO: Clear manifest
}
package storage

import (
	"errors"
	"gogcli/manifest"
)

type ActionResult struct {
	gameId int
	fileKind string
	action manifest.FileAction
	fileName string
	fileSize int64
	fileChecksum string
	err error
	end bool
}

func addFileAction(
	gameId int, 
	fileKind string, 
	action manifest.FileAction, 
	result chan ActionResult, 
	actionErr chan error,
	s Storage,
	d Downloader,
) {
	r := ActionResult{
		gameId: gameId,
		fileKind: fileKind,
		action: action,
		fileName: action.Name,
		end: false,
	}
	handle, fSize, _, err := d.Download(gameId, action)
	if err != nil {
		r.err = err
		actionErr <- err
	} else {
		r.fileSize = fSize
		fChecksum, uploadErr := s.UploadFile(handle, gameId, fileKind, action.Name, fSize)
		if err != nil {
			r.err = uploadErr
			actionErr <- uploadErr
		} else {
			r.fileChecksum = fChecksum
			result <- r
			actionErr <- nil
		}
	}
}

func launchActions(a *manifest.GameActions, s Storage, concurrency int, d Downloader, result chan ActionResult, actionErrsChan chan []error) {
	errs := make([]error, 0)
	actionErr := make(chan error)
	jobsRunning  := 0

	iterator := manifest.NewActionsInterator(*a)
	
	for true {
		if iterator.HasMore() && len(errs) == 0 {
			action := iterator.Next()
			if (!action.IsFileAction) {
				if action.GameAction == "add" {
					err := s.AddGame(action.GameId)
					if err != nil {
						errs = append(errs, err)
					}
				} else if action.GameAction == "remove" {
					err := s.RemoveGame(action.GameId)
					if err != nil {
						errs = append(errs, err)
					}
				}
			} else {
				fileActionPtr := action.FileActionPtr
				if (*fileActionPtr).Action == "add" {
					concurrency--
					jobsRunning++
					go addFileAction(
						action.GameId,
						(*fileActionPtr).Kind,
						(*fileActionPtr),
						result,
						actionErr,
						s,
						d,
					)
				} else if (*fileActionPtr).Action == "remove" {
					err := s.RemoveFile(action.GameId, (*fileActionPtr).Kind, (*fileActionPtr).Name)
					if err != nil {
						errs = append(errs, err)
					}
				}
			}
		}		

		allDone := ((!iterator.HasMore()) && jobsRunning == 0) || len(errs) > 0
		waitOnLingeringJobs := ((!iterator.HasMore()) || len(errs) > 0) && jobsRunning > 0
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

func keepManifestUpdated(m *manifest.Manifest, s Storage, result chan ActionResult, manifestErrsChan chan []error) {
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

func KeepActionsUpdated(g *manifest.GameActions, s Storage, action chan manifest.Action, actionErrsChan chan []error) {

}

func processGameActions(m *manifest.Manifest, a *manifest.GameActions, s Storage, concurrency int, d Downloader) []error {
	//Missing: keep game actions updated in storage
	actionErrsChan := make(chan []error)
	manifestErrsChan := make(chan []error)
	result := make(chan ActionResult)
	go launchActions(a.DeepCopy(), s, concurrency, d, result, actionErrsChan)
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

func UploadManifest(m *manifest.Manifest, s Storage, concurrency int, d Downloader) []error {
	var hasActions bool
	var actions *manifest.GameActions
	var err error

	hasActions, err = s.HasActions()
	if err != nil {
		return []error{err}
	}
	if hasActions {
		return []error{errors.New("An unfinished manifest apply is already in progress. Aborting.")}
	}

	actions, err = PlanManifest(m, s)
	if err != nil {
		return []error{err}
	}

	err = s.StoreActions(actions)
	if err != nil {
		return []error{err}
	}

	err = s.StoreManifest(m)
	if err != nil {
		return []error{err}	
	}

	errs := processGameActions(m, actions, s, concurrency, d)
	if len(errs) == 0 {
		err := s.RemoveActions()
		if err != nil {
			return []error{err}
		}
	}
	return errs
}
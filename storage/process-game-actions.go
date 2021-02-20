package storage

import "gogcli/manifest"

type DoneAction struct {
	action manifest.Action
	end bool
}

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
	actionResult chan ActionResult,
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
			actionResult <- r
			actionErr <- nil
		}
	}
}

func launchActions(iterator *manifest.ActionsIterator, s Storage, concurrency int, d Downloader, result chan ActionResult, doneAction chan DoneAction, actionErrsChan chan []error) {
	errs := make([]error, 0)
	actionErr := make(chan error)
	jobsRunning  := 0
	
	for true {
		if iterator.ShouldContinue() && len(errs) == 0 {
			action := iterator.Next()
			if (!action.IsFileAction) {
				if action.GameAction == "add" {
					err := s.AddGame(action.GameId)
					if err != nil {
						errs = append(errs, err)
					} else {
						doneAction <- DoneAction{action: action, end: false}
					}
				} else if action.GameAction == "remove" {
					err := s.RemoveGame(action.GameId)
					if err != nil {
						errs = append(errs, err)
					} else {
						doneAction <- DoneAction{action: action, end: false}
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
					} else {
						doneAction <- DoneAction{action: action, end: false}
					}
				}
			}
		}		

		allDone := ((!iterator.ShouldContinue()) && jobsRunning == 0) || len(errs) > 0
		waitOnLingeringJobs := ((!iterator.ShouldContinue()) || len(errs) > 0) && jobsRunning > 0
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

func keepManifestUpdated(m *manifest.Manifest, s Storage, result chan ActionResult, doneAction chan DoneAction, errsChan chan []error) {
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
			} else {
				action := manifest.Action{
					GameId: r.gameId,
					IsFileAction: true,
					FileActionPtr: &r.action,
					GameAction: "",
				}
				doneAction <- DoneAction{action: action, end: false}
			}
		}
	}

	errsChan <- errs
}

func KeepActionsUpdated(g *manifest.GameActions, s Storage, doneAction chan DoneAction, errsChan chan []error) {
	errs := make([]error, 0)
	for true {
		d := <- doneAction
		if d.end {
			break
		}
		g.ApplyAction(d.action)
		err := s.StoreActions(g)
		if err != nil {
			errs = append(errs, err)
		}
	}
	errsChan <- errs
}

func processGameActions(m *manifest.Manifest, a *manifest.GameActions, s Storage, concurrency int, d Downloader, gamesMax int) []error {
	actionErrsChan := make(chan []error)
	actionResult := make(chan ActionResult)
	manifestUpdateErrsChan := make(chan []error)
	actionsUpdateErrsChan := make(chan []error)
	doneAction := make(chan DoneAction)
	
	iterator := manifest.NewActionsInterator(*a, gamesMax)
	go launchActions(iterator, s, concurrency, d, actionResult, doneAction, actionErrsChan)
	go keepManifestUpdated(m, s, actionResult, doneAction, manifestUpdateErrsChan)
	go KeepActionsUpdated(a.DeepCopy(), s, doneAction, actionsUpdateErrsChan)
	actionErrs := <- actionErrsChan
	actionResult <- ActionResult{end: true}
	manifestUpdateErrs := <- manifestUpdateErrsChan
	doneAction <- DoneAction{end: true}
	actionsUpdateErrs := <- actionsUpdateErrsChan
	errs := make([]error, len(actionErrs) + len(manifestUpdateErrs) + len(actionsUpdateErrs))
	for idx, err := range actionErrs {
		errs[idx] = err
	}
	for idx, err := range manifestUpdateErrs {
		errs[idx + len(actionErrs)] = err
	}
	for idx, err := range actionsUpdateErrs {
		errs[idx + len(actionErrs) + len(manifestUpdateErrs)] = err	
	}

	if len(errs) == 0 && (!iterator.HasMore()) {
		err := s.RemoveActions()
		if err != nil {
			return []error{err}
		}
	}

	return errs
}
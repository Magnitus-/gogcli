package storage

import (
	"errors"
	"fmt"
	"log"
	"gogcli/logging"
	"gogcli/manifest"
	"os"
	"strings"
)

type DoneAction struct {
	action manifest.Action
	end bool
}

type ActionResult struct {
	gameId int64
	fileKind string
	action manifest.FileAction
	fileName string
	fileSize int64
	fileChecksum string
	err error
	end bool
}

type ActionsProcessor struct {
	concurrency            int
	retries                int
	gamesMax               int
	gamesSort              manifest.ActionsIteratorSort
    logger                 *logging.Logger
	actionErrChan          chan error
	actionsErrsChan        chan []error
	actionResultChan       chan ActionResult
	manifestUpdateErrsChan chan []error
	actionsUpdateErrsChan  chan []error
	doneActionChan         chan DoneAction
}

func GetActionsProcessor(
	concurrency int, 
	retries int, 
	gamesMax int, 
	gamesSort manifest.ActionsIteratorSort,
	logSource *logging.Source,
) (ActionsProcessor) {
	return ActionsProcessor{
		concurrency: concurrency,
		retries: retries,
		gamesMax: gamesMax,
		gamesSort: gamesSort,
		logger: logSource.CreateLogger(os.Stdout, "DOWNLOADS PROCESSING: ", log.Lshortfile),
		actionErrChan:  make(chan error),
		actionsErrsChan: make(chan []error),
		actionResultChan: make(chan ActionResult),
		manifestUpdateErrsChan: make(chan []error),
		actionsUpdateErrsChan: make(chan []error),
		doneActionChan: make(chan DoneAction),
	}
}


func (p ActionsProcessor) addFileAction(
	gameId int64,
	fileKind string,
	fileInfo manifest.FileInfo,
	action manifest.FileAction, 
	s Storage,
	d Downloader,
) {
	fn := fmt.Sprintf("addFileAction(gameId=%d, fileInfo={Kind=%s, Name=%s, ...}, ...)", gameId, fileInfo.Kind, fileInfo.Name)
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
		p.actionErrChan <- err
		return 
	}
	defer handle.Close()

	if fileInfo.Size > 0 && fileInfo.Size != fSize {
		msg := fmt.Sprintf("%s -> Download file size of %d does not match expected file size of %d", fn, fSize, fileInfo.Size)  
		r.err = errors.New(msg)
		p.actionErrChan <- errors.New(msg)
		return
	}
	
	r.fileSize = fSize
	fChecksum, uploadErr := s.UploadFile(handle, gameId, fileKind, action.Name, fSize)
	if err != nil {
		r.err = uploadErr
		p.actionErrChan <- uploadErr
		return
	}

	if fileInfo.Checksum != "" && fileInfo.Checksum != fChecksum {
		msg := fmt.Sprintf("%s -> Download file checksum of %s does not match expected file checksum of %s", fn, fChecksum, fileInfo.Checksum)  
		r.err = errors.New(msg)
		p.actionErrChan <- errors.New(msg)
		return
	}

	if fileInfo.Checksum == "" && strings.HasSuffix(fileInfo.Name, ".zip") {
		err = ValidateZipArchive(s, gameId, fileKind, fileInfo.Name)
		if err != nil {
			msg := fmt.Sprintf("%s -> Error occured while validating Zip archive %s: %s", fn, fileInfo.Name, err.Error())  
			r.err = errors.New(msg)
			p.actionErrChan <- errors.New(msg)
			return
		}
	}

	r.fileChecksum = fChecksum
	p.actionResultChan <- r
	p.actionErrChan <- nil
}

func (p ActionsProcessor) launchActions( m *manifest.Manifest, iterator *manifest.ActionsIterator, s Storage, d Downloader) {
	errs := make([]error, 0)
	jobsRunning  := 0
	concurrency := p.concurrency
	
	for true {
		if iterator.ShouldContinue() && len(errs) == 0 {
			action, nextErr := iterator.Next()
			if nextErr != nil {
				errs = append(errs, nextErr)
			} else if (!action.IsFileAction) {
				if action.GameAction == "add" {
					err := s.AddGame(action.GameId)
					if err != nil {
						errs = append(errs, err)
					} else {
						p.doneActionChan <- DoneAction{action: action, end: false}
					}
				} else if action.GameAction == "remove" {
					err := s.RemoveGame(action.GameId)
					if err != nil {
						errs = append(errs, err)
					} else {
						p.doneActionChan <- DoneAction{action: action, end: false}
					}
				}
			} else {
				fileActionPtr := action.FileActionPtr
				if (*fileActionPtr).Action == "add" {
					fileInfo, err := (*m).GetFileActionFileInfo(action.GameId, (*fileActionPtr))
					if err != nil {
						errs = append(errs, err)
					} else {
						concurrency--
						jobsRunning++
						go p.addFileAction(
							action.GameId,
							(*fileActionPtr).Kind,
							fileInfo,
							(*fileActionPtr),
							s,
							d,
						)
					}
				} else if (*fileActionPtr).Action == "remove" {
					err := s.RemoveFile(action.GameId, (*fileActionPtr).Kind, (*fileActionPtr).Name)
					if err != nil {
						errs = append(errs, err)
					} else {
						p.doneActionChan <- DoneAction{action: action, end: false}
					}
				}
			}
		}		
		endWhenPossible := (!iterator.ShouldContinue()) || (len(errs) > 0)
		allDone := endWhenPossible && jobsRunning == 0
		if allDone {
			break
		} else if (concurrency <= 0 && jobsRunning > 0) || (endWhenPossible && jobsRunning > 0) {
			err := <- p.actionErrChan
			if err != nil {
				errs = append(errs, err)
			}
			jobsRunning--
		}
	}

	p.actionsErrsChan <- errs
}

func (p ActionsProcessor) keepManifestUpdated(m *manifest.Manifest, s Storage) {
	errs := make([]error, 0)
	for true {
		r := <- p.actionResultChan
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
				p.doneActionChan <- DoneAction{action: action, end: false}
			}
		}
	}

	p.manifestUpdateErrsChan <- errs
}

func (p ActionsProcessor) keepActionsUpdated(g *manifest.GameActions, s Storage) {
	errs := make([]error, 0)
	for true {
		d := <- p.doneActionChan
		if d.end {
			break
		}
		g.ApplyAction(d.action)
		err := s.StoreActions(g)
		if err != nil {
			errs = append(errs, err)
		}
	}
	p.actionsUpdateErrsChan <- errs
}

func (p ActionsProcessor) ProcessGameActions(m *manifest.Manifest, a *manifest.GameActions, s Storage, d Downloader) []error {
	iterator := manifest.NewActionsIterator(*a, p.gamesMax)
	iterator.Sort(p.gamesSort, m)
	go p.launchActions(m, iterator, s, d)
	go p.keepManifestUpdated(m, s)
	go p.keepActionsUpdated(a.DeepCopy(), s)
	actionErrs := <- p.actionsErrsChan
	p.actionResultChan <- ActionResult{end: true}
	manifestUpdateErrs := <- p.manifestUpdateErrsChan
	p.doneActionChan <- DoneAction{end: true}
	actionsUpdateErrs := <- p.actionsUpdateErrsChan
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
		err = s.RemoveSource()
		if err != nil {
			return []error{err}
		}
	}

	return errs
}
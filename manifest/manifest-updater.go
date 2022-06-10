package manifest

import (
	"errors"
	"fmt"
	"log"
	"os"
	"gogcli/logging"
)

type ManifestGamesWriterState struct {
	Manifest Manifest
	GameIds  []int64
}

type ManifestGamesWriter struct {
	State  ManifestGamesWriterState
	logger *logging.Logger
}

type ManifestGamesWriterResult struct {
	Warnings []error
	Errors   []error
}

type ManifestGameGetterResult struct {
	Game ManifestGame
	Warnings []error
	Errors   []error
}

func NewManifestGamesWriter(state ManifestGamesWriterState, logSource *logging.Source) *ManifestGamesWriter {
	return &ManifestGamesWriter{
		State: state,
		logger: logSource.CreateLogger(os.Stdout, "[manifest writer] ", log.Lmsgprefix),
	}
}

type ManifestGameGetter func(chan struct{}, []int64) chan ManifestGameGetterResult
type ManifestWriterStatePersister func(state ManifestGamesWriterState) error

func (w *ManifestGamesWriter) Update(getter ManifestGameGetter, persister ManifestWriterStatePersister) ManifestGamesWriterResult {
	result := ManifestGamesWriterResult{
		Warnings: []error{},
		Errors: []error{},
	}

	done := make(chan struct{})
	defer close(done)

	gameResults := getter(done, (*w).State.GameIds)

	for gameResult := range gameResults {
		if len(gameResult.Errors) > 0 {
			result.Errors = append(result.Errors, gameResult.Errors...)
			return result
		}

		if len(gameResult.Warnings) > 0 {
			result.Warnings = append(result.Warnings, gameResult.Warnings...)
		}

		gameIds := RemoveIdFromList((*w).State.GameIds, gameResult.Game.Id)
		(*w).State.GameIds = gameIds
		(*w).State.Manifest.Games = append((*w).State.Manifest.Games, gameResult.Game)
		persistErr := persister((*w).State)
		if persistErr != nil {
			result.Errors = append(result.Errors, errors.New(fmt.Sprintf("Failed to persist state due to error: %s", persistErr.Error())))
			return result
		}
		(*w).logger.Info(fmt.Sprintf("Got all info on game with id %d. %d games left to process", gameResult.Game.Id, len((*w).State.GameIds)))	
	}

	return result
}

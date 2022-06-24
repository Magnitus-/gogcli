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

type ManifestGameGetterGame struct {
	Game ManifestGame
	Warnings []error
	Errors   []error
}

type ManifestGameGetterGameIds struct {
	Ids   []int64
	Error error
}

func NewManifestGamesWriterState(filter ManifestFilter, GameIds []int64) ManifestGamesWriterState {
	m := NewEmptyManifest(filter)
	return ManifestGamesWriterState{
		Manifest: *m,
		GameIds: GameIds,
	}
}

func NewManifestGamesWriter(state ManifestGamesWriterState, logSource *logging.Source) *ManifestGamesWriter {
	return &ManifestGamesWriter{
		State: state,
		logger: logSource.CreateLogger(os.Stdout, "[manifest writer] ", log.Lmsgprefix),
	}
}

type ManifestGameGetter func(<-chan struct{}, []int64, ManifestFilter) (<-chan ManifestGameGetterGame, <-chan ManifestGameGetterGameIds)
type ManifestWriterStatePersister func(state ManifestGamesWriterState) error

func (w *ManifestGamesWriter) Write(getter ManifestGameGetter, persister ManifestWriterStatePersister) ManifestGamesWriterResult {
	result := ManifestGamesWriterResult{
		Warnings: []error{},
		Errors: []error{},
	}

	done := make(chan struct{})
	defer close(done)

	gameCh, gameIdsCh := getter(done, (*w).State.GameIds, (*w).State.Manifest.Filter)

	IdsResult := <- gameIdsCh
	if IdsResult.Error != nil {
		result.Errors = append(result.Errors, IdsResult.Error)
		return result
	}

	(*w).State.GameIds = IdsResult.Ids
	(*w).logger.Info(fmt.Sprintf("Generating/Updating manifest for %d games", len((*w).State.GameIds)))	

	for gameResult := range gameCh {
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
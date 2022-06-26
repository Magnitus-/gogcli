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
	Warnings []string
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
		Warnings: []string{},
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

func (w *ManifestGamesWriter) Write(getter ManifestGameGetter, persister ManifestWriterStatePersister) []error {
	errs := []error{}

	done := make(chan struct{})
	defer close(done)

	gameCh, gameIdsCh := getter(done, (*w).State.GameIds, (*w).State.Manifest.Filter)

	IdsResult := <- gameIdsCh
	if IdsResult.Error != nil {
		errs = append(errs, IdsResult.Error)
		return errs
	}

	(*w).State.GameIds = IdsResult.Ids
	(*w).logger.Info(fmt.Sprintf("Generating/Updating manifest for %d games", len((*w).State.GameIds)))	

	for gameResult := range gameCh {
		if len(gameResult.Errors) > 0 {
			errs = append(errs, gameResult.Errors...)
			return errs
		}

		if len(gameResult.Warnings) > 0 {
			for _, warning := range gameResult.Warnings {
				(*w).State.Warnings = append((*w).State.Warnings, warning.Error())
			}
		}

		gameIds := RemoveIdFromList((*w).State.GameIds, gameResult.Game.Id)
		(*w).State.GameIds = gameIds
		(*w).State.Manifest.Games = append((*w).State.Manifest.Games, gameResult.Game)
		persistErr := persister((*w).State)
		if persistErr != nil {
			errs = append(errs, errors.New(fmt.Sprintf("Failed to persist state due to error: %s", persistErr.Error())))
			return errs
		}
		(*w).logger.Info(fmt.Sprintf("Got all info on game with id %d. %d games left to process", gameResult.Game.Id, len((*w).State.GameIds)))	
	}

	return errs
}

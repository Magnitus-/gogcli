package metadata

import (
	"errors"
	"fmt"
	"log"
	"os"
	"gogcli/logging"
)

type MetadataGamesWriterState struct {
	Metadata Metadata
	GameIds  []int64
	Warnings []string
}

type MetadataGamesWriter struct {
	State  MetadataGamesWriterState
	logger *logging.Logger
}

type MetadataGamesWriterResult struct {
	Warnings []error
	Errors   []error
}

type MetadataGameGetterGame struct {
	Game MetadataGame
	Warnings []error
	Errors   []error
}

type MetadataGameGetterGameIds struct {
	Ids   []int64
	Error error
}

func NewMetadataGamesWriterState(GameIds []int64) MetadataGamesWriterState {
	m := NewEmptyMetadata()
	return MetadataGamesWriterState{
		Metadata: *m,
		GameIds: GameIds,
		Warnings: []string{},
	}
}

func NewMetadataGamesWriter(state MetadataGamesWriterState, logSource *logging.Source) *MetadataGamesWriter {
	return &MetadataGamesWriter{
		State: state,
		logger: logSource.CreateLogger(os.Stdout, "[metadata writer] ", log.Lmsgprefix),
	}
}

type MetadataGameGetter func(<-chan struct{}, []int64) (<-chan MetadataGameGetterGame, <-chan MetadataGameGetterGameIds)
type MetadataWriterStatePersister func(state MetadataGamesWriterState) error

func (w *MetadataGamesWriter) Write(getter MetadataGameGetter, persister MetadataWriterStatePersister) []error {
	errs := []error{}

	done := make(chan struct{})
	defer close(done)

	gameCh, gameIdsCh := getter(done, (*w).State.GameIds)

	IdsResult := <- gameIdsCh
	if IdsResult.Error != nil {
		errs = append(errs, IdsResult.Error)
		return errs
	}

	(*w).State.GameIds = IdsResult.Ids
	(*w).logger.Info(fmt.Sprintf("Generating metadata for %d games", len((*w).State.GameIds)))	

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
		(*w).State.Metadata.Games = append((*w).State.Metadata.Games, gameResult.Game)
		persistErr := persister((*w).State)
		if persistErr != nil {
			errs = append(errs, errors.New(fmt.Sprintf("Failed to persist state due to error: %s", persistErr.Error())))
			return errs
		}
		(*w).logger.Info(fmt.Sprintf("Got all metadata info on game with id %d. %d games left to process", gameResult.Game.Id, len((*w).State.GameIds)))	
	}

	return errs
}


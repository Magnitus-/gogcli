package manifest

import (
	"errors"
	"fmt"
)

type Action struct {
	GameId int64
	IsFileAction bool
	FileActionPtr *FileAction
	GameAction string
}

type ActionsIterator struct {
	gameActions           GameActions
	gameIds               []int64
	currentGameActionDone bool
	installerNames        []string
	extraNames            []string
	maxGames              int
	processedGames        int
}

func NewActionsInterator(a GameActions, maxGames int) *ActionsIterator {
	gameIds := a.GetGameIds()
	currentGameAction := a[gameIds[0]]
	new := &ActionsIterator{
		gameActions: a,
		gameIds: gameIds,
		currentGameActionDone: currentGameAction.Action == "update",
		installerNames: currentGameAction.GetInstallerNames(),
		extraNames: currentGameAction.GetExtraNames(),
		maxGames: maxGames,
		processedGames: 0,
	}

	return new
}

func (i *ActionsIterator) Stringify() string {
	return fmt.Sprintf(
		"{'gameId': %d, 'gameActionDone': %t, 'installersLeft': %d, 'extrasLeft': %d, 'processedGames': %d, 'maxGames': %d}",
		(*i).gameIds[0],
		(*i).currentGameActionDone,
		len((*i).installerNames),
		len((*i).extraNames),
		(*i).processedGames,
		(*i).maxGames,
	)
}

func (i *ActionsIterator) ShouldContinue() bool {
	return i.HasMore() && (i.maxGames == -1 || i.processedGames < i.maxGames)
}

func (i *ActionsIterator) HasMore() bool {
	if len(i.gameIds) == 0 {
		return false
	}

	moreFileActions := len(i.gameIds) > 1 || len(i.installerNames) > 0 || len(i.extraNames) > 0
	if moreFileActions {
		return true
	}

	return !i.currentGameActionDone
}

func (i *ActionsIterator) Next() (Action, error) {
	if !i.ShouldContinue() {
		return Action{
			GameId: -1,
			IsFileAction: false,
			FileActionPtr: nil,
			GameAction: "",
		}, errors.New("*ActionsIterator.Next() -> End of iterator, cannot fetch anymore")
	}

	currentGameId := i.gameIds[0]
	currentGame := i.gameActions[currentGameId]
	
	remainingFileActions := len(i.extraNames) + len(i.installerNames)
	onlyOneFileActionRemains := i.currentGameActionDone && remainingFileActions == 1
	onlyOneGameActionRemains := (!i.currentGameActionDone) && remainingFileActions == 0
	if onlyOneFileActionRemains || onlyOneGameActionRemains {
		i.processedGames++
	}

	if (!i.currentGameActionDone) && currentGame.Action == "add" {
		i.currentGameActionDone = true
		return Action{
			GameId: currentGameId,
			IsFileAction: false,
			FileActionPtr: nil,
			GameAction: "add",
		}, nil
	}

	if len(i.installerNames) > 0 {
		name := i.installerNames[0]
		i.installerNames = i.installerNames[1:]
		fileAction := currentGame.InstallerActions[name]
		return Action{
			GameId: currentGameId,
			IsFileAction: true,
			FileActionPtr: &fileAction,
			GameAction: "",
		}, nil
	}

	if len(i.extraNames) > 0 {
		name := i.extraNames[0]
		i.extraNames = i.extraNames[1:]
		fileAction := currentGame.ExtraActions[name]
		return Action{
			GameId: currentGameId,
			IsFileAction: true,
			FileActionPtr: &fileAction,
			GameAction: "",
		}, nil
	}

	if (!i.currentGameActionDone) && currentGame.Action == "remove" {
		i.currentGameActionDone = true
		return Action{
			GameId: currentGameId,
			IsFileAction: false,
			FileActionPtr: nil,
			GameAction: "remove",
		}, nil
	}

	i.gameIds = i.gameIds[1:]
	i.currentGameActionDone = i.gameActions[i.gameIds[0]].Action == "update"
	currentGameAction := i.gameActions[i.gameIds[0]]
	i.installerNames = currentGameAction.GetInstallerNames()
	i.extraNames = currentGameAction.GetExtraNames()
	return i.Next()
}
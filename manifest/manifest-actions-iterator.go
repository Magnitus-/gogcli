package manifest

import (
	"errors"
	"fmt"
	"sort"
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

type ActionsIteratorSort struct {
	gameIds []int64
	criteria string
	ascending bool
}

func NewActionIteratorSort(gameIds []int64, criteria string, ascending bool) ActionsIteratorSort {
	return ActionsIteratorSort{
		gameIds,
		criteria,
		ascending,
	}
}

func (i *ActionsIterator) Sort(gamesSort ActionsIteratorSort, m* Manifest) {
	manifestGames := make(map[int64]ManifestGame)
	if gamesSort.criteria != "none" && gamesSort.criteria != "id" {
		for _, game := range (*m).Games {
			manifestGames[game.Id] = game
		}
	}
	
	if gamesSort.criteria != "none" {
		sort.Slice((*i).gameIds, func(x, y int) bool {
			gameIdX := (*i).gameIds[x]
			gameIdY := (*i).gameIds[y]
			if gamesSort.criteria == "size" {
				if gamesSort.ascending {
					return manifestGames[gameIdX].VerifiedSize < manifestGames[gameIdY].VerifiedSize
				}

				return manifestGames[gameIdY].VerifiedSize < manifestGames[gameIdX].VerifiedSize
			} else if gamesSort.criteria == "title" {
				if gamesSort.ascending {
					return manifestGames[gameIdX].Title < manifestGames[gameIdY].Title
				}
				
				return manifestGames[gameIdY].Title < manifestGames[gameIdX].Title
			}

			//We assume the criteria is "id" otherwise
			if gamesSort.ascending {
				return gameIdX < gameIdY
			} else {
				return gameIdY < gameIdX
			}
		})
	}

	for _, gameId := range gamesSort.gameIds {
		for idx, id := range (*i).gameIds {
			if id == gameId {
				end := (*i).gameIds[idx+1:]
				beginning := append([]int64{(*i).gameIds[idx]}, (*i).gameIds[0:idx]...)
				(*i).gameIds = append(beginning, end...)
				break
			}
		}
	}

	currentGameAction := (*i).gameActions[(*i).gameIds[0]]
	installerNames := currentGameAction.GetInstallerNames()
	sort.Strings(installerNames)
	extraNames := currentGameAction.GetExtraNames()
	sort.Strings(extraNames)
	(*i).currentGameActionDone = (currentGameAction.Action == "update")
	(*i).installerNames = installerNames
	(*i).extraNames = extraNames
}

func NewActionsIterator(a GameActions, maxGames int) *ActionsIterator {
	gameIds := a.GetGameIds()
	currentGameAction := a[gameIds[0]]
	installerNames := currentGameAction.GetInstallerNames()
	sort.Strings(installerNames)
	extraNames := currentGameAction.GetExtraNames()
	sort.Strings(extraNames)
	new := &ActionsIterator{
		gameActions: a,
		gameIds: gameIds,
		currentGameActionDone: currentGameAction.Action == "update",
		installerNames: installerNames,
		extraNames: extraNames,
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

func (i *ActionsIterator) GetProgress() (int, int, int) {
	return i.processedGames, len(i.gameIds), i.maxGames
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
	installerNames := currentGameAction.GetInstallerNames()
	sort.Strings(installerNames)
	extraNames := currentGameAction.GetExtraNames()
	sort.Strings(extraNames)
	i.installerNames = installerNames
	i.extraNames = extraNames
	return i.Next()
}
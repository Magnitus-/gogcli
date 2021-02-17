package manifest

type Action struct {
	GameId int
	IsFileAction bool
	FileActionPtr *FileAction
	GameAction string
}

type ActionsIterator struct {
	gameActions           GameActions
	gameIds               []int
	currentGameActionDone bool
	installerNames        []string
	extraNames            []string
}

func NewActionsInterator(a GameActions) *ActionsIterator {
	gameIds := a.GetGameIds()
	currentGameAction := a[gameIds[0]]
	new := &ActionsIterator{
		gameActions: a,
		gameIds: a.GetGameIds(),
		currentGameActionDone: false,
		installerNames: currentGameAction.GetInstallerNames(),
		extraNames: currentGameAction.GetExtraNames(),
	}

	return new
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

func (i *ActionsIterator) Next() Action {
	if !i.HasMore() {
		return Action{
			GameId: -1,
			IsFileAction: false,
			FileActionPtr: nil,
			GameAction: "",
		}
	}

	currentGameId := i.gameIds[0]
	if (!i.currentGameActionDone) && i.gameActions[currentGameId].Action == "add" {
		i.currentGameActionDone = true
		return Action{
			GameId: currentGameId,
			IsFileAction: false,
			FileActionPtr: nil,
			GameAction: "add",
		}
	}

	if len(i.installerNames) > 0 {
		name := i.installerNames[0]
		i.installerNames = i.installerNames[1:]
		fileAction := i.gameActions[currentGameId].InstallerActions[name]
		return Action{
			GameId: currentGameId,
			IsFileAction: true,
			FileActionPtr: &fileAction,
			GameAction: "",
		}
	}

	if len(i.extraNames) > 0 {
		name := i.extraNames[0]
		i.extraNames = i.extraNames[1:]
		fileAction := i.gameActions[currentGameId].ExtraActions[name]
		return Action{
			GameId: currentGameId,
			IsFileAction: true,
			FileActionPtr: &fileAction,
			GameAction: "",
		}
	}

	if (!i.currentGameActionDone) && i.gameActions[currentGameId].Action == "remove" {
		i.currentGameActionDone = true
		return Action{
			GameId: currentGameId,
			IsFileAction: false,
			FileActionPtr: nil,
			GameAction: "remove",
		}
	}

	i.gameIds = i.gameIds[1:]
	i.currentGameActionDone = i.gameActions[i.gameIds[0]].Action == "update"
	currentGameAction := i.gameActions[i.gameIds[0]]
	i.installerNames = currentGameAction.GetInstallerNames()
	i.extraNames = currentGameAction.GetExtraNames()
	return i.Next()
}
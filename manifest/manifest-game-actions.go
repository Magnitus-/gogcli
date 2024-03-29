package manifest

import (
	"errors"
	"fmt"
)

type GameActions map[int64]GameAction

const (
	ChecksumValidation string   = "mandatory"
	ChecksumValidationIfPresent = "ifPresent"
	ChecksumNoValidation        = "noValidation"
)

func (g *GameActions) Update(n *GameActions) error {
	for gameId, _ := range *n {
		if gameAction, ok := (*g)[gameId]; ok {
			newGameAction := (*n)[gameId]
			err := gameAction.Update(&newGameAction)
			if err != nil {
				return err
			}
			(*g)[gameId] = gameAction
		} else {
			(*g)[gameId] = (*n)[gameId]
		}
	}

	return nil
}

func (g *GameActions) ActionsLeft() int {
	actionsCount := 0
	for id, _ := range *g {
		gameAction := (*g)[id]
		actionsCount += gameAction.ActionsLeft()
	}
	return actionsCount
}

func (g *GameActions) ApplyAction(a Action) {
	gameId := a.Game.Id
	if a.IsFileAction {
		game := (*g)[gameId]
		if (*a.FileActionPtr).Kind == "installer" {
			delete(game.InstallerActions, (*a.FileActionPtr).Name)
		} else {
			delete(game.ExtraActions, (*a.FileActionPtr).Name)
		}
		if game.IsNoOp() {
			delete((*g), game.Id)
		} else {
			(*g)[gameId] = game
		}
	} else {
		if a.GameAction == "add" {
			game := (*g)[gameId]
			game.Action = "update"
			(*g)[gameId] = game
		} else {
			delete((*g), gameId)
		}
	}
}

func (g *GameActions) GetGameIds() []int64 {
	gameIds := make([]int64, len(*g))

	idx := 0
	for id, _ := range *g {
		gameIds[idx] = id
		idx++
	}

	return gameIds
}

func (g *GameActions) DeepCopy() *GameActions {
	new := GameActions(make(map[int64]GameAction))

	for id, _ := range *g {
		newGame := GameAction{
			Title:            (*g)[id].Title,
			Slug:             (*g)[id].Slug,
			Id:               (*g)[id].Id,
			Action:           (*g)[id].Action,
			InstallerActions: make(map[string]FileAction),
			ExtraActions:     make(map[string]FileAction),
		}

		for name, inst := range (*g)[id].InstallerActions {
			newGame.InstallerActions[name] = inst
		}

		for name, extr := range (*g)[id].ExtraActions {
			newGame.ExtraActions[name] = extr
		}

		new[id] = newGame
	}

	return &new
}

func planManifestGameAddOrRemove(m *ManifestGame, action string) (GameAction, error) {
	g := GameAction{
		Title:            (*m).Title,
		Slug:             (*m).Slug,
		Id:               (*m).Id,
		Action:           action,
		InstallerActions: make(map[string]FileAction),
		ExtraActions:     make(map[string]FileAction),
	}

	if action != "add" && action != "remove" {
		return g, errors.New(fmt.Sprintf("action %s is undefined", action))
	}

	for _, i := range (*m).Installers {
		g.InstallerActions[i.Name] = FileAction{
			Title:  i.Title,
			Name:   i.Name,
			Url:    i.Url,
			Kind:   "installer",
			Action: action,
		}
	}

	for _, e := range (*m).Extras {
		g.ExtraActions[e.Name] = FileAction{
			Title:  e.Title,
			Name:   e.Name,
			Url:    e.Url,
			Kind:   "extra",
			Action: action,
		}
	}

	return g, nil
}

func planManifestGameUpdate(curr *ManifestGame, next *ManifestGame, checksumValidation string, ignoreMetadata bool) GameAction {
	g := GameAction{
		Title:            (*curr).Title,
		Slug:             (*curr).Slug,
		Id:               (*curr).Id,
		Action:           "update",
		InstallerActions: make(map[string]FileAction),
		ExtraActions:     make(map[string]FileAction),
	}

	currentInstallers := make(map[string]ManifestGameInstaller)
	futureInstallers := make(map[string]ManifestGameInstaller)

	for _, i := range (*curr).Installers {
		currentInstallers[i.Name] = i
	}

	for _, i := range (*next).Installers {
		futureInstallers[i.Name] = i
	}

	for name, inst := range futureInstallers {
		if val, ok := currentInstallers[name]; ok {
			if !inst.IsEquivalentTo(&val, checksumValidation, ignoreMetadata) {
				//Overwrite
				g.InstallerActions[name] = FileAction{Title: inst.Title, Name: inst.Name, Url: inst.Url, Kind: "installer", Action: "add"}
			}
		} else {
			//Add missing file
			g.InstallerActions[name] = FileAction{Title: inst.Title, Name: inst.Name, Url: inst.Url, Kind: "installer", Action: "add"}
		}
	}

	for name, inst := range currentInstallers {
		if _, ok := futureInstallers[name]; !ok {
			//Remove dangling file
			g.InstallerActions[name] = FileAction{Title: inst.Title, Name: inst.Name, Url: inst.Url, Kind: "installer", Action: "remove"}
		}
	}

	currentExtras := make(map[string]ManifestGameExtra)
	futureExtras := make(map[string]ManifestGameExtra)

	for _, e := range (*curr).Extras {
		currentExtras[e.Name] = e
	}

	for _, e := range (*next).Extras {
		futureExtras[e.Name] = e
	}

	for name, extr := range futureExtras {
		if val, ok := currentExtras[name]; ok {
			if !extr.IsEquivalentTo(&val, checksumValidation, ignoreMetadata) {
				//Overwrite
				g.ExtraActions[name] = FileAction{Title: extr.Title, Name: extr.Name, Url: extr.Url, Kind: "extra", Action: "add"}
			}
		} else {
			//Add missing file
			g.ExtraActions[name] = FileAction{Title: extr.Title, Name: extr.Name, Url: extr.Url, Kind: "extra", Action: "add"}
		}
	}

	for name, extr := range currentExtras {
		if _, ok := futureExtras[name]; !ok {
			//Remove dangling file
			g.ExtraActions[name] = FileAction{Title: extr.Title, Name: extr.Name, Url: extr.Url, Kind: "extra", Action: "remove"}
		}
	}

	return g
}

func (curr *Manifest) Plan(next *Manifest, checksumValidation string, ignoreMetadata bool) *GameActions {
	actions := GameActions(make(map[int64]GameAction))
	currentGames := make(map[int64]ManifestGame)
	futureGames := make(map[int64]ManifestGame)

	for _, g := range (*curr).Games {
		currentGames[g.Id] = g
	}

	for _, g := range (*next).Games {
		futureGames[g.Id] = g
	}

	for id, game := range futureGames {
		if val, ok := currentGames[id]; !ok {
			actions[id], _ = planManifestGameAddOrRemove(&game, "add")
		} else {
			actions[id] = planManifestGameUpdate(&val, &game, checksumValidation, ignoreMetadata)
		}
	}

	for id, game := range currentGames {
		if _, ok := futureGames[id]; !ok {
			actions[id], _ = planManifestGameAddOrRemove(&game, "remove")
		}
	}

	for id, gameAction := range actions {
		if gameAction.ActionsLeft() == 0 {
			delete(actions, id)
		}
	}

	return &actions
}

func (m *Manifest) GetFileActionFileInfo(gameInfo GameInfo, action FileAction) (FileInfo, error) {
	for idx, _ := range (*m).Games {
		if (*m).Games[idx].Id == gameInfo.Id {
			game := (*m).Games[idx]
			if action.Kind == "installer" {
				installer, err := game.GetInstallerNamed(action.Name)
				if err != nil {
					return FileInfo{}, err
				}
				return FileInfo{
					Game:     gameInfo,
					Kind:     "installer",
					Name:     installer.Name,
					Checksum: installer.Checksum,
					Size:     installer.VerifiedSize,
					Url:      installer.Url,
				}, nil
			} else {
				extra, err := game.GetExtraNamed(action.Name)
				if err != nil {
					return FileInfo{}, err
				}
				return FileInfo{
					Game:     gameInfo,
					Kind:     "extra",
					Name:     extra.Name,
					Checksum: extra.Checksum,
					Size:     extra.VerifiedSize,
					Url:      extra.Url,
				}, err
			}
		}
	}

	msg := fmt.Sprintf("*Manifest.GetFileActionFileInfo(gameId=%d, ...) -> Game with given id not found in manifest", gameInfo.Id)
	return FileInfo{}, errors.New(msg)
}

package manifest

import (
	"errors"
	"fmt"
)

type FileAction struct {
	Title        string
	Name         string
	Url          string
	Kind         string
	Action       string
}

type GameAction struct {
	Title            string
	Id               int
	Action           string
	InstallerActions map[string]FileAction
	ExtraActions     map[string]FileAction
}

func (g *GameAction) IsNoOp() bool {
	return (*g).Action == "update" && (!(*g).HasFileActions())
}

func (g *GameAction) GetInstallerNames() []string {
	installerNames := make([]string, len((*g).InstallerActions))
	
	idx := 0
	for name, _ := range (*g).InstallerActions {
		installerNames[idx] = name
		idx++
	}

	return installerNames
}

func (g *GameAction) GetExtraNames() []string {
	extraNames := make([]string, len((*g).ExtraActions))
	
	idx := 0
	for name, _ := range (*g).ExtraActions {
		extraNames[idx] = name
		idx++
	}

	return extraNames
}

func (g *GameAction) HasFileActions() bool {
	return len((*g).InstallerActions) > 0 || len((*g).ExtraActions) > 0
}

func (g *GameAction) ExtractFileAction() (FileAction, string, error) {
	var fetchedAction FileAction
	if len((*g).InstallerActions) > 0 {
		for k, _ := range (*g).InstallerActions {
			fetchedAction = (*g).InstallerActions[k]
			delete((*g).InstallerActions, k)
			return fetchedAction, "installer", nil
		}
	} else if len((*g).ExtraActions) > 0 {
		for k, _ := range (*g).ExtraActions {
			fetchedAction = (*g).ExtraActions[k]
			delete((*g).ExtraActions, k)
			return fetchedAction, "extra", nil
		}
	}

	return fetchedAction, "", errors.New("ExtractFileAction() -> No action left to extract")
}

type GameActions map[int]GameAction

func (g *GameActions) ApplyAction(a Action) {
	if a.IsFileAction {
		game := (*g)[a.GameId]
		if (*a.FileActionPtr).Kind == "installer" { 
			delete(game.InstallerActions, (*a.FileActionPtr).Name)
		} else {
			delete(game.ExtraActions, (*a.FileActionPtr).Name)
		}
		if game.IsNoOp() {
			delete((*g), game.Id)
		} else {
			(*g)[a.GameId] = game
		}
	} else {
		if a.GameAction == "add" {
			game := (*g)[a.GameId]
			game.Action = "update"
			(*g)[a.GameId] = game
		} else {
			delete((*g), a.GameId)
		}
	}
}

func (g *GameActions) GetGameIds() []int {
	gameIds := make([]int, len(*g))

	idx := 0
	for id, _ := range *g {
		gameIds[idx] = id
		idx++
	}

	return gameIds
}

func (g *GameActions) DeepCopy() (*GameActions) {
	new := GameActions(make(map[int]GameAction))

	for id, _ := range (*g) {
		newGame := GameAction{
			Title: (*g)[id].Title,
			Id: (*g)[id].Id,
			Action: (*g)[id].Action,
			InstallerActions: make(map[string]FileAction),
			ExtraActions: make(map[string]FileAction),
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

func planManifestGameUpdate(curr *ManifestGame, next *ManifestGame) GameAction {
	g := GameAction{
		Title:            (*curr).Title,
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
			if !inst.isEquivalentTo(&val) {
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
			if !extr.isEquivalentTo(&val) {
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

func (curr *Manifest) Plan(next *Manifest) *GameActions {
	actions := GameActions(make(map[int]GameAction))
	currentGames := make(map[int]ManifestGame)
	futureGames := make(map[int]ManifestGame)

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
			actions[id] = planManifestGameUpdate(&val, &game)
		}
	}

	for id, game := range currentGames {
		if _, ok := futureGames[id]; !ok {
			actions[id], _ = planManifestGameAddOrRemove(&game, "remove")
		}
	}

	return &actions
}
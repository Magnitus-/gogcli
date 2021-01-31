package manifest

import (
	"errors"
	"fmt"
)

type FileAction struct {
	Title        string
	PreviousName string
	Name         string
	Url          string
	Action       string
}

type GameAction struct {
	Title            string
	Id               int
	Action           string
	InstallerActions map[string]FileAction
	ExtraActions     map[string]FileAction
}

type GameActions map[int]GameAction

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
		g.InstallerActions[i.Url] = FileAction{
			Title:  i.Title,
			Name:   i.Name,
			Url:    i.Url,
			Action: action,
		}
	}

	for _, e := range (*m).Extras {
		g.ExtraActions[e.Url] = FileAction{
			Title:  e.Title,
			Name:   e.Name,
			Url:    e.Url,
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
		currentInstallers[i.Url] = i
	}

	for _, i := range (*next).Installers {
		futureInstallers[i.Url] = i
	}

	for url, inst := range futureInstallers {
		if val, ok := currentInstallers[url]; ok {
			if !inst.isEquivalentTo(&val) {
				//Overwrite
				g.InstallerActions[url] = FileAction{Title: inst.Title, PreviousName: val.Name, Name: inst.Name, Url: inst.Url, Action: "replace"}
			}
		} else {
			//Add missing file
			g.InstallerActions[url] = FileAction{Title: inst.Title, Name: inst.Name, Url: inst.Url, Action: "add"}
		}
	}

	for url, inst := range currentInstallers {
		if _, ok := futureInstallers[url]; !ok {
			//Remove dangling file
			g.InstallerActions[url] = FileAction{Title: inst.Title, Name: inst.Name, Url: inst.Url, Action: "remove"}
		}
	}

	currentExtras := make(map[string]ManifestGameExtra)
	futureExtras := make(map[string]ManifestGameExtra)

	for _, e := range (*curr).Extras {
		currentExtras[e.Url] = e
	}

	for _, e := range (*next).Extras {
		futureExtras[e.Url] = e
	}

	for url, extr := range futureExtras {
		if val, ok := currentExtras[url]; ok {
			if !extr.isEquivalentTo(&val) {
				//Overwrite
				g.ExtraActions[url] = FileAction{Title: extr.Title, PreviousName: val.Name, Name: extr.Name, Url: extr.Url, Action: "replace"}
			}
		} else {
			//Add missing file
			g.ExtraActions[url] = FileAction{Title: extr.Title, Name: extr.Name, Url: extr.Url, Action: "add"}
		}
	}

	for url, extr := range currentExtras {
		if _, ok := futureExtras[url]; !ok {
			//Remove dangling file
			g.ExtraActions[url] = FileAction{Title: extr.Title, Name: extr.Name, Url: extr.Url, Action: "remove"}
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
package manifest

type FileAction struct {
	Title        string
	Name         string
	Url          string
	Kind         string
	Action       string
}

type GameAction struct {
	Title            string
	Id               int64
	Action           string
	InstallerActions map[string]FileAction
	ExtraActions     map[string]FileAction
}

func (g *GameAction) IsNoOp() bool {
	return (*g).ActionsLeft() == 0
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

func (g *GameAction) CountFileActions() int {
	return len((*g).InstallerActions) + len((*g).ExtraActions)
}

func (g *GameAction) ActionsLeft() int {
	actionsCount := g.CountFileActions()
	if (*g).Action != "update" {
		actionsCount++
	}
	return actionsCount
}
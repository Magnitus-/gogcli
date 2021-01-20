package manifest

type FileAction struct {
	Name   string
	Url    string
	Action string
}

type GameAction struct {
	Id               int
	Action           string
	InstallerActions map[string]FileAction
	ExtraActions     map[string]FileAction
}

type GameActions map[int]GameAction

package manifest

type ProtectedGameFiles struct {
	Installers []string
	Extras     []string
}

func (p *ProtectedGameFiles) AddGameFile(file FileInfo) {
	if file.Kind == "extra" {
		if !containsStr((*p).Extras, file.Name) {
			(*p).Extras = append((*p).Extras, file.Name)
		}
		return
	}

	if !containsStr((*p).Installers, file.Name) {
		(*p).Extras = append((*p).Installers, file.Name)
	}
}

func (p *ProtectedGameFiles) RemoveGameFile(file FileInfo) {
	if file.Kind == "extra" {
		(*p).Extras = RemoveStrFromList((*p).Extras, file.Name)
		return
	}

	(*p).Installers = RemoveStrFromList((*p).Installers, file.Name)
}

type ProtectedManifestFiles map[int64]ProtectedGameFiles

func (p *ProtectedManifestFiles) AddGameFile(file FileInfo) {
	if (*p) == nil {
		(*p) = make(map[int64]ProtectedGameFiles)
	}

	game, ok := (*p)[file.Game.Id]
	if !ok {
		game = ProtectedGameFiles{
			Installers: []string{},
			Extras:     []string{},
		}
	}
	
	game.AddGameFile(file)
	(*p)[file.Game.Id] = game
}

func (p *ProtectedManifestFiles) RemoveGameFile(file FileInfo) {
	if (*p) == nil {
		(*p) = make(map[int64]ProtectedGameFiles)
	}
	
	game, ok := (*p)[file.Game.Id]
	if !ok {
		return
	}

	game.RemoveGameFile(file)
	(*p)[file.Game.Id] = game
}
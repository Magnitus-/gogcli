package manifest

type ProtectedGameFiles struct {
	Installers []string
	Extras     []string
}

type ProtectedManifestFiles map[int64]ProtectedGameFiles

func (p *ProtectedManifestFiles) AddGameFile(file FileInfo) {
	/*(*p)[file.Game.Id] = GameFileGuard{
		GameId: file.Game.Id,
		FileType: file.Kind,
		FileName: file.Name,
	}*/
}

func (p *ProtectedManifestFiles) RemoveGameFile(file FileInfo) {
	delete(*p, file.Game.Id)
}
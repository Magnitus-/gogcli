package manifest

type FileInfo struct {
	GameId   int
	Kind     string
	Name     string
	Checksum string
	Size     int
	Url      string
}

type ManifestFileIterator struct {
	manifestPtr      *Manifest
	currentGame      int
	currentInstaller int
	currentExtra     int
}

func NewManifestFileInterator(m *Manifest) ManifestFileIterator {
    new := ManifestFileIterator{
		manifestPtr: m,
		currentGame: 0,
		currentInstaller: 0,
		currentExtra: 0,
	}

	return new
}

func (i *ManifestFileIterator) HasMore() bool {
	//Not the last game
	if (*i).currentGame < len((*(*i).manifestPtr).Games) - 1 {
		return true
	}

	notLastInstaller := (*i).currentInstaller < len((*(*i).manifestPtr).Games[(*i).currentGame].Installers)
	notLastExtra := (*i).currentExtra < len((*(*i).manifestPtr).Games[(*i).currentGame].Extras)
	return notLastInstaller || notLastExtra
}

func (i *ManifestFileIterator) Next() FileInfo {
	if !i.HasMore() {
		return FileInfo{
			GameId: -1,
			Kind: "",
			Name: "",
			Checksum: "",
			Size: 0,
			Url: "",
		}
	}

	currentGame := (*(*i).manifestPtr).Games[(*i).currentGame]
	if (*i).currentInstaller < len(currentGame.Installers) {
		new := FileInfo{
			GameId: currentGame.Id,
			Kind: "installer",
			Name: currentGame.Installers[(*i).currentInstaller].Name,
			Checksum: currentGame.Installers[(*i).currentInstaller].Checksum,
			Size: currentGame.Installers[(*i).currentInstaller].VerifiedSize,
			Url: currentGame.Installers[(*i).currentInstaller].Url,
		}
		(*i).currentInstaller++
		return new
	} else if (*i).currentExtra < len(currentGame.Extras) {
		new := FileInfo{
			GameId: currentGame.Id,
			Kind: "extra",
			Name: currentGame.Extras[(*i).currentExtra].Name,
			Checksum: currentGame.Extras[(*i).currentExtra].Checksum,
			Size: currentGame.Extras[(*i).currentExtra].VerifiedSize,
			Url: currentGame.Extras[(*i).currentExtra].Url,
		}
		(*i).currentExtra++
		return new
	} else {
		(*i).currentGame++
		return i.Next()
	}
}
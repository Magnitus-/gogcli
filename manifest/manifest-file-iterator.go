package manifest

import "errors"

type GameInfo struct {
	Id int64
	Slug          string
	Title         string
}

type FileInfo struct {
	Game     GameInfo
	Kind     string
	Name     string
	Checksum string
	Size     int64
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
		manifestPtr:      m,
		currentGame:      0,
		currentInstaller: 0,
		currentExtra:     0,
	}

	return new
}

func (i *ManifestFileIterator) HasMore() bool {
	//Not the last game
	if (*i).currentGame < (len((*(*i).manifestPtr).Games) - 1) {
		return true
	}

	//If its the last game, check to see if there is an installer or extra left to fetch
	currentGame := (*(*i).manifestPtr).Games[(*i).currentGame]
	notLastInstaller := (*i).currentInstaller < len(currentGame.Installers)
	notLastExtra := (*i).currentExtra < len(currentGame.Extras)
	return notLastInstaller || notLastExtra
}

func (i *ManifestFileIterator) Next() (FileInfo, error) {
	if !i.HasMore() {
		return FileInfo{
			Game:     GameInfo{Id: -1, Slug: "", Title: ""},
			Kind:     "",
			Name:     "",
			Checksum: "",
			Size:     0,
			Url:      "",
		}, errors.New("*ManifestFileIterator.Next() -> End of iterator, cannot fetch anymore")
	}

	currentGame := (*(*i).manifestPtr).Games[(*i).currentGame]
	if (*i).currentInstaller < len(currentGame.Installers) {
		new := FileInfo{
			Game:     GameInfo{Id: currentGame.Id, Slug: currentGame.Slug, Title: currentGame.Title},
			Kind:     "installer",
			Name:     currentGame.Installers[(*i).currentInstaller].Name,
			Checksum: currentGame.Installers[(*i).currentInstaller].Checksum,
			Size:     currentGame.Installers[(*i).currentInstaller].VerifiedSize,
			Url:      currentGame.Installers[(*i).currentInstaller].Url,
		}
		(*i).currentInstaller++
		return new, nil
	} else if (*i).currentExtra < len(currentGame.Extras) {
		new := FileInfo{
			Game:   GameInfo{Id: currentGame.Id, Slug: currentGame.Slug, Title: currentGame.Title},
			Kind:     "extra",
			Name:     currentGame.Extras[(*i).currentExtra].Name,
			Checksum: currentGame.Extras[(*i).currentExtra].Checksum,
			Size:     currentGame.Extras[(*i).currentExtra].VerifiedSize,
			Url:      currentGame.Extras[(*i).currentExtra].Url,
		}
		(*i).currentExtra++
		return new, nil
	} else {
		(*i).currentGame++
		(*i).currentInstaller = 0
		(*i).currentExtra = 0
		return i.Next()
	}
}

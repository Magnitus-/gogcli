package manifest

import (
	"errors"
	"fmt"
	"strings"
)

type ManifestGameExtra struct {
	Url           string
	Title         string
	Name          string
	Type          string
	Info          int
	EstimatedSize string
	VerifiedSize  int
	Checksum      string
}

func (e *ManifestGameExtra) hasOneOfTypeTerms(typeTerms []string) bool {
	for _, t := range typeTerms {
		if strings.Contains((*e).Type, t) {
			return true
		}
	}
	return false
}

func (e *ManifestGameExtra) isEquivalentTo(o *ManifestGameExtra) bool {
	sameName := (*e).Name == (*o).Name
	sameTitle := (*e).Title == (*o).Title
	sameUrl := (*e).Url == (*o).Url
	sameVerifiedSize := (*o).VerifiedSize != 0 && (*e).VerifiedSize == (*o).VerifiedSize
	sameChecksum := (*o).Checksum != "" && (*e).Checksum == (*o).Checksum
	return sameName && sameTitle && sameUrl && sameVerifiedSize && sameChecksum
}

func (e *ManifestGameExtra) getEstimatedSizeInBytes() (int, error) {
	return GetEstimateToBytes((*e).EstimatedSize)
}

type ManifestGameInstaller struct {
	Language      string
	Os            string
	Url           string
	Title         string
	Name          string
	Version       string
	Date          string
	EstimatedSize string
	VerifiedSize  int
	Checksum      string
}

func (i *ManifestGameInstaller) hasOneOfOses(oses []string) bool {
	for _, os := range oses {
		if os == i.Os {
			return true
		}
	}
	return false
}

func (i *ManifestGameInstaller) hasOneOfLanguages(languages []string) bool {
	for _, l := range languages {
		if l == i.Language {
			return true
		}
	}
	return false
}

func (i *ManifestGameInstaller) isEquivalentTo(o *ManifestGameInstaller) bool {
	sameName := (*i).Name == (*o).Name
	sameTitle := (*i).Title == (*o).Title
	sameUrl := (*i).Url == (*o).Url
	sameVerifiedSize := (*o).VerifiedSize != 0 && (*i).VerifiedSize == (*o).VerifiedSize
	sameChecksum := (*o).Checksum != "" && (*i).Checksum == (*o).Checksum
	return sameName && sameTitle && sameUrl && sameVerifiedSize && sameChecksum
}

func (i *ManifestGameInstaller) getEstimatedSizeInBytes() (int, error) {
	return GetEstimateToBytes((*i).EstimatedSize)
}

type ManifestGame struct {
	Id            int
	Title         string
	CdKey         string
	Tags          []string
	Installers    []ManifestGameInstaller
	Extras        []ManifestGameExtra
	EstimatedSize string
	VerifiedSize  int
}

func (g *ManifestGame) trimInstallers(oses []string, languages []string, keepAny bool) {
	filteredInstallers := make([]ManifestGameInstaller, 0)

	if keepAny {
		if len(oses) == 0 && len(languages) == 0 {
			//Save some needless computation
			return
		}

		for _, i := range (*g).Installers {
			hasOneOfOses := len(oses) == 0 || i.hasOneOfOses(oses)
			hasOneOfLanguages := len(languages) == 0 || i.hasOneOfLanguages(languages)
			if hasOneOfOses && hasOneOfLanguages {
				filteredInstallers = append(filteredInstallers, i)
			}
		}
	}
	(*g).Installers = filteredInstallers
}

func (g *ManifestGame) trimExtras(typeTerms []string, keepAny bool) {
	filteredExtras := make([]ManifestGameExtra, 0)

	if keepAny {
		if len(typeTerms) == 0 {
			return
		}

		for _, e := range (*g).Extras {
			if e.hasOneOfTypeTerms(typeTerms) {
				filteredExtras = append(filteredExtras, e)
			}
		}
	}
	(*g).Extras = filteredExtras
}

func (g *ManifestGame) hasTitleTerm(titleTerm string) bool {
	return titleTerm == "" || strings.Contains((*g).Title, titleTerm)
}

func (g *ManifestGame) hasOneOfTags(tags []string) bool {
	for _, t := range tags {
		for _, gt := range (*g).Tags {
			if t == gt {
				return true
			}
		}
	}
	return false
}

func (g *ManifestGame) isEmpty() bool {
	return len((*g).Installers) == 0 && len((*g).Extras) == 0
}

func (g *ManifestGame) computeEstimatedSize() (int, error) {
	accumulate := 0
	for _, inst := range (*g).Installers {
		size, err := inst.getEstimatedSizeInBytes()
		if err != nil {
			return 0, err
		}
		accumulate += size
	}

	for _, extr := range (*g).Extras {
		size, err := extr.getEstimatedSizeInBytes()
		if err != nil {
			return 0, err
		}
		accumulate += size
	}

	(*g).EstimatedSize = GetBytesToEstimate(accumulate)
	return accumulate, nil
}

func (g *ManifestGame) fillMissingFileInfo(fileKind string, fileUrl string, fileName string, fileSize int, fileChecksum string) error {
	if fileKind == "installer" {
		for idx, _ := range (*g).Installers {
			if (*g).Installers[idx].Url == fileUrl {
				(*g).Installers[idx].Name = fileName
				(*g).Installers[idx].VerifiedSize = fileSize
				(*g).Installers[idx].Checksum = fileChecksum
				return nil
			}
		}

		return errors.New(fmt.Sprintf("File with url %s was not found in the installers of game with id %d", fileUrl, (*g).Id))
	} else if fileKind == "extra" {
		for idx, _ := range (*g).Extras {
			if (*g).Extras[idx].Url == fileUrl {
				(*g).Extras[idx].Name = fileName
				(*g).Extras[idx].VerifiedSize = fileSize
				(*g).Extras[idx].Checksum = fileChecksum
				return nil
			}
		}

		return errors.New(fmt.Sprintf("File with url %s was not found in the extras of game with id %d", fileUrl, (*g).Id))
	}

	return errors.New(fmt.Sprintf("%s is not a valid kind of file", fileKind))
}

type Manifest struct {
	Games         []ManifestGame
	EstimatedSize string
	VerifiedSize  int
}

func NewEmptyManifest() *Manifest {
	return &Manifest{Games: make([]ManifestGame, 0), EstimatedSize: "0 MB"}
}

func (m *Manifest) TrimGames(titleTerm string, tags []string) {
	filteredGames := make([]ManifestGame, 0)

	if titleTerm == "" && len(tags) == 0 {
		//Save some needless computation
		return
	}

	for _, g := range (*m).Games {
		hasTitleTerm := titleTerm == "" || g.hasTitleTerm(titleTerm)
		hasOneOfTags := len(tags) == 0 || g.hasOneOfTags(tags)
		if hasTitleTerm && hasOneOfTags {
			filteredGames = append(filteredGames, g)
		}
	}

	(*m).Games = filteredGames
}

func (m *Manifest) TrimInstallers(oses []string, languages []string, keepAny bool) {
	filteredGames := make([]ManifestGame, 0)

	if len(oses) == 0 && len(languages) == 0 && keepAny {
		//Save some needless computation
		return
	}

	for _, g := range (*m).Games {
		g.trimInstallers(oses, languages, keepAny)
		if !g.isEmpty() {
			filteredGames = append(filteredGames, g)
		}
	}

	(*m).Games = filteredGames
}

func (m *Manifest) TrimExtras(typeTerms []string, keepAny bool) {
	filteredGames := make([]ManifestGame, 0)

	if len(typeTerms) == 0 && keepAny {
		//Save some needless computation
		return
	}

	for _, g := range (*m).Games {
		g.trimExtras(typeTerms, keepAny)
		if !g.isEmpty() {
			filteredGames = append(filteredGames, g)
		}
	}

	(*m).Games = filteredGames
}

//For new games
func (m *Manifest) AddGames(games []ManifestGame) {
	(*m).Games = append((*m).Games, games...)
}

//For game updates
func (m *Manifest) ReplaceGames(games []ManifestGame) {
	filteredGames := make([]ManifestGame, 0)
	replaceMap := make(map[int]ManifestGame)

	for _, game := range games {
		replaceMap[game.Id] = game
	}

	for _, game := range (*m).Games {
		if repl, ok := replaceMap[game.Id]; ok {
			filteredGames = append(filteredGames, repl)
		} else {
			filteredGames = append(filteredGames, game)
		}
	}

	(*m).Games = filteredGames
}

func (m *Manifest) ComputeEstimatedSize() (int, error) {
	accumulate := 0

	for idx, _ := range (*m).Games {
		size, err := (*m).Games[idx].computeEstimatedSize()
		if err != nil {
			return 0, err
		}
		accumulate += size
	}

	(*m).EstimatedSize = GetBytesToEstimate(accumulate)
	return accumulate, nil
}

func (m *Manifest) FillMissingFileInfo(gameId int, fileKind string, fileUrl string, fileName string, fileSize int, fileChecksum string) error {
	fn := fmt.Sprintf(
		"Manifest.FillMissingFileInfo(gameId=%d, fileKind=%s, fileUrl=%s, fileName=%s, fileSize=%d, fileChecksum=%s)",
		gameId,
		fileKind,
		fileUrl,
		fileName,
		fileSize,
		fileChecksum,
	)
	for idx, _ := range (*m).Games {
		if (*m).Games[idx].Id == gameId {
			err := (*m).Games[idx].fillMissingFileInfo(fileKind, fileUrl, fileName, fileSize, fileChecksum)
			if err != nil {
				return errors.New(fmt.Sprintf("%s -> Error filling game's missing info: %s", err.Error()))
			}
			return nil
		}
	}

	return errors.New(fmt.Sprintf("%s -> Provided game id could not be found in the manifest", fn))
}
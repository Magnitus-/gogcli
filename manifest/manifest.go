package manifest

import (
	"errors"
	"fmt"
)

type Manifest struct {
	Games         []ManifestGame
	EstimatedSize string
	VerifiedSize  int64
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

func (m *Manifest) ComputeEstimatedSize() (int64, error) {
	accumulate := int64(0)

	for idx, _ := range (*m).Games {
		size, err := (*m).Games[idx].computeEstimatedSize()
		if err != nil {
			return int64(0), err
		}
		accumulate += size
	}

	(*m).EstimatedSize = GetBytesToEstimate(accumulate)
	return accumulate, nil
}

func (m *Manifest) FillMissingFileInfo(gameId int, fileKind string, fileName string, fileSize int64, fileChecksum string) error {
	fn := fmt.Sprintf(
		"Manifest.FillMissingFileInfo(gameId=%d, fileKind=%s, fileName=%s, fileSize=%d, fileChecksum=%s)",
		gameId,
		fileKind,
		fileName,
		fileSize,
		fileChecksum,
	)
	for idx, _ := range (*m).Games {
		if (*m).Games[idx].Id == gameId {
			err := (*m).Games[idx].fillMissingFileInfo(fileKind, fileName, fileSize, fileChecksum)
			if err != nil {
				return errors.New(fmt.Sprintf("%s -> Error filling game's missing info: %s", err.Error()))
			}
			return nil
		}
	}

	return errors.New(fmt.Sprintf("%s -> Provided game id could not be found in the manifest", fn))
}

func (m *Manifest) GetUrlMappedInstallers() map[string]*ManifestGameInstaller {
	installers := make(map[string]*ManifestGameInstaller)

	for idx, _ := range (*m).Games {
		for idx2, _ := range (*m).Games[idx].Installers {
			installers[(*m).Games[idx].Installers[idx2].Url] = &(*m).Games[idx].Installers[idx2]
		}
	}
	
	return installers
}

func (m *Manifest) GetUrlMappedExtras() map[string]*ManifestGameExtra {
	extras := make(map[string]*ManifestGameExtra)

	for idx, _ := range (*m).Games {
		for idx2, _ := range (*m).Games[idx].Extras {
			extras[(*m).Games[idx].Extras[idx2].Url] = &(*m).Games[idx].Extras[idx2]
		}
	}
	
	return extras
}
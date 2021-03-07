package manifest

import (
	"errors"
	"fmt"
)

type Manifest struct {
	Games         []ManifestGame
	EstimatedSize string
	VerifiedSize  int64
	Filter ManifestFilter
}

func (m *Manifest) ImprintMissingChecksums(prev *Manifest) error {
	prevGames := make(map[int64]ManifestGame)

	for _, game := range (*prev).Games {
		prevGames[game.Id] = game
	}

	for idx, game := range (*m).Games {
		if prevGame, ok := prevGames[game.Id]; ok {
			err := game.ImprintMissingChecksums(&prevGame)
			if err != nil {
				return err
			}
			(*m).Games[idx] = game
		}
	}

	return nil
}

func (m *Manifest) Trim() {
	m.TrimGames()	
	m.TrimInstallers()	
	m.TrimExtras()
}

func NewEmptyManifest(f ManifestFilter) *Manifest {
	return &Manifest{
		Games: make([]ManifestGame, 0), 
		EstimatedSize: "0 MB",
		VerifiedSize: 0,
		Filter: f,
	}
}

func (m *Manifest) Finalize() {
	filteredGames := make([]ManifestGame, 0)
	for _, g := range (*m).Games {
		if !g.IsEmpty() {
			filteredGames = append(filteredGames, g)
		}
	}
	(*m).Games = filteredGames

	m.ComputeEstimatedSize()	
	m.ComputeVerifiedSize()
}

func (m *Manifest) TrimGames() {
	titles := (*m).Filter.Titles
	tags := (*m).Filter.Tags
	filteredGames := make([]ManifestGame, 0)

	if len(titles) == 0 && len(tags) == 0 {
		//Save some needless computation
		return
	}

	for _, g := range (*m).Games {
		hasTitleTerm := len(titles) == 0 || g.HasTitleTerms(titles)
		hasOneOfTags := len(tags) == 0 || g.HasOneOfTags(tags)
		if hasTitleTerm && hasOneOfTags {
			filteredGames = append(filteredGames, g)
		}
	}

	(*m).Games = filteredGames
}

func (m *Manifest) TrimInstallers() {
	oses := (*m).Filter.Oses
	languages := (*m).Filter.Languages
	keepAny := (*m).Filter.Installers
	filteredGames := make([]ManifestGame, 0)

	if len(oses) == 0 && len(languages) == 0 && keepAny {
		//Save some needless computation
		return
	}

	for _, g := range (*m).Games {
		g.TrimInstallers(oses, languages, keepAny)
		if !g.IsEmpty() {
			filteredGames = append(filteredGames, g)
		}
	}

	(*m).Games = filteredGames
}

func (m *Manifest) TrimExtras() {
	typeTerms := (*m).Filter.ExtraTypes
	keepAny := (*m).Filter.Extras
	filteredGames := make([]ManifestGame, 0)

	if len(typeTerms) == 0 && keepAny {
		//Save some needless computation
		return
	}

	for _, g := range (*m).Games {
		g.TrimExtras(typeTerms, keepAny)
		if !g.IsEmpty() {
			filteredGames = append(filteredGames, g)
		}
	}

	(*m).Games = filteredGames
}

func (m *Manifest) OverwriteGames(games []ManifestGame) {
	filteredGames := make([]ManifestGame, 0)
	replaceMap := make(map[int64]ManifestGame)
	existingMap := make(map[int64]ManifestGame)

	for _, game := range games {
		replaceMap[game.Id] = game
	}

	for _, game := range (*m).Games {
		existingMap[game.Id] = game
		if repl, ok := replaceMap[game.Id]; ok {
			filteredGames = append(filteredGames, repl)
		} else {
			filteredGames = append(filteredGames, game)
		}
	}

	for _, game := range games {
		if _, ok := existingMap[game.Id]; !ok {
			filteredGames = append(filteredGames, game)
		}
	}

	(*m).Games = filteredGames
}

func (m *Manifest) ComputeVerifiedSize() int64 {
	accumulate := int64(0)

	for idx, _ := range (*m).Games {
		accumulate += (*m).Games[idx].ComputeVerifiedSize()
	}

	(*m).VerifiedSize = accumulate
	return accumulate
}

func (m *Manifest) ComputeEstimatedSize() (int64, error) {
	accumulate := int64(0)

	for idx, _ := range (*m).Games {
		size, err := (*m).Games[idx].ComputeEstimatedSize()
		if err != nil {
			return int64(0), err
		}
		accumulate += size
	}

	(*m).EstimatedSize = GetBytesToEstimate(accumulate)
	return accumulate, nil
}

func (m *Manifest) FillMissingFileInfo(gameId int64, fileKind string, fileName string, fileSize int64, fileChecksum string) error {
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
			err := (*m).Games[idx].FillMissingFileInfo(fileKind, fileName, fileSize, fileChecksum)
			if err != nil {
				return errors.New(fmt.Sprintf("%s -> Error filling game's missing info: %s", fn, err.Error()))
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
package manifest

import "strings"

type ManifestGameExtra struct {
	Url          string
	Name         string
	Type         string
	Info         int
	Size         string
	VerifiedSize int
	Checksum     string
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
	sameUrl := (*e).Url == (*o).Url
	sameVerifiedSize := (*o).VerifiedSize != 0 && (*e).VerifiedSize == (*o).VerifiedSize
	sameChecksum := (*o).Checksum != "" && (*e).Checksum == (*o).Checksum
	return sameName && sameUrl && sameVerifiedSize && sameChecksum
}

type ManifestGameInstaller struct {
	Language     string
	Os           string
	Url          string
	Name         string
	Version      string
	Date         string
	Size         string
	VerifiedSize int
	Checksum     string
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
	sameUrl := (*i).Url == (*o).Url
	sameVerifiedSize := (*o).VerifiedSize != 0 && (*i).VerifiedSize == (*o).VerifiedSize
	sameChecksum := (*o).Checksum != "" && (*i).Checksum == (*o).Checksum
	return sameName && sameUrl && sameVerifiedSize && sameChecksum
}

type ManifestGame struct {
	Id         int
	Title      string
	CdKey      string
	Tags       []string
	Installers []ManifestGameInstaller
	Extras     []ManifestGameExtra
	Size       string
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

type Manifest struct {
	Games []ManifestGame
	Size  string
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

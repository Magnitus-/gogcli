package manifest

import "strings"

type ManifestGameExtra struct {
	Url          string
	Name         string
	Type         string
	Info         int
	Size         string
	ComputedSize int
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

type ManifestGameInstaller struct {
	Language     string
	Os           string
	Url          string
	Name         string
	Version      string
	Date         string
	Size         string
	ComputedSize int
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
	var filteredInstallers []ManifestGameInstaller
	if keepAny {
		if len(oses) == 0 && len(languages) == 0 {
			//Save some needless computation
			return
		}

		for _, i := range (*g).Installers {
			hasOneOfOses := len(oses) == 0 || i.hasOneOfLanguages(oses)
			hasOneOfLanguages := len(languages) == 0 || i.hasOneOfLanguages(languages)
			if hasOneOfOses && hasOneOfLanguages {
				filteredInstallers = append(filteredInstallers, i)
			}
		}
	}
	(*g).Installers = filteredInstallers
}

func (g *ManifestGame) trimExtras(typeTerms []string, keepAny bool) {
	var filteredExtras []ManifestGameExtra
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
	return len((*g).Installers) > 0 || len((*g).Extras) > 0
}

type Manifest struct {
	Games []ManifestGame
	Size  string
}

func (m *Manifest) TrimGames(titleTerm string, tags []string) {
	var filteredGames []ManifestGame

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
	var filteredGames []ManifestGame

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
	var filteredGames []ManifestGame

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

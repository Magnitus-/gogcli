package manifest

import (
	"errors"
	"fmt"
	"strings"
)

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
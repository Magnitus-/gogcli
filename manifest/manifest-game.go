package manifest

import (
	"errors"
	"fmt"
	"strings"
)

type ManifestGame struct {
	Id            int64
	Title         string
	CdKey         string
	Tags          []string
	Installers    []ManifestGameInstaller
	Extras        []ManifestGameExtra
	EstimatedSize string
	VerifiedSize  int64
}

func (g *ManifestGame) getInstallerNamed(name string) (ManifestGameInstaller, error) {
	for idx, _ := range (*g).Installers {
		if (*g).Installers[idx].Name == name {
			return (*g).Installers[idx], nil
		}
	}

	msg := fmt.Sprintf("*ManifestGame.getInstallerNamed(name=%) -> No installer by that name", name)
	return ManifestGameInstaller{}, errors.New(msg)
}

func (g *ManifestGame) getExtraNamed(name string) (ManifestGameExtra, error) {
	for idx, _ := range (*g).Extras {
		if (*g).Extras[idx].Name == name {
			return (*g).Extras[idx], nil
		}
	}

	msg := fmt.Sprintf("*ManifestGame.getExtraNamed(name=%) -> No extra by that name", name)
	return ManifestGameExtra{}, errors.New(msg)
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

func (g *ManifestGame) computeVerifiedSize() int64 {
	accumulate := int64(0)
	for _, inst := range (*g).Installers {
		accumulate += inst.VerifiedSize
	}

	for _, extr := range (*g).Extras {
		accumulate += extr.VerifiedSize
	}

	(*g).VerifiedSize = accumulate
	return accumulate
}

func (g *ManifestGame) computeEstimatedSize() (int64, error) {
	accumulate := int64(0)
	for _, inst := range (*g).Installers {
		size, err := inst.getEstimatedSizeInBytes()
		if err != nil {
			return int64(0), err
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

func (g *ManifestGame) fillMissingFileInfo(fileKind string, fileName string, fileSize int64, fileChecksum string) error {
	if fileKind == "installer" {
		for idx, _ := range (*g).Installers {
			if (*g).Installers[idx].Name == fileName {
				(*g).Installers[idx].VerifiedSize = fileSize
				(*g).Installers[idx].Checksum = fileChecksum
				return nil
			}
		}

		return errors.New(fmt.Sprintf("File with name %s was not found in the installers of game with id %d", fileName, (*g).Id))
	} else if fileKind == "extra" {
		for idx, _ := range (*g).Extras {
			if (*g).Extras[idx].Name == fileName {
				(*g).Extras[idx].VerifiedSize = fileSize
				(*g).Extras[idx].Checksum = fileChecksum
				return nil
			}
		}

		return errors.New(fmt.Sprintf("File with name %s was not found in the extras of game with id %d", fileName, (*g).Id))
	}

	return errors.New(fmt.Sprintf("%s is not a valid kind of file", fileKind))
}
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

func (g *ManifestGame) ImprintMissingChecksums(prev *ManifestGame) error {
	if (*g).Id != (*prev).Id {
		return errors.New("imprintMissingChecksums(...) -> Game ids do not match")
	}

	previousInstallers := make(map[string]ManifestGameInstaller)
	previousExtras := make(map[string]ManifestGameExtra)

	for _, installer := range (*prev).Installers {
		previousInstallers[installer.Name] = installer 
	}

	for _, extra := range (*prev).Extras {
		previousExtras[extra.Name] = extra
	}

	for idx, installer := range (*g).Installers {
		if prevInstaller, ok := previousInstallers[installer.Name]; ok {
			if installer.IsEquivalentTo(&prevInstaller, true) {
				if installer.Checksum == "" &&  prevInstaller.Checksum != "" {
					installer.Checksum = prevInstaller.Checksum
					(*g).Installers[idx] = installer
				}
			}
		}
	}

	for idx, extra := range (*g).Extras {
		if prevExtra, ok := previousExtras[extra.Name]; ok {
			if extra.IsEquivalentTo(&prevExtra, true) {
				if extra.Checksum == "" &&  prevExtra.Checksum != "" {
					extra.Checksum = prevExtra.Checksum
					(*g).Extras[idx] = extra
				}
			}
		}
	}

	return nil
}

func (g *ManifestGame) GetInstallerNamed(name string) (ManifestGameInstaller, error) {
	for idx, _ := range (*g).Installers {
		if (*g).Installers[idx].Name == name {
			return (*g).Installers[idx], nil
		}
	}

	msg := fmt.Sprintf("*ManifestGame.GetInstallerNamed(name=%s) -> No installer by that name", name)
	return ManifestGameInstaller{}, errors.New(msg)
}

func (g *ManifestGame) getExtraNamed(name string) (ManifestGameExtra, error) {
	for idx, _ := range (*g).Extras {
		if (*g).Extras[idx].Name == name {
			return (*g).Extras[idx], nil
		}
	}

	msg := fmt.Sprintf("*ManifestGame.getExtraNamed(name=%s) -> No extra by that name", name)
	return ManifestGameExtra{}, errors.New(msg)
}

func (g *ManifestGame) TrimInstallers(oses []string, languages []string, keepAny bool) {
	filteredInstallers := make([]ManifestGameInstaller, 0)

	if keepAny {
		if len(oses) == 0 && len(languages) == 0 {
			//Save some needless computation
			return
		}

		for _, i := range (*g).Installers {
			hasOneOfOses := len(oses) == 0 || i.HasOneOfOses(oses)
			hasOneOfLanguages := len(languages) == 0 || i.HasOneOfLanguages(languages)
			if hasOneOfOses && hasOneOfLanguages {
				filteredInstallers = append(filteredInstallers, i)
			}
		}
	}
	(*g).Installers = filteredInstallers
}

func (g *ManifestGame) TrimExtras(typeTerms []string, keepAny bool) {
	filteredExtras := make([]ManifestGameExtra, 0)

	if keepAny {
		if len(typeTerms) == 0 {
			return
		}

		for _, e := range (*g).Extras {
			if e.HasOneOfTypeTerms(typeTerms) {
				filteredExtras = append(filteredExtras, e)
			}
		}
	}
	(*g).Extras = filteredExtras
}

func (g *ManifestGame) HasTitleTerms(titleTerms []string) bool {
	if len(titleTerms) == 0 {
		return true
	}

	for idx, _ := range titleTerms {
		if strings.Contains(strings.ToLower((*g).Title), strings.ToLower(titleTerms[idx])) {
			return true
		}
	}
	
	return false
}

func (g *ManifestGame) HasOneOfTags(tags []string) bool {
	for _, t := range tags {
		for _, gt := range (*g).Tags {
			if t == gt {
				return true
			}
		}
	}
	return false
}

func (g *ManifestGame) IsEmpty() bool {
	return len((*g).Installers) == 0 && len((*g).Extras) == 0
}

func (g *ManifestGame) ComputeVerifiedSize() int64 {
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

func (g *ManifestGame) ComputeEstimatedSize() (int64, error) {
	accumulate := int64(0)
	for _, inst := range (*g).Installers {
		size, err := inst.GetEstimatedSizeInBytes()
		if err != nil {
			return int64(0), err
		}
		accumulate += size
	}

	for _, extr := range (*g).Extras {
		size, err := extr.GetEstimatedSizeInBytes()
		if err != nil {
			return 0, err
		}
		accumulate += size
	}

	(*g).EstimatedSize = GetBytesToEstimate(accumulate)
	return accumulate, nil
}

func (g *ManifestGame) FillMissingFileInfo(fileKind string, fileName string, fileSize int64, fileChecksum string) error {
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
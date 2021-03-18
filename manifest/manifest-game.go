package manifest

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

type GameFilenameDuplicates struct {
	Id             int64
	Installers []string
	Extras     []string
}

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

func (g *ManifestGame) CompressIdenticalInstallers() {
	mappedInstallers := map[string]ManifestGameInstaller{}
	
	for _, installer := range (*g).Installers {
		key := fmt.Sprintf("%s|%s|%d", installer.Name, installer.Checksum, installer.VerifiedSize)
		if preExistingInst, ok := mappedInstallers[key]; ok {
			preExistingInst.Languages = ConcatStringSlicesUnique(installer.Languages, preExistingInst.Languages)
			mappedInstallers[key] = preExistingInst
		} else {
			mappedInstallers[key] = installer
		}
	}

	installers := make([]ManifestGameInstaller, len(mappedInstallers))
	idx := 0
	for _, installer := range mappedInstallers {
		installers[idx] = installer
		idx++
	}

	(*g).Installers = installers
}

func (g *ManifestGame) RenameDuplicateFilenames() GameFilenameDuplicates {
	sort.Slice((*g).Installers, func(x, y int) bool {
		if (*g).Installers[x].Name != (*g).Installers[y].Name {
			return (*g).Installers[x].Name < (*g).Installers[y].Name
		} else if (*g).Installers[x].Title != (*g).Installers[y].Title {
			return (*g).Installers[x].Title < (*g).Installers[y].Title
		}
		return (*g).Installers[x].Url < (*g).Installers[y].Url
	})

	duplicateInstallers := make([]string, 0)
	currentName := ""
	occurenceCounter := 0
	for idx, installer := range (*g).Installers {
		if installer.Name == currentName {
			if occurenceCounter == 0 {
				duplicateInstallers = append(duplicateInstallers, installer.Name)
			}
			occurenceCounter++
			suffix := filepath.Ext(installer.Name)
			base := strings.TrimSuffix(installer.Name, suffix)
			installer.Name = fmt.Sprintf("%s[%d]%s", base, occurenceCounter, suffix)
			(*g).Installers[idx] = installer
		} else {
			currentName = installer.Name
			occurenceCounter = 0
		}
	}

	sort.Slice((*g).Extras, func(x, y int) bool {
		if (*g).Extras[x].Name != (*g).Extras[y].Name {
			return (*g).Extras[x].Name < (*g).Extras[y].Name
		} else if (*g).Extras[x].Title != (*g).Extras[y].Title {
			return (*g).Extras[x].Title < (*g).Extras[y].Title
		}
		return (*g).Extras[x].Url < (*g).Extras[y].Url
	})

	duplicateExtras := make([]string, 0)
	currentName = ""
	occurenceCounter = 0
	for idx, extra := range (*g).Extras {
		if extra.Name == currentName {
			if occurenceCounter == 0 {
				duplicateExtras = append(duplicateExtras, extra.Name)
			}
			occurenceCounter++
			suffix := filepath.Ext(extra.Name)
			base := strings.TrimSuffix(extra.Name, suffix)
			extra.Name = fmt.Sprintf("%s[%d]%s", base, occurenceCounter, suffix)
			(*g).Extras[idx] = extra
		} else {
			currentName = extra.Name
			occurenceCounter = 0
		}
	}

	return GameFilenameDuplicates{
		Id: (*g).Id,
		Installers: duplicateInstallers,
		Extras: duplicateExtras,
	}
}

func (g *ManifestGame) TrimIncompleteFiles() {
	installers := make([]ManifestGameInstaller, 0)
	extras := make([]ManifestGameExtra, 0)

	for _, installer := range (*g).Installers {
		if installer.Name != "" && installer.VerifiedSize > 0 {
			installers = append(installers, installer)
		}
	}

	for _, extra := range (*g).Extras {
		if extra.Name != "" && extra.VerifiedSize > 0 {
			extras = append(extras, extra)
		}
	}

	(*g).Installers = installers
	(*g).Extras = extras
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

func (g *ManifestGame) GetExtraNamed(name string) (ManifestGameExtra, error) {
	for idx, _ := range (*g).Extras {
		if (*g).Extras[idx].Name == name {
			return (*g).Extras[idx], nil
		}
	}

	msg := fmt.Sprintf("*ManifestGame.GetExtraNamed(name=%s) -> No extra by that name", name)
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
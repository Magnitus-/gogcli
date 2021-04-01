package manifest

import "strings"

type ManifestGameExtra struct {
	Url           string
	Title         string
	Name          string
	Type          string
	Info          int
	EstimatedSize string
	VerifiedSize  int64
	Checksum      string
}

func (e *ManifestGameExtra) HasOneOfTypeTerms(typeTerms []string) bool {
	for _, t := range typeTerms {
		if strings.Contains((*e).Type, t) {
			return true
		}
	}
	return false
}

func (e *ManifestGameExtra) IsEquivalentTo(o *ManifestGameExtra, emptyChecksumOk bool, ignoreMetadata bool) bool {
	sameName := (*e).Name == (*o).Name
	sameTitle := ((*e).Title == (*o).Title) || ignoreMetadata
	sameUrl := ((*e).Url == (*o).Url) || ignoreMetadata
	sameVerifiedSize := (*o).VerifiedSize != 0 && (*e).VerifiedSize == (*o).VerifiedSize
	checksumIsEmptyAndItsOk := emptyChecksumOk && ((*e).Checksum == "" || (*o).Checksum == "")
	sameChecksum := (*o).Checksum != "" && (*e).Checksum == (*o).Checksum
	return sameName && sameTitle && sameUrl && sameVerifiedSize && (sameChecksum || checksumIsEmptyAndItsOk)
}

func (e *ManifestGameExtra) GetEstimatedSizeInBytes() (int64, error) {
	return GetEstimateToBytes((*e).EstimatedSize)
}
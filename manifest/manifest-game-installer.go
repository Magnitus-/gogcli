package manifest

type ManifestGameInstaller struct {
	Languages     []string
	Os            string
	Url           string
	Title         string
	Name          string
	Version       string
	Date          string
	EstimatedSize string
	VerifiedSize  int64
	Checksum      string
}

func (i *ManifestGameInstaller) HasOneOfOses(oses []string) bool {
	for _, os := range oses {
		if os == i.Os {
			return true
		}
	}
	return false
}

func (i *ManifestGameInstaller) HasOneOfLanguages(languages []string) bool {
	for _, l := range languages {
		for _, l2 := range i.Languages {
			if l == l2 {
				return true
			}
		}
	}
	return false
}

func (i *ManifestGameInstaller) IsEquivalentTo(o *ManifestGameInstaller, checksumValidation string, ignoreMetadata bool) bool {
	sameName := (*i).Name == (*o).Name
	sameTitle := ((*i).Title == (*o).Title) || ignoreMetadata
	sameUrl := ((*i).Url == (*o).Url) || ignoreMetadata
	sameVerifiedSize := (*o).VerifiedSize != 0 && (*i).VerifiedSize == (*o).VerifiedSize
	checksumIsEmptyAndItsOk := checksumValidation == ChecksumValidationIfPresent && ((*i).Checksum == "" || (*o).Checksum == "")
	sameChecksum := (*o).Checksum != "" && (*i).Checksum == (*o).Checksum
	return sameName && sameTitle && sameUrl && sameVerifiedSize && (checksumValidation == ChecksumNoValidation || sameChecksum || checksumIsEmptyAndItsOk)
}

func (i *ManifestGameInstaller) GetEstimatedSizeInBytes() (int64, error) {
	return GetEstimateToBytes((*i).EstimatedSize)
}

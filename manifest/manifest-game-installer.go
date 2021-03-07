package manifest

type ManifestGameInstaller struct {
	Language      string
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
		if l == i.Language {
			return true
		}
	}
	return false
}

func (i *ManifestGameInstaller) IsEquivalentTo(o *ManifestGameInstaller, emptyChecksumOk bool) bool {
	sameName := (*i).Name == (*o).Name
	sameTitle := (*i).Title == (*o).Title
	sameUrl := (*i).Url == (*o).Url
	sameVerifiedSize := (*o).VerifiedSize != 0 && (*i).VerifiedSize == (*o).VerifiedSize
	checksumIsEmptyAndItsOk := emptyChecksumOk && ((*i).Checksum == "" || (*o).Checksum == "")
	sameChecksum := (*o).Checksum != "" && (*i).Checksum == (*o).Checksum
	return sameName && sameTitle && sameUrl && sameVerifiedSize && (sameChecksum || checksumIsEmptyAndItsOk)
}

func (i *ManifestGameInstaller) GetEstimatedSizeInBytes() (int64, error) {
	return GetEstimateToBytes((*i).EstimatedSize)
}

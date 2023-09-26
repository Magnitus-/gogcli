package metadata

type MetadataFilter struct {
	Titles          []string
}

func NewMetadataFilter(titles []string) MetadataFilter {
	return MetadataFilter{Titles: titles}
}
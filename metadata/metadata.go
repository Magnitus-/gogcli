package metadata

type GameMetadataDescription struct {
	Summary    string
	Full       string
	Highlights string
}

type Image struct {
	Name     string
	Size     int64
	Checksum string
	Url      string
	Tag      string
}

type Video struct {
	ThumbnailUrl string
	VideoUrl     string
	Provider     string
}

type MetadataGame struct {
	Id                    int64
	Title                 string
	Tags                  []string
	Description           GameMetadataDescription
	ProductCards          []Image
	OtherProductImages    []Image
    Screenshots           []Image
	Videos                []Video
	Slug                  string
	ReleaseDate           string
	Rating                int
	Category              string
	Dlcs                  int
	Features              []string
	Changelog             string
}

type Metadata struct {
	Games []MetadataGame
	Size  int64
}

func NewEmptyMetadata() *Metadata {
	m := Metadata{
		Games: []MetadataGame{},
		Size: 0,
	}
	return &m
}
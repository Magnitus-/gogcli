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
}

type Video struct {
	ThumbnailUrl string
	VideoUrl     string
	Provider     string
}

type GameMetadata struct {
	Id                    int64
	Title                 string
	Description           GameMetadataDescription
	ProductCards          []Image
	OtherProductImages    []Image
    Screenshots           []Image
	Videos                []Video
	Slug                  string
	ReleaseDate           string
	Rating                int64
	Category              string
	Dlcs                  int
	Features              []string
	Changelog             string
}
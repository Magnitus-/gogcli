package metadata

type GameMetadataDescription struct {
	Summary    string
	Full       string
	Highlights string
}

type GameMetadataImage struct {
	Name     string
	Size     int64
	Checksum string
	Url      string
	Tag      string
}

type GameMetadataVideo struct {
	ThumbnailUrl string
	VideoUrl     string
	Provider     string
}

type GameMetadataProductImages struct {
	Background          GameMetadataImage
	Logo                GameMetadataImage
	Icon                GameMetadataImage
}

type GameMetadataScreenShot struct {
	List GameMetadataImage
	Main GameMetadataImage
}

type MetadataGame struct {
	Id             int64
	Title          string
	Tags           []string
	ListingImage   GameMetadataImage
	Description    GameMetadataDescription
	ProductImages  GameMetadataProductImages
	Screenshots    []GameMetadataScreenShot
	Videos         []GameMetadataVideo
	Slug           string
	ReleaseDate    string
	Rating         int
	Category       string
	Dlcs           int
	Changelog      string
	HasProductInfo bool
}

func (g *MetadataGame) GetImagesPointers() []*GameMetadataImage {
	result := []*GameMetadataImage{
		&(*g).ListingImage,
		&(*g).ProductImages.Background,
		&(*g).ProductImages.Logo,
		&(*g).ProductImages.Icon,
	}

	for idx, _ := range (*g).Screenshots {
		result = append(result, &(*g).Screenshots[idx].List)
		result = append(result, &(*g).Screenshots[idx].Main)
	}
	return result
}

type Metadata struct {
	Games []MetadataGame
	Size  int64
}

func NewEmptyMetadata() *Metadata {
	m := Metadata{
		Games: []MetadataGame{},
		Size:  0,
	}
	return &m
}

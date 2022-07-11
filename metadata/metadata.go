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
	Logo2x              GameMetadataImage
	Icon                GameMetadataImage
	SidebarIcon         GameMetadataImage
	SidebarIcon2x       GameMetadataImage
	MenuNotificationAv  GameMetadataImage
	MenuNotificationAv2 GameMetadataImage
}

type GameMetadataScreenShot []GameMetadataImage

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
		&(*g).ProductImages.Logo2x,
		&(*g).ProductImages.Icon,
		&(*g).ProductImages.SidebarIcon,
		&(*g).ProductImages.SidebarIcon2x,
		&(*g).ProductImages.MenuNotificationAv,
		&(*g).ProductImages.MenuNotificationAv2,
	}

	for idx, _ := range (*g).Screenshots {
		for idxInner := range (*g).Screenshots[idx] {
			result = append(result, &(*g).Screenshots[idx][idxInner])
		}
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

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
	Id            int64
	Title         string
	Tags          []string
	ListingImage  GameMetadataImage
	Description   GameMetadataDescription
	ProductCards  []GameMetadataImage
	ProductImages GameMetadataProductImages
	Screenshots   []GameMetadataScreenShot
	Videos        []GameMetadataVideo
	Slug          string
	ReleaseDate   string
	Rating        int
	Category      string
	Dlcs          int
	Features      []string
	Changelog     string
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

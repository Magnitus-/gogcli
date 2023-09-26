package metadata

import (
	"strings"
)

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
	result := []*GameMetadataImage{}

	if (*g).ListingImage.Url != "" {
		result = append(result, &(*g).ListingImage)
	}

	if (*g).ProductImages.Background.Url != "" {
		result = append(result, &(*g).ProductImages.Background)
	}

	if (*g).ProductImages.Logo.Url != "" {
		result = append(result, &(*g).ProductImages.Logo)
	}

	if (*g).ProductImages.Icon.Url != "" {
		result = append(result, &(*g).ProductImages.Icon)
	}

	for idx, _ := range (*g).Screenshots {
		result = append(result, &(*g).Screenshots[idx].List)
		result = append(result, &(*g).Screenshots[idx].Main)
	}
	return result
}

func (g *MetadataGame) HasTitleTerms(titleTerms []string) bool {
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

type Metadata struct {
	Games      []MetadataGame
	Filter     MetadataFilter
	SkipImages []string
	Size       int64
}

func NewEmptyMetadata(filter MetadataFilter, skipImages []string) *Metadata {
	m := Metadata{
		Games: []MetadataGame{},
		Filter: filter,
		SkipImages: skipImages,
		Size:  0,
	}
	return &m
}

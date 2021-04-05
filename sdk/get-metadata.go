package sdk

import "gogcli/metadata"

func addOwnedGamesPagesToMetadata(m *metadata.Metadata, pages []OwnedGamesPage) {
	for _, page := range pages {
		for _, product := range page.Products {
			g := metadata.MetadataGame{
				Id:    product.Id,
				Title: product.Title,
				Slug: product.Slug,
				Category: product.Category,
				Rating: product.Rating,
				Dlcs: product.DlcCount,
			}
			tags := []string{}
			for _, tag := range product.Tags {
				for _, tagDetails := range page.Tags {
					if tagDetails.Id == tag {
						tags = append(tags, tagDetails.Name)
					}
				}
			}
			g.Tags = tags
			(*m).Games = append(
				(*m).Games,
				g,
			)
		}
	}
}

/*
Left:
type MetadataGame struct {
	ProductCards          []Image
	OtherProductImages    []Image
    Screenshots           []Image
	Features              []string
}
*/

func updateMetadataWithProducts(m *metadata.Metadata, products []Product) {
	productsMap := map[int64]Product{}
	for _, p := range products {
		productsMap[p.Id] = p
	}

	for idx, _ := range (*m).Games {
		game := (*m).Games[idx]
		product := productsMap[game.Id]
		game.Description = metadata.GameMetadataDescription{
			Summary: product.Description.Lead,
			Full: product.Description.Full,
			Highlights: product.Description.Whats_cool_about_it,
		}
		game.ReleaseDate = product.Release_date
		game.Changelog = product.Changelog
		
		videos := []metadata.Video{}
		for _, vid := range product.Videos {
			videos = append(videos, metadata.Video{
				ThumbnailUrl: vid.Thumbnail_url,
				VideoUrl: vid.Video_url,
				Provider: vid.Provider,
			})
		}
		game.Videos = videos
		
		otherProductImages := []metadata.Image{}
		game.OtherProductImages = otherProductImages

		screenshots := []metadata.Image{}
		game.Screenshots = screenshots

		(*m).Games[idx] = game
	}
}

func (s *Sdk) GetMetadata(concurrency int, pause int, tolerateDangles bool) (metadata.Metadata, []error, []error) {
	m := metadata.NewEmptyMetadata()

	pages, errs := s.GetAllOwnedGamesPages("", concurrency, pause)
	if len(errs) > 0 {
		return *m, errs, []error{}
	}

	addOwnedGamesPagesToMetadata(m, pages)

	gameIds := make([]int64, len(m.Games))
	for i := 0; i < len(m.Games); i++ {
		gameIds[i] = m.Games[i].Id
	}

	products, productsErrs, productWarnings := s.GetManyProducts(gameIds, concurrency, pause)
	if len(productsErrs) > 0 {
		return *m, productsErrs, productWarnings
	}

	updateMetadataWithProducts(m, products)

	return *m, []error{}, productWarnings
}
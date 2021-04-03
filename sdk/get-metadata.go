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
			(*m).Games = append(
				(*m).Games,
				g,
			)
		}
	}
}

func (s *Sdk) GetMetadata(concurrency int, pause int, tolerateDangles bool) (metadata.Metadata, []error, []error) {
	m := metadata.NewEmptyMetadata()

	pages, errs := s.GetAllOwnedGamesPages("", concurrency, pause)
	if len(errs) > 0 {
		return *m, errs, []error{}
	}

	addOwnedGamesPagesToMetadata(m, pages)

	return (*m), []error{}, []error{}
}
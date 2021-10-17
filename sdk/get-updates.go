package sdk

import "gogcli/manifest"

func (s *Sdk) GetUpdates(concurrency int, pause int) (manifest.Updates, []error) {
	updates := manifest.NewEmptyUpdates()

	pages, err := s.GetAllOwnedGamesPages("", concurrency, pause)
	if err != nil {
		return *updates, err
	}

	for _, page := range pages {
		for _, product := range page.Products {
			if product.IsNew {
				updates.AddNewGame(product.Id)
			} else if product.Updates > 0 {
				updates.AddUpdatedGame(product.Id)
			}
		}
	}

	return *updates, nil
}

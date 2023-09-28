package sdk

import "gogcli/gameupdates"

func (s *Sdk) GetUpdates(concurrency int, pause int) (gameupdates.Updates, []error) {
	updates := gameupdates.NewEmptyUpdates()

	pages, err := s.GetAllOwnedGamesPagesSync("", concurrency, pause)
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

package sdk

import (
	"gogcli/manifest"
)

func (o *OwnedGamesPage) GetMissingGames(m *manifest.Manifest) []product {
	missingGames := []product{}

	for _, product := range o.Products {
		found := false
		for _, game := range m.Games {
			if game.Id == product.Id {
				found = true
				break
			}
		}

		if !found {
			missingGames = append(missingGames, product)
		}
	}

	return missingGames
}

func GetMissingGames(ownedGames []OwnedGamesPage, m *manifest.Manifest) []product {
	missingGames := []product{}

	for idx, _ := range ownedGames {
		missingGames = append(missingGames, ownedGames[idx].GetMissingGames(m)...)
	}

	return missingGames
}
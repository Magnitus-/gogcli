package sdk

import (
	"gogcli/manifest"
)

func addOwnedGamesPagesToManifest(m *manifest.Manifest, pages []OwnedGamesPage) {
	for _, page := range pages {
		for _, product := range page.Products {
			g := manifest.ManifestGame{
				Id:    product.Id,
				Title: product.Title,
				Slug:  product.Slug,
			}
			(*m).Games = append(
				(*m).Games,
				g,
			)
		}
	}
}

func updateManifestWithGameDetails(m *manifest.Manifest, gameDetails []GameDetailsWithId) {
	for _, gd := range gameDetails {
		for gidx, _ := range (*m).Games {
			if gd.id == (*m).Games[gidx].Id {
				(*m).Games[gidx].CdKey = gd.game.CdKey

				(*m).Games[gidx].Tags = make([]string, len(gd.game.Tags))
				for i, _ := range gd.game.Tags {
					(*m).Games[gidx].Tags[i] = gd.game.Tags[i].Name
				}

				for _, i := range gd.game.Downloads {
					(*m).Games[gidx].Installers = append(
						(*m).Games[gidx].Installers,
						manifest.ManifestGameInstaller{
							Languages:     []string{languageToAscii(i.Language)},
							Os:            i.Os,
							Url:           i.ManualUrl,
							Title:         i.Name,
							Version:       i.Version,
							Date:          i.Date,
							EstimatedSize: i.Size,
						},
					)
				}

				for _, e := range gd.game.Extras {
					(*m).Games[gidx].Extras = append(
						(*m).Games[gidx].Extras,
						manifest.ManifestGameExtra{
							Url:           e.ManualUrl,
							Title:         e.Name,
							Type:          e.Type,
							Info:          e.Info,
							EstimatedSize: e.Size,
						},
					)
				}

				for _, d := range gd.game.Dlcs {
					for _, i := range d.Downloads {
						(*m).Games[gidx].Installers = append(
							(*m).Games[gidx].Installers,
							manifest.ManifestGameInstaller{
								Languages:     []string{languageToAscii(i.Language)},
								Os:            i.Os,
								Url:           i.ManualUrl,
								Title:         i.Name,
								Version:       i.Version,
								Date:          i.Date,
								EstimatedSize: i.Size,
							},
						)
					}

					for _, e := range d.Extras {
						(*m).Games[gidx].Extras = append(
							(*m).Games[gidx].Extras,
							manifest.ManifestGameExtra{
								Url:           e.ManualUrl,
								Title:         e.Name,
								Type:          e.Type,
								Info:          e.Info,
								EstimatedSize: e.Size,
							},
						)
					}
				}
			}
		}
	}
}

func (s *Sdk) GetManifest(f manifest.ManifestFilter, concurrency int, pause int, tolerateDangles bool, tolerateBadMetadata bool) (manifest.Manifest, []error, []error) {
	m := manifest.NewEmptyManifest(f)

	pages, errs := s.GetAllOwnedGamesPagesSync("", concurrency, pause)
	if len(errs) > 0 {
		return *m, errs, []error{}
	}

	addOwnedGamesPagesToManifest(m, pages)
	m.TrimGames()

	gameIds := make([]int64, len(m.Games))
	for i := 0; i < len(m.Games); i++ {
		gameIds[i] = m.Games[i].Id
	}

	details, detailsErrs := s.GetManyGameDetails(gameIds, concurrency, pause)
	if len(detailsErrs) > 0 {
		return *m, detailsErrs, []error{}
	}

	updateManifestWithGameDetails(m, details)
	m.Trim()

	errs, warnings := s.fillManifestFiles(m, concurrency, pause, tolerateDangles, tolerateBadMetadata)
	return *m, errs, warnings
}

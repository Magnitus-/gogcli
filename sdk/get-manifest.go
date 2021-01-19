package sdk

import (
	"gogcli/manifest"
	"strings"
)

var LANGUAGE_MAP map[string]string

func init() {
	LANGUAGE_MAP = make(map[string]string)
	LANGUAGE_MAP["english"] = "English"
	LANGUAGE_MAP["french"] = "fran\\u00e7ais"
	LANGUAGE_MAP["dutch"] = "nederlands"
	LANGUAGE_MAP["spanish"] = "espa\\u00f1ol"
	LANGUAGE_MAP["portuguese_brazilian"] = "Portugu\\u00eas do Brasil"
	LANGUAGE_MAP["russian"] = "\\u0440\\u0443\\u0441\\u0441\\u043a\\u0438\\u0439"
	LANGUAGE_MAP["korean"] = "\\ud55c\\uad6d\\uc5b4"
	LANGUAGE_MAP["chinese_simplified"] = "\\u4e2d\\u6587(\\u7b80\\u4f53)"
	LANGUAGE_MAP["japanese"] = "\\u65e5\u672c\\u8a9e"
	LANGUAGE_MAP["polish"] = "polski"
	LANGUAGE_MAP["italian"] = "italiano"
	LANGUAGE_MAP["german"] = "Deutsch"
	LANGUAGE_MAP["czech"] = "\\u010desk\\u00fd"
	LANGUAGE_MAP["hungarian"] = "magyar"
	LANGUAGE_MAP["portuguese"] = "portugu\\u00eas"
	LANGUAGE_MAP["danish"] = "Dansk"
	LANGUAGE_MAP["finnish"] = "suomi"
	LANGUAGE_MAP["swedish"] = "svenska"
	LANGUAGE_MAP["turkish"] = "T\\u00fcrk\\u00e7e"
	LANGUAGE_MAP["arabic"] = "\\u0627\\u0644\\u0639\\u0631\\u0628\\u064a\\u0629"
	LANGUAGE_MAP["romanian"] = "rom\\u00e2n\\u0103"
}

func languageToAscii(unicodeRepresentation string) string {
	for k, v := range LANGUAGE_MAP {
		if strings.EqualFold(v, unicodeRepresentation) {
			return k
		}
	}
	return "unknown"
}

func addOwnedGamesPagesToManifest(m *manifest.Manifest, pages []OwnedGamesPage) {
	for _, page := range pages {
		for _, product := range page.Products {
			g := manifest.ManifestGame{
				Id:    product.Id,
				Title: product.Title,
			}
			(*m).Games = append(
				(*m).Games,
				g,
			)
		}
	}
}

func addGameDetailsToManifest(m *manifest.Manifest, gameDetails []GameDetailsWithId) {
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
							Language: languageToAscii(i.Language),
							Os:       i.Os,
							Url:      i.ManualUrl,
							Name:     i.Name,
							Version:  i.Version,
							Date:     i.Date,
							Size:     i.Size,
						},
					)
				}

				for _, e := range gd.game.Extras {
					(*m).Games[gidx].Extras = append(
						(*m).Games[gidx].Extras,
						manifest.ManifestGameExtra{
							Url:  e.ManualUrl,
							Name: e.Name,
							Type: e.Type,
							Info: e.Info,
							Size: e.Size,
						},
					)
				}

				for _, d := range gd.game.Dlcs {
					for _, i := range d.Downloads {
						(*m).Games[gidx].Installers = append(
							(*m).Games[gidx].Installers,
							manifest.ManifestGameInstaller{
								Language: languageToAscii(i.Language),
								Os:       i.Os,
								Url:      i.ManualUrl,
								Name:     i.Name,
								Version:  i.Version,
								Date:     i.Date,
								Size:     i.Size,
							},
						)
					}

					for _, e := range d.Extras {
						(*m).Games[gidx].Extras = append(
							(*m).Games[gidx].Extras,
							manifest.ManifestGameExtra{
								Url:  e.ManualUrl,
								Name: e.Name,
								Type: e.Type,
								Info: e.Info,
								Size: e.Size,
							},
						)
					}
				}
			}
		}
	}
}

func (s *Sdk) GetManifest(search string, concurrency int, pause int, debug bool) (manifest.Manifest, []error) {
	var m manifest.Manifest

	pages, errs := s.GetAllOwnedGamesPages(search, concurrency, pause, debug)
	if len(errs) > 0 {
		return m, errs
	}

	addOwnedGamesPagesToManifest(&m, pages)

	gameIds := make([]int, len(m.Games))
	for i := 0; i < len(m.Games); i++ {
		gameIds[i] = m.Games[i].Id
	}

	details, detailsErrs := s.GetManyGameDetails(gameIds, concurrency, pause, debug)
	if len(detailsErrs) > 0 {
		return m, detailsErrs
	}

	addGameDetailsToManifest(&m, details)

	return m, nil
}

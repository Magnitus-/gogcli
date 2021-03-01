package sdk

import (
	"gogcli/manifest"
	"strings"
)

var LANGUAGE_MAP map[string]string
func getLanguageMap() map[string]string {
	langMap := make(map[string]string)
	langMap["english"] = "English"
	langMap["french"] = "fran\\u00e7ais"
	langMap["dutch"] = "nederlands"
	langMap["spanish"] = "espa\\u00f1ol"
	langMap["portuguese_brazilian"] = "Portugu\\u00eas do Brasil"
	langMap["russian"] = "\\u0440\\u0443\\u0441\\u0441\\u043a\\u0438\\u0439"
	langMap["korean"] = "\\ud55c\\uad6d\\uc5b4"
	langMap["chinese_simplified"] = "\\u4e2d\\u6587(\\u7b80\\u4f53)"
	langMap["japanese"] = "\\u65e5\u672c\\u8a9e"
	langMap["polish"] = "polski"
	langMap["italian"] = "italiano"
	langMap["german"] = "Deutsch"
	langMap["czech"] = "\\u010desk\\u00fd"
	langMap["hungarian"] = "magyar"
	langMap["portuguese"] = "portugu\\u00eas"
	langMap["danish"] = "Dansk"
	langMap["finnish"] = "suomi"
	langMap["swedish"] = "svenska"
	langMap["turkish"] = "T\\u00fcrk\\u00e7e"
	langMap["arabic"] = "\\u0627\\u0644\\u0639\\u0631\\u0628\\u064a\\u0629"
	langMap["romanian"] = "rom\\u00e2n\\u0103"
	return langMap
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
							Language:      languageToAscii(i.Language),
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
								Language:      languageToAscii(i.Language),
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

func (s *Sdk) GetManifest(f manifest.ManifestFilter, concurrency int, pause int, tolerateDangles bool) (manifest.Manifest, []error, []error) {
	m := manifest.NewEmptyManifest(f)

	pages, errs := s.GetAllOwnedGamesPages("", concurrency, pause)
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

	addGameDetailsToManifest(m, details)
	m.Trim()

	installersMap := m.GetUrlMappedInstallers()
	installerUrls := make([]string, len(installersMap))
	idx := 0
	for k, _ := range installersMap {
		installerUrls[idx] = k
		idx++;
	}

	downloadInfos, fileInfoErrs, danglingInstallerErrs := s.GetManyDownloadFileInfo(installerUrls, concurrency, pause, tolerateDangles)
	if len(fileInfoErrs) > 0 {
		return *m, fileInfoErrs, danglingInstallerErrs
	}
	for _, downloadFileInfo := range downloadInfos {
		(*installersMap[downloadFileInfo.url]).Name = downloadFileInfo.name
		(*installersMap[downloadFileInfo.url]).Checksum = downloadFileInfo.checksum
		(*installersMap[downloadFileInfo.url]).VerifiedSize = downloadFileInfo.size
	} 

	extrasMap := m.GetUrlMappedExtras()
	extraUrls := make([]string, len(extrasMap))
	idx = 0
	for k, _ := range extrasMap {
		extraUrls[idx] = k
		idx++;
	}

	var danglingExtraErrs []error
	downloadInfos, fileInfoErrs, danglingExtraErrs = s.GetManyDownloadFileInfo(extraUrls, concurrency, pause, tolerateDangles)
	if len(fileInfoErrs) > 0 {
		return *m, fileInfoErrs, append(danglingInstallerErrs, danglingExtraErrs...)
	}
	for _, downloadFileInfo := range downloadInfos {
		(*extrasMap[downloadFileInfo.url]).Name = downloadFileInfo.name
		(*extrasMap[downloadFileInfo.url]).Checksum = downloadFileInfo.checksum
		(*extrasMap[downloadFileInfo.url]).VerifiedSize = downloadFileInfo.size
	} 

	return *m, []error{}, append(danglingInstallerErrs, danglingExtraErrs...)
}

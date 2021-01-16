package sdk

import (
	"gogcli/manifest"
)

var LANGUAGE_MAP map[string]string

func init() {
	LANGUAGE_MAP = make(map[string]string)
	LANGUAGE_MAP["english"] = "English"
	LANGUAGE_MAP["french"] = "fran\u00e7ais"
	LANGUAGE_MAP["dutch"] = "nederlands"
	LANGUAGE_MAP["spanish"] = "espa\u00f1ol"
	LANGUAGE_MAP["portuguese_brazilian"] = "Portugu\u00eas do Brasil"
	LANGUAGE_MAP["russian"] = "\u0440\u0443\u0441\u0441\u043a\u0438\u0439"
	LANGUAGE_MAP["korean"] = "\ud55c\uad6d\uc5b4"
	LANGUAGE_MAP["chinese_simplified"] = "\u4e2d\u6587(\u7b80\u4f53)"
	LANGUAGE_MAP["japanese"] = "\u65e5\u672c\u8a9e"
	LANGUAGE_MAP["polish"] = "polski"
	LANGUAGE_MAP["italian"] = "italiano"
	LANGUAGE_MAP["german"] = "Deutsch"
	LANGUAGE_MAP["czech"] = "\u010desk\u00fd"
	LANGUAGE_MAP["hungarian"] = "magyar"
	LANGUAGE_MAP["portuguese"] = "portugu\u00eas"
	LANGUAGE_MAP["danish"] = "Dansk"
	LANGUAGE_MAP["finnish"] = "suomi"
	LANGUAGE_MAP["swedish"] = "svenska"
	LANGUAGE_MAP["turkish"] = "T\u00fcrk\u00e7e"
	LANGUAGE_MAP["arabic"] = "\u0627\u0644\u0639\u0631\u0628\u064a\u0629"
	LANGUAGE_MAP["romanian"] = "rom\u00e2n\u0103"
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

func (s *Sdk) GetManifest(search string, concurrency int, pause int, debug bool) (manifest.Manifest, []error) {
	var m manifest.Manifest

	pages, errs := s.GetAllOwnedGamesPages(search, concurrency, pause, debug)
	if len(errs) > 0 {
		return m, errs
	}

	addOwnedGamesPagesToManifest(&m, pages)
	return m, nil
}

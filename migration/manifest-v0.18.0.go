package migration

import (
	"gogcli/manifest"
	"gogcli/sdk"
)

type ManifestGameV0_18 struct {
	Id            int64
	Title         string
	CdKey         string
	Tags          []string
	Installers    []manifest.ManifestGameInstaller
	Extras        []manifest.ManifestGameExtra
	EstimatedSize string
	VerifiedSize  int64
}

type ManifestV0_18 struct {
	Games         []ManifestGameV0_18
	EstimatedSize string
	VerifiedSize  int64
	Filter        manifest.ManifestFilter
}

func addSlugsToManifest(m *manifest.Manifest, pages []sdk.OwnedGamesPage) {
	mappedSlugs := map[int64]string{}
	for _, page := range pages {
		for _, product := range page.Products {
			mappedSlugs[product.Id] = product.Slug
		}
	}

	for idx, game := range m.Games {
	    game.Slug = mappedSlugs[game.Id]
		m.Games[idx] = game
	}
}

func (m *ManifestV0_18) Migrate(s *sdk.Sdk) (*manifest.Manifest, []error) {
	migrated := manifest.Manifest{
		Games:         make([]manifest.ManifestGame, len((*m).Games)),
		EstimatedSize: (*m).EstimatedSize,
		VerifiedSize:  (*m).VerifiedSize,
		Filter:        (*m).Filter,
	}

	for idx, game := range (*m).Games {
		gameCopy := manifest.ManifestGame{
			Id:            game.Id,
			Title:         game.Title,
			CdKey:         game.CdKey,
			Tags:          game.Tags,
			Installers:    game.Installers,
			Extras:        game.Extras,
			EstimatedSize: game.EstimatedSize,
			VerifiedSize:  game.VerifiedSize,
		}
		migrated.Games[idx] = gameCopy
	}

	pages, errs := s.GetAllOwnedGamesPages("", 10, 200)
	if len(errs) > 0 {
		return &migrated, errs
	}

	addSlugsToManifest(&migrated, pages)

	migrated.Finalize()
	return &migrated, []error{}
}
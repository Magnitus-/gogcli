package sdk

import "gogcli/manifest"

func addGameDetailsToManifest(m *manifest.Manifest, gameDetails []GameDetailsWithId) {
	for _, gd := range gameDetails {
		game := manifest.ManifestGame{
			Id: gd.id,
			Title: gd.game.Title,
			CdKey: gd.game.CdKey,
			Tags: make([]string, len(gd.game.Tags)),
			Installers: make([]manifest.ManifestGameInstaller, 0),
			Extras: make([]manifest.ManifestGameExtra, 0),
		}
		for i, _ := range gd.game.Tags {
			game.Tags[i] = gd.game.Tags[i].Name
		}
		for _, i := range gd.game.Downloads {
			game.Installers = append(
				game.Installers,
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
			game.Extras = append(
				game.Extras,
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
				game.Installers = append(
					game.Installers,
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
				game.Extras = append(
					game.Extras,
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
		(*m).Games = append((*m).Games, game)
	}
}

func (s *Sdk) GetManifestFromIds(f manifest.ManifestFilter, gameIds []int64, concurrency int, pause int, tolerateDangles bool) (*manifest.Manifest, []error, []error) {
	m := manifest.NewEmptyManifest(f)
	
	details, detailsErrs := s.GetManyGameDetails(gameIds, concurrency, pause)
	if len(detailsErrs) > 0 {
		return m, detailsErrs, []error{}
	}

	addGameDetailsToManifest(m, details)
	m.Trim()

	errs, warnings := s.fillManifestFiles(m, concurrency, pause, tolerateDangles)
	return m, errs, warnings
}
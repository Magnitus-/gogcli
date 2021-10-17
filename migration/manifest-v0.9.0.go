package migration

import "gogcli/manifest"

type ManifestGameInstallerV0_9 struct {
	Language      string
	Os            string
	Url           string
	Title         string
	Name          string
	Version       string
	Date          string
	EstimatedSize string
	VerifiedSize  int64
	Checksum      string
}

type ManifestGameV0_9 struct {
	Id            int64
	Title         string
	CdKey         string
	Tags          []string
	Installers    []ManifestGameInstallerV0_9
	Extras        []manifest.ManifestGameExtra
	EstimatedSize string
	VerifiedSize  int64
}

type ManifestV0_9 struct {
	Games         []ManifestGameV0_9
	EstimatedSize string
	VerifiedSize  int64
	Filter        manifest.ManifestFilter
}

func (m *ManifestV0_9) Migrate() *manifest.Manifest {
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
			Installers:    make([]manifest.ManifestGameInstaller, len(game.Installers)),
			Extras:        game.Extras,
			EstimatedSize: game.EstimatedSize,
			VerifiedSize:  game.VerifiedSize,
		}
		for idx2, inst := range game.Installers {
			gameCopy.Installers[idx2] = manifest.ManifestGameInstaller{
				Languages:     []string{inst.Language},
				Os:            inst.Os,
				Url:           inst.Url,
				Title:         inst.Title,
				Name:          inst.Name,
				Version:       inst.Version,
				Date:          inst.Date,
				EstimatedSize: inst.EstimatedSize,
				VerifiedSize:  inst.VerifiedSize,
				Checksum:      inst.Checksum,
			}
		}
		migrated.Games[idx] = gameCopy
	}

	migrated.Finalize()
	return &migrated
}

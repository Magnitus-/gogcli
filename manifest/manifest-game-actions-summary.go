package manifest

type GamesActionsSummary struct {
	GameAdditions      int
	GameDeletions      int
	GameUpdates        int
	InstallerUpserts   int
	InstallerDeletions int
	ExtraUpserts       int
	ExtraDeletions     int
}

type GameFilesActionsSummary struct {
	InstallerUpserts   int
	InstallerDeletions int
	ExtraUpserts       int
	ExtraDeletions     int
}

func (g *GameAction) GetFilesActionSummary() GameFilesActionsSummary {
	summary := GameFilesActionsSummary{InstallerUpserts: 0, InstallerDeletions: 0, ExtraUpserts: 0, ExtraDeletions: 0}
	for _, fileAction := range (*g).InstallerActions {
		if fileAction.Action == "add" {
			summary.InstallerUpserts += 1
		} else {
			summary.InstallerDeletions += 1
		}
	}

	for _, fileAction := range (*g).ExtraActions {
		if fileAction.Action == "add" {
			summary.ExtraUpserts += 1
		} else {
			summary.ExtraDeletions += 1
		}
	}

	return summary
}

func (g *GameActions) GetSummary() GamesActionsSummary {
	summary := GamesActionsSummary{
		GameAdditions:      0,
		GameDeletions:      0,
		GameUpdates:        0,
		InstallerUpserts:   0,
		InstallerDeletions: 0,
		ExtraUpserts:       0,
		ExtraDeletions:     0,
	}

	for _, gameAction := range *g {
		if gameAction.Action == "add" {
			summary.GameAdditions += 1
		} else if gameAction.Action == "remove" {
			summary.GameDeletions += 1
		} else {
			summary.GameUpdates += 1
		}
		gameSummary := gameAction.GetFilesActionSummary()
		summary.InstallerUpserts += gameSummary.InstallerUpserts
		summary.InstallerDeletions += gameSummary.InstallerDeletions
		summary.ExtraUpserts += gameSummary.ExtraUpserts
		summary.ExtraDeletions += gameSummary.ExtraDeletions
	}

	return summary
}

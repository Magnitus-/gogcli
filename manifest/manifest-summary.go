package manifest

type GameSummary struct {
	Id           int64
	Title        string
	Size         int64
	SizeAsString string
	Installers   int
	Extras       int
}

type ManifestSummary struct {
	Games               int
	Files               int
	Installers          int
	Extras              int
	Size                int64
	SizeAsString        string
	SizeAverage         int64
	SizeAverageAsString string
	LargestGame         GameSummary
	SmallestGame        GameSummary
}

func (m *Manifest) GetSummary() ManifestSummary {
	filesCount := 0
	installersCount := 0
	extrasCount := 0
	var largestGame GameSummary
	var smallestGame GameSummary

	for _, game := range (*m).Games {
		filesCount += (len(game.Installers) + len(game.Extras))
		installersCount += len(game.Installers)
		extrasCount += len(game.Extras)

		if largestGame.Id == 0 || (game.VerifiedSize > largestGame.Size) {
			largestGame.Id = game.Id
			largestGame.Title = game.Title
			largestGame.Size = game.VerifiedSize
			largestGame.Installers = len(game.Installers)
			largestGame.Extras = len(game.Extras)
		}

		if smallestGame.Id == 0 || (game.VerifiedSize < smallestGame.Size) {
			smallestGame.Id = game.Id
			smallestGame.Title = game.Title
			smallestGame.Size = game.VerifiedSize
			smallestGame.Installers = len(game.Installers)
			smallestGame.Extras = len(game.Extras)
		}
	}

	largestGame.SizeAsString = GetBytesToEstimate(largestGame.Size)
	smallestGame.SizeAsString = GetBytesToEstimate(smallestGame.Size)

	return ManifestSummary{
		Games:               len((*m).Games),
		Files:               filesCount,
		Installers:          installersCount,
		Extras:              extrasCount,
		Size:                (*m).VerifiedSize,
		SizeAsString:        GetBytesToEstimate((*m).VerifiedSize),
		SizeAverage:         (*m).VerifiedSize / int64(len((*m).Games)),
		SizeAverageAsString: GetBytesToEstimate((*m).VerifiedSize / int64(len((*m).Games))),
		LargestGame:         largestGame,
		SmallestGame:        smallestGame,
	}
}

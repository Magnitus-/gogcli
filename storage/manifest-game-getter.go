package storage

import (
	//"sync"

	"gogcli/manifest"
)

/*
type ManifestGameGetterGame struct {
	Game ManifestGame
	Warnings []error
	Errors   []error
}

type ManifestGameGetterGameIds struct {
	Ids   []int64
	Error error
}
*/

func getStorageGameIds(s Storage) <-chan manifest.ManifestGameGetterGameIds {
	gameIdsCh := make(chan manifest.ManifestGameGetterGameIds)
	
	go func() {
		defer close(gameIdsCh)
		ids, err := s.GetGameIds()
		gameIdsCh <- manifest.ManifestGameGetterGameIds{
			Ids: ids,
			Error: err,
		}
	}()

	return gameIdsCh
}

func duplicateGameIdsChan(gameIdsCh <-chan manifest.ManifestGameGetterGameIds) (<-chan manifest.ManifestGameGetterGameIds, <-chan manifest.ManifestGameGetterGameIds) {
	gameIdsCh1 := make(chan manifest.ManifestGameGetterGameIds)
	gameIdsCh2 := make(chan manifest.ManifestGameGetterGameIds)

	go func() {
		defer close(gameIdsCh1)
		defer close(gameIdsCh2)

		res := <-gameIdsCh
		gameIdsCh1 <-res
		gameIdsCh2 <-res
	}()

	return gameIdsCh1, gameIdsCh2
}

func getStorageGames(s Storage, done <-chan struct{}, gameIdsCh <-chan manifest.ManifestGameGetterGameIds) <-chan manifest.ManifestGameGetterGame {
	gamesCh := make(chan manifest.ManifestGameGetterGame)

	go func() {
		defer close(gamesCh)

		select {
		case res := <-gameIdsCh:
			if res.Error != nil {
				return
			}

			for _, id := range res.Ids {
				select {
				case <-done:
					return
				default:				
				}
	
				files, filesErr := s.GetGameFiles(id)
				if filesErr != nil {
					gamesCh <- manifest.ManifestGameGetterGame{
						Game: manifest.ManifestGame{},
						Warnings: []error{},
						Errors: []error{filesErr},
					}
					return
				}
	
				game := manifest.ManifestGame{
					Id:         id,
					Tags:       []string{},
					Installers: []manifest.ManifestGameInstaller{},
					Extras:     []manifest.ManifestGameExtra{},
				}

				for _, file := range files {
					if file.Kind == "installer" {
						game.Installers = append(game.Installers, manifest.ManifestGameInstaller{Name: file.Name})
					} else {
						game.Extras = append(game.Extras, manifest.ManifestGameExtra{Name: file.Name})
					}
				}

				gamesCh <- manifest.ManifestGameGetterGame{
					Game: game,
					Warnings: []error{},
					Errors: []error{},
				}
			}
		case <-done:
			return
		}
	}()

	return gamesCh
}

func addStorageGamesFilesMetadata(s Storage, concurrency int, done <-chan struct{}, game <-chan manifest.ManifestGameGetterGame) <-chan manifest.ManifestGameGetterGame {
	gamesCh := make(chan manifest.ManifestGameGetterGame)

	go func() {
		defer close(gamesCh)
	}()

	return gamesCh
}

func GenerateManifestGameGetter(s Storage, concurrency int) manifest.ManifestGameGetter {
	return func(done <-chan struct{}, gameIds []int64, filter manifest.ManifestFilter) (<-chan manifest.ManifestGameGetterGame, <-chan manifest.ManifestGameGetterGameIds) {
		gameIdsResultCh, gameIdsCh2 := duplicateGameIdsChan(getStorageGameIds(s))
		gameResultCh :=  addStorageGamesFilesMetadata(s, concurrency, done, getStorageGames(s, done, gameIdsCh2))
		return gameResultCh, gameIdsResultCh
	}
}
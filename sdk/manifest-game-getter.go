package sdk

import (
	"gogcli/manifest"
	
	"strings"
	"sync"
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

type GameResult struct {
	Game  manifest.ManifestGame
	Error error
}

type GameIdsResult struct {
	Ids   []int64
	Error error
}

func OwnedGamePagesToGames(done <-chan struct{}, ownedGamesPageCh <-chan OwnedGamesPageReturn, gameIds []int64, filter manifest.ManifestFilter) <-chan GameResult {
	gameCh := make(chan GameResult)
	
	go func() {
		defer close(gameCh)
		titles := filter.Titles
		for true {
			select {
			case pageRes, ok := <- ownedGamesPageCh:
				if !ok {
					return
				}

				if pageRes.err != nil {
					gameCh <- GameResult{
						Game: manifest.ManifestGame{},
						Error: pageRes.err,
					}
					continue
				}

				for _, product := range pageRes.page.Products {
					if len(gameIds) > 0 && (!contains(gameIds, product.Id)) {
						continue
					}

					game := manifest.ManifestGame{
						Id:    product.Id,
						Title: product.Title,
						Slug:  product.Slug,
					}

					if len(titles) > 0 && (!game.HasTitleTerms(titles)) {
						continue
					}

					gameCh <- GameResult{
						Game: game,
						Error: nil,
					}
				}
			case <-done:
				return
			}
		}
	}()

	return gameCh
}

func (s *Sdk) AddGameDetailsToGames(done <-chan struct{}, inGameCh <-chan GameResult, concurrency int, pause int, filter manifest.ManifestFilter) <-chan GameResult {
	var wg sync.WaitGroup
	outGameCh := make(chan GameResult)

	for idx := 0; idx < concurrency; idx++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tags := filter.Tags
			for true {
				select {
				case gameRes, ok := <-inGameCh:
					if !ok {
						return
					}

					if gameRes.Error != nil {
						outGameCh <- gameRes
					}

					gd, err := s.GetGameDetails(gameRes.Game.Id)
					if err != nil {
						gameRes.Error = err
						outGameCh <- gameRes
					}
					game := gameRes.Game

					game.Tags = make([]string, len(gd.Tags))
					for i, _ := range gd.Tags {
						game.Tags[i] = gd.Tags[i].Name
					}

					if len(tags) > 0 && (!game.HasOneOfTags(tags)) {
						continue
					}

					game.CdKey = gd.CdKey

					for _, i := range gd.Downloads {
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

					for _, e := range gd.Extras {
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

					for _, d := range gd.Dlcs {
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

					outGameCh <- GameResult{
						Game: game,
						Error: nil,
					}
				case <-done:
					return
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(outGameCh)
	}()

	return outGameCh
}

func TapGameIds(done <-chan struct{}, inGameCh <-chan GameResult) (<-chan GameResult, <-chan GameIdsResult) {
	outGameCh := make(chan GameResult)
	outGameIdsCh := make(chan GameIdsResult)

	go func() {
		defer close(outGameIdsCh)
		defer close(outGameCh)

		games := []manifest.ManifestGame{}
		gameIds := []int64{}
		for true {
			select {
			case gameRes, ok := <-inGameCh:
				if !ok {
					outGameIdsCh <- GameIdsResult{Ids: gameIds, Error: nil}
					for _, game := range games {
						outGameCh <- GameResult{Game: game, Error: nil}
					}
					return
				}
				
				if gameRes.Error != nil {
					outGameIdsCh <- GameIdsResult{Ids: []int64{}, Error: gameRes.Error}
				}

				games = append(games, gameRes.Game)
				gameIds = append(gameIds, gameRes.Game.Id)
			case <-done:
				return
			}
		}
	}()

	return outGameCh, outGameIdsCh
}

func (s *Sdk) GenerateManifestGameGetter(f manifest.ManifestFilter, concurrency int, pause int, tolerateDangles bool, tolerateBadMetadata bool) manifest.ManifestGameGetter {
	return func(done <-chan struct{}, gameIds []int64, filter manifest.ManifestFilter) (<-chan manifest.ManifestGameGetterGame, <-chan manifest.ManifestGameGetterGameIds) {
		gameResultCh := make(chan manifest.ManifestGameGetterGame)
		gameIdsResultCh := make(chan manifest.ManifestGameGetterGameIds)

		gamesCh, gameIdsCh := TapGameIds(
			done,
			s.AddGameDetailsToGames(
				done, 
				OwnedGamePagesToGames(
					done, 
					s.GetAllOwnedGamesPages(done, "", concurrency, pause), 
					gameIds, 
					filter,
				),
				concurrency, 
				pause,
				filter,
			),
		)

		go func() {
			defer close(gameIdsResultCh)
			defer close(gameResultCh)

			select {
			case gameIdsRes := <-gameIdsCh:
				gameIdsResultCh <- manifest.ManifestGameGetterGameIds{
					Ids: gameIdsRes.Ids,
					Error: gameIdsRes.Error,
				}
			case <-done:
				return
			}

			games := []manifest.ManifestGame{}
			for true {
				select {
				case gameRes, ok := <-gamesCh:
					if !ok {
						break
					}

					if gameRes.Error != nil {
						gameResultCh <- manifest.ManifestGameGetterGame{
							Game: gameRes.Game,
							Warnings: []error{},
							Errors: []error{gameRes.Error},
						}
						continue
					}

					games = append(games, gameRes.Game)
				case <-done:
					return
				}
			}
		}()
		
		return gameResultCh, gameIdsResultCh
	}
}
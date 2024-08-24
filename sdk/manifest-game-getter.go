package sdk

import (
	"gogcli/manifest"
	
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
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

type ManifestGameResult struct {
	Game  manifest.ManifestGame
	Error error
}

type ManifestGameIdsResult struct {
	Ids   []int64
	Error error
}

type GameManyErrorsResult struct {
	Game  manifest.ManifestGame
	Warnings []error
	Errors   []error
}

func OwnedGamePagesToManifestGames(done <-chan struct{}, ownedGamesPageCh <-chan OwnedGamesPageReturn, gameIds []int64, filter manifest.ManifestFilter) <-chan ManifestGameResult {
	gameCh := make(chan ManifestGameResult)
	
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
					gameCh <- ManifestGameResult{
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

					gameCh <- ManifestGameResult{
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

func (s *Sdk) AddGameDetailsToGames(done <-chan struct{}, inGameCh <-chan ManifestGameResult, concurrency int, pause int, filter manifest.ManifestFilter) <-chan ManifestGameResult {
	var wg sync.WaitGroup
	outGameCh := make(chan ManifestGameResult)

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
						continue
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
						if i.Size == "0 MB" {
							continue
						}

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
						if e.Size == "0 MB" {
							continue
						}

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
							if i.Size == "0 MB" {
								continue
							}
							
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
							if e.Size == "0 MB" {
								continue
							}

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

					outGameCh <- ManifestGameResult{
						Game: game,
						Error: nil,
					}
				case <-done:
					return
				}

				time.Sleep(time.Duration(pause) * time.Millisecond)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(outGameCh)
	}()

	return outGameCh
}

func TapManifestGameIds(done <-chan struct{}, inGameCh <-chan ManifestGameResult) (<-chan ManifestGameResult, <-chan ManifestGameIdsResult) {
	outGameCh := make(chan ManifestGameResult)
	outGameIdsCh := make(chan ManifestGameIdsResult)

	go func() {
		defer close(outGameIdsCh)
		defer close(outGameCh)

		games := []manifest.ManifestGame{}
		gameIds := []int64{}
		for true {
			select {
			case gameRes, ok := <-inGameCh:
				if !ok {
					outGameIdsCh <- ManifestGameIdsResult{Ids: gameIds, Error: nil}
					for _, game := range games {
						outGameCh <- ManifestGameResult{Game: game, Error: nil}
					}
					return
				}
				
				if gameRes.Error != nil {
					outGameIdsCh <- ManifestGameIdsResult{Ids: []int64{}, Error: gameRes.Error}
					continue
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

func (s *Sdk) AddFileInfoToGames(done <-chan struct{}, inGameCh <-chan ManifestGameResult, concurrency int, pause int, tolerateDangles bool, tolerateBadMetadata bool, filter manifest.ManifestFilter) <-chan GameManyErrorsResult {
	var wg sync.WaitGroup
	outGameCh := make(chan GameManyErrorsResult)

	for idx := 0; idx < concurrency; idx++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for true {
				select {
				case gameRes, ok := <-inGameCh:
					if !ok {
						return
					}

					if gameRes.Error != nil {
						outGameCh <- GameManyErrorsResult{
							Game: gameRes.Game,
							Warnings: []error{},
							Errors: []error{gameRes.Error},
						}
					}

					warnings := []error{}
					errs := []error{}

					game := gameRes.Game
					game.TrimFilesFromFilter(filter)
					if game.IsEmpty() {
						break
					}

					for idx, installer := range game.Installers {
						if len(errs) > 0 {
							break
						}

						info := s.GetFileInfo(installer.Url, tolerateBadMetadata)
						if info.Error != nil {
							if info.BadMetadata && tolerateBadMetadata {
								(*s).logger.Warning(fmt.Sprintf("Bad metadata for %s: File metadata was still fetched using much longer workaround method.", info.Url))
								err := errors.New(fmt.Sprintf("Bad metadata workaround: %s", info.Error.Error()))
								warnings = append(warnings, err)
								installer.Name = info.Name
								installer.Checksum = info.Checksum
								installer.VerifiedSize = info.Size
								game.Installers[idx] = installer
							} else if info.Dangling && tolerateDangles {
								(*s).logger.Warning(fmt.Sprintf("Bad download link for %s: File was not added to manifest.", info.Url))
								err := errors.New(fmt.Sprintf("Skipped File: %s", info.Error.Error()))
								warnings = append(warnings, err)
							} else {
								errs = append(errs, info.Error)
							}
							continue
						}
						installer.Name = info.Name
						installer.Checksum = info.Checksum
						installer.VerifiedSize = info.Size
						game.Installers[idx] = installer
					}

					for idx, extra := range game.Extras {
						if len(errs) > 0 {
							break
						}

						info := s.GetFileInfo(extra.Url, tolerateBadMetadata)
						if info.Error != nil {
							if info.BadMetadata && tolerateBadMetadata {
								(*s).logger.Warning(fmt.Sprintf("Bad metadata for %s: File metadata was still fetched using much longer workaround method.", info.Url))
								err := errors.New(fmt.Sprintf("Bad metadata workaround: %s", info.Error.Error()))
								warnings = append(warnings, err)
								extra.Name = info.Name
								extra.Checksum = info.Checksum
								extra.VerifiedSize = info.Size
								game.Extras[idx] = extra
							} else if info.Dangling && tolerateDangles {
								(*s).logger.Warning(fmt.Sprintf("Bad download link for %s: File was not added to manifest.", info.Url))
								err := errors.New(fmt.Sprintf("Skipped File: %s", info.Error.Error()))
								warnings = append(warnings, err)
							} else {
								errs = append(errs, info.Error)
							}
							continue
						}
						extra.Name = info.Name
						extra.Checksum = info.Checksum
						extra.VerifiedSize = info.Size
						game.Extras[idx] = extra
					}

					outGameCh <- GameManyErrorsResult{
						Game: game,
						Warnings: warnings,
						Errors: errs,
					}
				case <-done:
					return
				}

				time.Sleep(time.Duration(pause) * time.Millisecond)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(outGameCh)
	}()

	return outGameCh
}

func (s *Sdk) GenerateManifestGameGetter(concurrency int, pause int, tolerateDangles bool, tolerateBadMetadata bool) manifest.ManifestGameGetter {
	return func(done <-chan struct{}, gameIds []int64, filter manifest.ManifestFilter) (<-chan manifest.ManifestGameGetterGame, <-chan manifest.ManifestGameGetterGameIds) {
		gameResultCh := make(chan manifest.ManifestGameGetterGame)
		gameIdsResultCh := make(chan manifest.ManifestGameGetterGameIds)

		gamesCh, gameIdsCh := TapManifestGameIds(
			done,
			s.AddGameDetailsToGames(
				done, 
				OwnedGamePagesToManifestGames(
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

		gamesFinalCh := s.AddFileInfoToGames(done, gamesCh, concurrency, pause, tolerateDangles, tolerateBadMetadata, filter)

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

			for true {
				select {
				case gameRes, ok := <-gamesFinalCh:
					if !ok {
						return
					}

					gameResultCh <- manifest.ManifestGameGetterGame{
						Game: gameRes.Game,
						Warnings: gameRes.Warnings,
						Errors: gameRes.Errors,
					}
				case <-done:
					return
				}
			}
		}()
		
		return gameResultCh, gameIdsResultCh
	}
}
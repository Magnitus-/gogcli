package sdk

import (
	"gogcli/metadata"
	
	"fmt"
	"sync"
	"time"
)

type MetadataGameResult struct {
	Game     metadata.MetadataGame
	Warnings []error
	Error    error
}

type MetadataGameIdsResult struct {
	Ids   []int64
	Error error
}

func processProductImageUrl(u string) string {
	return fmt.Sprintf("%s%s", "https:", ensureFileNameSuffix(u, ".jpg"))
}

func OwnedGamePagesToMetadataGames(done <-chan struct{}, ownedGamesPageCh <-chan OwnedGamesPageReturn, gameIds []int64) <-chan MetadataGameResult {
	gameCh := make(chan MetadataGameResult)
	
	go func() {
		defer close(gameCh)
		for true {
			select {
			case pageRes, ok := <- ownedGamesPageCh:
				if !ok {
					return
				}

				if pageRes.err != nil {
					gameCh <- MetadataGameResult{
						Game:     metadata.MetadataGame{},
						Warnings: []error{},
						Error:    pageRes.err,
					}
					continue
				}

				for _, product := range pageRes.page.Products {
					if len(gameIds) > 0 && (!contains(gameIds, product.Id)) {
						continue
					}

					game := metadata.MetadataGame{
						Id:           product.Id,
						Title:        product.Title,
						Slug:         product.Slug,
						Category:     product.Category,
						Rating:       product.Rating,
						Dlcs:         product.DlcCount,
						ListingImage: metadata.GameMetadataImage{
							Url: processProductImageUrl(product.Image),
							Tag: "Listing",
						},
					}

					gameCh <- MetadataGameResult{
						Game:     game,
						Warnings: []error{},
						Error:    nil,
					}
				}
			case <-done:
				return
			}
		}
	}()

	return gameCh
}

func TapMetadataGameIds(done <-chan struct{}, inGameCh <-chan MetadataGameResult) (<-chan MetadataGameResult, <-chan MetadataGameIdsResult) {
	outGameCh := make(chan MetadataGameResult)
	outGameIdsCh := make(chan MetadataGameIdsResult)

	go func() {
		defer close(outGameIdsCh)
		defer close(outGameCh)

		games := []metadata.MetadataGame{}
		gameIds := []int64{}
		for true {
			select {
			case gameRes, ok := <-inGameCh:
				if !ok {
					outGameIdsCh <- MetadataGameIdsResult{Ids: gameIds, Error: nil}
					for _, game := range games {
						outGameCh <- MetadataGameResult{Game: game, Warnings: []error{}, Error: nil}
					}
					return
				}
				
				if gameRes.Error != nil {
					outGameIdsCh <- MetadataGameIdsResult{Ids: []int64{}, Error: gameRes.Error}
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

func (s *Sdk) AddProductsInfoToMetadataGames(done <-chan struct{}, inGameCh <-chan MetadataGameResult, concurrency int, pause int, skipImages []string) <-chan MetadataGameResult {
	var wg sync.WaitGroup
	outGameCh := make(chan MetadataGameResult)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for true {
				select {
				case gameRes, ok := <- inGameCh:
					if !ok {
						return
					}

					if gameRes.Error != nil {
						outGameCh <- gameRes
						continue
					}

					product, dangling, err := s.GetProduct(gameRes.Game.Id)
					if err != nil {
						if dangling {
							gameRes.Warnings = []error{err}
							gameRes.Game.HasProductInfo = false
							outGameCh <- gameRes
							continue
						}
						gameRes.Error = err
						outGameCh <- gameRes
						continue
					}

					game := gameRes.Game
					game.HasProductInfo = true

					game.Description = metadata.GameMetadataDescription{
						Summary:    product.Description.Lead,
						Full:       product.Description.Full,
						Highlights: product.Description.Whats_cool_about_it,
					}
					game.ReleaseDate = product.Release_date
					game.Changelog = product.Changelog
			
					videos := []metadata.GameMetadataVideo{}
					for _, vid := range product.Videos {
						videos = append(videos, metadata.GameMetadataVideo{
							ThumbnailUrl: vid.Thumbnail_url,
							VideoUrl:     vid.Video_url,
							Provider:     vid.Provider,
						})
					}
					game.Videos = videos
					
					backgroundUrl := ""
					if (!containsStr(skipImages, processProductImageUrl(product.Images.Background))) && product.Images.Background != "" {
						backgroundUrl = processProductImageUrl(product.Images.Background)
					}

					logoUrl := ""
					if (!containsStr(skipImages, processProductImageUrl(product.Images.Logo2x))) && product.Images.Logo2x != "" {
						logoUrl = processProductImageUrl(product.Images.Logo2x)
					}

					iconUrl := ""
					if (!containsStr(skipImages, processProductImageUrl(product.Images.Icon))) && product.Images.Icon != "" {
						iconUrl = processProductImageUrl(product.Images.Icon)
					}

					game.ProductImages = metadata.GameMetadataProductImages{
						Background: metadata.GameMetadataImage{
							Url: backgroundUrl,
							Tag: "Background",
						},
						Logo: metadata.GameMetadataImage{
							Url: logoUrl,
							Tag: "Logo",
						},
						Icon: metadata.GameMetadataImage{
							Url: iconUrl,
							Tag: "Icon",
						},
					}
			
					screenshots := []metadata.GameMetadataScreenShot{}
					for _, prodScreenshot := range product.Screenshots {
						isBad := false
						screenshot := metadata.GameMetadataScreenShot{}
						for _, prodScreenshotRes := range prodScreenshot.Formatted_images {
							if containsStr(skipImages, prodScreenshotRes.Image_url) {
								isBad = true
								break
							}

							if prodScreenshotRes.Formatter_name == "ggvgm" {
								screenshot.List = metadata.GameMetadataImage{
									Url: prodScreenshotRes.Image_url,
									Tag: "List",
								}
							}
							if prodScreenshotRes.Formatter_name == "ggvgl_2x" {
								screenshot.Main = metadata.GameMetadataImage{
									Url: prodScreenshotRes.Image_url,
									Tag: "Main",
								}
							}
						}

						if !isBad {
							screenshots = append(screenshots, screenshot)
						}
					}
					game.Screenshots = screenshots

					gameRes.Game = game
					outGameCh <- gameRes
				case <- done:
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

func fillGameMetadataImageInfo(gameImage *metadata.GameMetadataImage, info *ImageInfo) {
	(*gameImage).Name = (*info).Name
	(*gameImage).Size = (*info).Size
	(*gameImage).Checksum = (*info).Checksum
}

func (s *Sdk) AddImagesInfoToMetadataGames(done <-chan struct{}, inGameCh <-chan MetadataGameResult, concurrency int, pause int) <-chan MetadataGameResult {
	var wg sync.WaitGroup
	outGameCh := make(chan MetadataGameResult)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for true {
				select{
				case gameRes, ok := <- inGameCh:
					if !ok {
						return
					}

					if gameRes.Error != nil || !gameRes.Game.HasProductInfo {
						outGameCh <- gameRes
						continue
					}

					game := gameRes.Game

					gameImagesPtrs := game.GetImagesPointers()
					for _, ptr := range gameImagesPtrs {
						info, err := s.GetImageInfoWithFallback(ptr)
						if err != nil {
							gameRes.Error = err
							outGameCh <- gameRes
							break
						}
						fillGameMetadataImageInfo(ptr, &info)
					}

					gameRes.Game = game
					outGameCh <- gameRes
				case <- done:
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

func (s *Sdk) GenerateMetadataGameGetter(concurrency int, pause int, tolerateDangles bool) metadata.MetadataGameGetter {
	return func(done <-chan struct{}, gameIds []int64, skipImages []string) (<-chan metadata.MetadataGameGetterGame, <-chan metadata.MetadataGameGetterGameIds) {
		gameResultCh := make(chan metadata.MetadataGameGetterGame)
		gameIdsResultCh := make(chan metadata.MetadataGameGetterGameIds)

		gamesCh, gameIdsCh := TapMetadataGameIds(
			done,
			OwnedGamePagesToMetadataGames(
				done, 
				s.GetAllOwnedGamesPages(done, "", concurrency, pause), 
				gameIds,
			),
		)

		gamesCh = s.AddImagesInfoToMetadataGames(
			done,
			s.AddProductsInfoToMetadataGames(done, gamesCh, concurrency, pause, skipImages),
			concurrency,
			pause,
		)

		go func() {
			defer close(gameIdsResultCh)
			defer close(gameResultCh)

			select {
			case gameIdsRes := <-gameIdsCh:
				gameIdsResultCh <- metadata.MetadataGameGetterGameIds{
					Ids: gameIdsRes.Ids,
					Error: gameIdsRes.Error,
				}
			case <-done:
				return
			}

			for true {
				select {
				case gameRes, ok := <-gamesCh:
					if !ok {
						return
					}

					if gameRes.Error != nil {
						gameResultCh <- metadata.MetadataGameGetterGame{
							Game: gameRes.Game,
							Warnings: []error{},
							Errors: []error{gameRes.Error},
						}
						continue
					}

					gameResultCh <- metadata.MetadataGameGetterGame{
						Game: gameRes.Game,
						Warnings: gameRes.Warnings,
						Errors: []error{},
					}
				case <-done:
					return
				}
			}
		}()

		return gameResultCh, gameIdsResultCh
	}
}
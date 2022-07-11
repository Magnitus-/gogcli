package sdk

import (
	"gogcli/metadata"
	
	//"errors"
	"fmt"
	"strings"
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
	return fmt.Sprintf("%s%s", "https:", ensureFileNameSuffix(u, ".png"))
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
							Tag: "Logo",
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

/*
Left:
type MetadataGame struct {
	ProductCards          []Image
	Features              []string
}
*/

func (s *Sdk) AddProductsInfoToMetadataGames(done <-chan struct{}, inGameCh <-chan MetadataGameResult, concurrency int, pause int) <-chan MetadataGameResult {
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
							outGameCh <- gameRes
							continue
						}
						gameRes.Error = err
						outGameCh <- gameRes
						continue
					}

					game := gameRes.Game

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
			
					game.ProductImages = metadata.GameMetadataProductImages{
						Background: metadata.GameMetadataImage{
							Url: processProductImageUrl(product.Images.Background),
							Tag: "Background",
						},
						Logo: metadata.GameMetadataImage{
							Url: processProductImageUrl(product.Images.Logo),
							Tag: "Logo",
						},
						Logo2x: metadata.GameMetadataImage{
							Url: processProductImageUrl(product.Images.Logo2x),
							Tag: "Logo2x",
						},
						Icon: metadata.GameMetadataImage{
							Url: processProductImageUrl(product.Images.Icon),
							Tag: "Icon",
						},
						SidebarIcon: metadata.GameMetadataImage{
							Url: processProductImageUrl(product.Images.SidebarIcon),
							Tag: "SidebarIcon",
						},
						SidebarIcon2x: metadata.GameMetadataImage{
							Url: processProductImageUrl(product.Images.SidebarIcon2x),
							Tag: "SidebarIcon2x",
						},
						MenuNotificationAv: metadata.GameMetadataImage{
							Url: processProductImageUrl(product.Images.MenuNotificationAv),
							Tag: "MenuNotificationAv",
						},
						MenuNotificationAv2: metadata.GameMetadataImage{
							Url: processProductImageUrl(product.Images.MenuNotificationAv2),
							Tag: "MenuNotificationAv2",
						},
					}
			
					screenshots := []metadata.GameMetadataScreenShot{}
					for _, prodScreenshot := range product.Screenshots {
						screenshot := metadata.GameMetadataScreenShot{}
						for _, prodScreenshotRes := range prodScreenshot.Formatted_images {
							screenshot = append(screenshot, metadata.GameMetadataImage{
								Url: strings.Replace(prodScreenshotRes.Image_url, ".jpg", ".png", -1),
								Tag: prodScreenshotRes.Formatter_name,
							})
						}
						screenshots = append(screenshots, screenshot)
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

func (s *Sdk) GenerateMetadataGameGetter(concurrency int, pause int, tolerateDangles bool) metadata.MetadataGameGetter {
	return func(done <-chan struct{}, gameIds []int64) (<-chan metadata.MetadataGameGetterGame, <-chan metadata.MetadataGameGetterGameIds) {
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

		gamesCh = s.AddProductsInfoToMetadataGames(done, gamesCh, concurrency, pause)

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
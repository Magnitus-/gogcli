package sdk

import (
	"gogcli/metadata"
	
	/*"errors"
	"fmt"
	"strings"
	"sync"
	"time"*/
)

type MetadataGameResult struct {
	Game  metadata.MetadataGame
	Error error
}

type MetadataGameIdsResult struct {
	Ids   []int64
	Error error
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
						Game: metadata.MetadataGame{},
						Error: pageRes.err,
					}
					continue
				}

				for _, product := range pageRes.page.Products {
					if len(gameIds) > 0 && (!contains(gameIds, product.Id)) {
						continue
					}

					game := metadata.MetadataGame{
						Id:       product.Id,
						Title:    product.Title,
						Slug:     product.Slug,
						Category: product.Category,
						Rating:   product.Rating,
						Dlcs:     product.DlcCount,
					}

					gameCh <- MetadataGameResult{
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
						outGameCh <- MetadataGameResult{Game: game, Error: nil}
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

					gameResultCh <- metadata.MetadataGameGetterGame{
						Game: gameRes.Game,
						Warnings: []error{},
						Errors: []error{gameRes.Error},
					}
				case <-done:
					return
				}
			}
		}()

		return gameResultCh, gameIdsResultCh
	}
}
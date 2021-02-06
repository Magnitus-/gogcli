package sdk

import "time"

type GamesDetailsReturn struct {
	game GameDetails
	id   int
	err  error
}

func (s *Sdk) GetGameDetailsAsync(gameId int, returnVal chan GamesDetailsReturn) {
	g, err := s.GetGameDetails(gameId)
	returnVal <- GamesDetailsReturn{game: g, id: gameId, err: err}
}

type GameDetailsWithId struct {
	game GameDetails
	id   int
}

func (s *Sdk) GetManyGameDetails(gameIds []int, concurrency int, pause int) ([]GameDetailsWithId, []error) {
	var errs []error
	var games []GameDetailsWithId
	c := make(chan GamesDetailsReturn)

	i := 0
	for i < len(gameIds) {
		beginning := i
		target := min(len(gameIds), i+concurrency)
		for i < target {
			go s.GetGameDetailsAsync(gameIds[i], c)
			i++
		}

		y := beginning
		for y < target {
			returnVal := <-c
			if returnVal.err != nil {
				errs = append(errs, returnVal.err)
			} else {
				games = append(games, GameDetailsWithId{game: returnVal.game, id: returnVal.id})
			}
			y++
		}

		if len(errs) > 0 {
			return games, errs
		}

		if i < len(gameIds) {
			time.Sleep(time.Duration(pause) * time.Millisecond)
		}
	}

	return games, nil
}

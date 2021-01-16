package sdk

import (
	"fmt"
	"time"
)

type OwnedGamesPageReturn struct {
	page OwnedGamesPage
	err  error
}

func (s *Sdk) GetOwnedGamesPageAsync(page int, search string, debug bool, returnVal chan OwnedGamesPageReturn) {
	o, err := s.GetOwnedGames(page, search, debug)
	returnVal <- OwnedGamesPageReturn{page: o, err: err}
}

func (s *Sdk) GetAllOwnedGamesPages(search string, concurrency int, pause int, debug bool) ([]OwnedGamesPage, []error) {
	fnCall := fmt.Sprintf(
		"GetAllOwnedGamesPages(search=%s, concurrency=%d, pause=%d",
		search,
		concurrency,
		pause,
	)
	var pageCount int
	var currentPage int
	var pages []OwnedGamesPage
	var errs []error
	var callVal OwnedGamesPageReturn
	c := make(chan OwnedGamesPageReturn)

	if debug {
		(*s).logger.Println(
			fmt.Sprintf("%s -> fetching first page", fnCall),
		)
	}
	go s.GetOwnedGamesPageAsync(1, search, debug, c)
	callVal = <-c
	if callVal.err != nil {
		errs = append(errs, callVal.err)
		return pages, errs
	}

	if callVal.page.TotalPages == 0 {
		return pages, errs
	}

	pages = append(pages, callVal.page)
	pageCount = callVal.page.TotalPages
	currentPage = callVal.page.Page + 1

	if debug {
		(*s).logger.Println(
			fmt.Sprintf("%s -> total pages to fetch: %d", fnCall, pageCount),
		)
	}

	for currentPage < pageCount {
		maxPage := min(currentPage+concurrency, pageCount)

		if debug {
			(*s).logger.Println(
				fmt.Sprintf("%s -> fetching page %d to %d", fnCall, currentPage, maxPage),
			)
		}

		for i := currentPage + 1; i <= maxPage; i++ {
			go s.GetOwnedGamesPageAsync(i, search, debug, c)
		}

		for i := currentPage + 1; i <= maxPage; i++ {
			callVal = <-c
			if callVal.err != nil {
				errs = append(errs, callVal.err)
			} else {
				pages = append(pages, callVal.page)
			}
		}

		if len(errs) > 0 {
			return pages, errs
		}

		currentPage = maxPage
		time.Sleep(time.Duration(pause) * time.Millisecond)
	}
	return pages, errs
}

func min(x int, y int) int {
	if x < y {
		return x
	}
	return y
}

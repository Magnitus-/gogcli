package sdk

import (
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
	var pageCount int
	var currentPage int
	var pages []OwnedGamesPage
	var errs []error
	var callVal OwnedGamesPageReturn
	c := make(chan OwnedGamesPageReturn)

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

	for currentPage < pageCount {
		maxPage := min(currentPage+concurrency, pageCount)

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
